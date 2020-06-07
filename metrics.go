package main

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	"stiebeleltron-exporter/pkg/metrics"
	"stiebeleltron-exporter/pkg/stiebeleltron"
	"time"
)

var (
	scrapeErrorCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: metrics.Namespace,
		Name:      "scrape_errors_total",
		Help:      "Scrape errors can be used to monitor whether ISG is responsive",
	})
	parseCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: metrics.Namespace,
		Name:      "parse_errors_total",
		Help:      "Parsing errors when extracting the ISG HTML pages for metric properties",
	})
	scrapeDurationGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: metrics.Namespace,
		Name:      "scrape_duration_seconds",
		Help:      "Total scrape duration in seconds",
	})
)

func MergeProperties(metrics map[string]*metrics.MetricProperty, defaults map[string]stiebeleltron.Property) map[string]*metrics.MetricProperty {
	for key, value := range metrics {
		def, exists := defaults[key]
		if exists {
			value.PropertyGroup = getOrDefault(value.PropertyGroup, def.GetGroup())
			value.SearchString = getOrDefault(value.SearchString, def.GetSearchString())
		}
	}
	return metrics
}

func PrepareGauges(m map[string]*metrics.MetricProperty) {
	for _, value := range m {
		value.Gauge = promauto.NewGauge(prometheus.GaugeOpts{
			Namespace:   metrics.Namespace,
			Help:        value.HelpText,
			ConstLabels: value.Labels,
			Name:        value.GaugeName,
		})
	}
}

func getOrDefault(value, def string) string {
	if value != "" {
		return value
	}
	return def
}

func ScrapeISG(c *stiebeleltron.ISGClient, m map[string]*metrics.MetricProperty) {
	start := time.Now()
	defer func() {
		scrapeDurationGauge.Set(time.Since(start).Seconds())
	}()

	ctx, _ := context.WithTimeout(context.Background(), config.ISG.Timeout)

	select {
	case <-ctx.Done():
		log.WithField("timeout", config.ISG.Timeout.Seconds()).Warn("Scrape timed out")
		scrapeErrorCounter.Inc()
	case <-scrapeISGasync(c, m):
		log.WithFields(log.Fields{
			"duration": time.Since(start).Seconds(),
		}).Debug("Scrape completed")
	}

}

func scrapeISGasync(c *stiebeleltron.ISGClient, m map[string]*metrics.MetricProperty) <-chan error {
	respChan := make(chan error, 1)
	go func() {
		err := c.LoadSystemInfoPage()
		if err != nil {
			scrapeErrorCounter.Inc()
			log.WithError(err).Error("Could not scrape system info page")
			respChan <- err
			return
		}
		parseErrors := c.ParsePage(setMetricDelegate(m))
		log.Debug("Parsed system info page")

		err = c.LoadHeatPumpInfoPage()
		if err != nil {
			scrapeErrorCounter.Inc()
			log.WithError(err).Error("Could not scrape heat pump info page")
			respChan <- err
			return
		}
		parseErrors = append(parseErrors, c.ParsePage(setMetricDelegate(m))...)
		log.Debug("Parsed heat pump info page")

		for _, parseError := range parseErrors {
			log.WithFields(log.Fields{
				"group":    parseError.Group,
				"property": parseError.Property,
				"value":    parseError.Value,
				"error":    parseError.Error,
			}).Warn("Could not parse property")
			parseCounter.Inc()
		}
		respChan <- err
	}()
	return respChan
}

func setMetricDelegate(m map[string]*metrics.MetricProperty) func(string, string, float64) {
	return func(group, key string, value float64) {
		for prop, ma := range m {
			if ma.PropertyGroup == group && ma.SearchString == key {
				if ma.ValueTransformer != nil {
					value = ma.ValueTransformer(value)
				}
				log.WithFields(log.Fields{
					"property":   prop,
					"parseError": value,
					"metric":     ma.GaugeName,
					"labels":     ma.Labels,
				}).Debug("Assigned value")
				ma.Gauge.Set(value)
				return
			}
		}
		log.WithFields(log.Fields{
			"group":    group,
			"property": key,
			"value":    value,
		}).Warn("Could not find a matching API property")
	}
}
