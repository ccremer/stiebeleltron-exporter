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
			URL     string
			Timeout time.Duration
			Headers []string `koanf:"header"`
		}
		BindAddr   string
		Properties map[string]MetricProperty
	}
	MetricProperty struct {
		GaugeName        string
		Labels           map[string]string
		HelpText         string
		Gauge            prometheus.Gauge
		PropertyGroup    string
		SearchString     string
		ValueTransformer func(v float64) float64
	}
)

// NewDefaultConfig retrieves the hardcoded configs with sane defaults
func NewDefaultConfig() *Configuration {
	c := &Configuration{}
	c.Log.Level = "info"
	c.ISG.URL = "http://isg.ip.or.hostname"
	c.ISG.Timeout = 5 * time.Second
	c.BindAddr = ":8080"
	//c.Properties = metrics.NewDefaultMetricProperties()
	return c
}
