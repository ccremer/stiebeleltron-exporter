package cfg

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
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
		Pages []Page
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
		Name        string
		Description string
		Multiplier  *float64
		Divisor     *float64
		Labels      prometheus.Labels
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
