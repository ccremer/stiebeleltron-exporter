package cfg

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"stiebeleltron-exporter/pkg/metrics"
)

type (
	// Configuration holds a strongly-typed tree of the configuration
	Configuration struct {
		Log struct {
			Level   string
			Verbose bool
		}
		ISG struct {
			URL            string
			Timeout        time.Duration
			Headers        []string `koanf:"header"`
			DefinitionPath string
		}
		BindAddr string
	}
	MetricDefinitions struct {
		Pages map[string]Page
	}
	Page struct {
		Groups    map[string]Group
		URLSuffix string
	}
	Group struct {
		SearchString string
		Metrics      []Metric
	}
	Metric struct {
		Name         string
		Description  string
		SearchString string
		Multiplier   *float64
		Divisor      *float64
		Labels       prometheus.Labels
	}
)

// NewDefaultExporterConfig retrieves the hardcoded configs with sane defaults
func NewDefaultExporterConfig() *Configuration {
	c := &Configuration{}
	c.Log.Level = "info"
	c.ISG.URL = "http://isg.ip.or.hostname"
	c.ISG.Timeout = 5 * time.Second
	c.BindAddr = ":8080"
	return c
}

func (m Metric) GetMultiplier() float64 {
	if m.Multiplier == nil {
		return 1
	}
	return *m.Multiplier
}

func (m Metric) GetDivisor() float64 {
	if m.Divisor == nil {
		return 1
	}
	return *m.Divisor
}

func (definitions MetricDefinitions) MapToPrometheusMetric() map[string][]*metrics.PrometheusMetric {
	m := make(map[string][]*metrics.PrometheusMetric, 0)
	for _, page := range definitions.Pages {
		perPageMetrics := make([]*metrics.PrometheusMetric, 0)
		for groupName, group := range page.Groups {
			for _, metric := range group.Metrics {
				promMetric := &metrics.PrometheusMetric{
					GaugeName:            metric.Name,
					Group:                groupName,
					GroupSearchString:    group.SearchString,
					PropertySearchString: metric.SearchString,
					HelpText:             metric.Description,
					Labels:               metric.Labels,
				}
				if metric.Divisor != nil {
					promMetric.ValueTransformer = metrics.NewDivisorTransformer(*metric.Divisor)
				}
				if metric.Multiplier != nil {
					promMetric.ValueTransformer = metrics.NewMultiplierTransformer(*metric.Multiplier)
				}
				promMetric.InitializeMetric()
				perPageMetrics = append(perPageMetrics, promMetric)
			}
		}
		m[page.URLSuffix] = perPageMetrics
	}
	return m
}
