package main

import (
	"context"
	"sync"
	"time"

	"github.com/ccremer/stiebeleltron-exporter/pkg/metrics"
	"github.com/ccremer/stiebeleltron-exporter/pkg/stiebeleltron"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
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

func scrapeISG(c *stiebeleltron.ISGClient, m map[string][]*metrics.PrometheusMetric) {
	start := time.Now()
	defer func() {
		scrapeDurationGauge.Set(time.Since(start).Seconds())
	}()

	ctx, cancel := context.WithTimeout(context.Background(), config.ISG.Timeout)
	defer cancel()

	select {
	case <-ctx.Done():
		log.WithField("timeout", config.ISG.Timeout.Seconds()).Warn("Scrape timed out")
		scrapeErrorCounter.Inc()
	case <-fanoutScrape(c, m):
		log.WithFields(log.Fields{
			"duration": time.Since(start).Seconds(),
		}).Debug("Scrape completed")
	}
}

func fanoutScrape(c *stiebeleltron.ISGClient, m map[string][]*metrics.PrometheusMetric) <-chan error {
	respChan := make(chan error, 1)
	go func() {
		wg := &sync.WaitGroup{}
		wg.Add(len(m))
		for urlSuffix, metricList := range m {
			list := make([]stiebeleltron.Property, len(metricList))
			for i := range metricList {
				list[i] = metricList[i]
			}
			go scrapeSinglePage(urlSuffix, list, c, respChan, wg)
		}
		wg.Wait()
		respChan <- nil
	}()
	return respChan
}

func scrapeSinglePage(urlSuffix string, metricList []stiebeleltron.Property, c *stiebeleltron.ISGClient, respChan chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	scrapeLog := log.WithFields(log.Fields{"page": urlSuffix})
	parseErrors, err := c.ParsePage(urlSuffix, metricList)
	if err != nil {
		scrapeErrorCounter.Inc()
		scrapeLog.WithError(err).Error("Could not scrape page")
		respChan <- err
		return
	}
	for _, parseError := range parseErrors {
		scrapeLog.WithFields(log.Fields{
			"property": parseError.Property,
			"value":    parseError.RawText,
			"error":    parseError.Error,
		}).Warn("Could not parse property")
		parseCounter.Inc()
	}
	scrapeLog.Debug("Parsed page")
}
