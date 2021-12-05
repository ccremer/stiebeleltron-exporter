package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type (
	Transformer func(v float64) float64
)

type PrometheusMetric struct {
	GaugeName            string
	Group                string
	GroupSearchString    string
	PropertySearchString string
	HelpText             string
	Labels               prometheus.Labels
	Gauge                prometheus.Gauge
	ValueTransformer     Transformer
}

var (
	Namespace = "stiebeleltron"
)

func (p *PrometheusMetric) GetGroup() string {
	return p.GroupSearchString
}

func (p *PrometheusMetric) GetSearchString() string {
	return p.PropertySearchString
}

func (p *PrometheusMetric) InitializeMetric() {
	p.Gauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   Namespace,
		Subsystem:   p.Group,
		Name:        p.GaugeName,
		Help:        p.HelpText,
		ConstLabels: p.Labels,
	})
	prometheus.MustRegister(p.Gauge)
}

func (p *PrometheusMetric) SetValue(v float64) {
	if p.ValueTransformer == nil {
		p.Gauge.Set(v)
		return
	}
	p.Gauge.Set(p.ValueTransformer(v))
}

func NewDivisorTransformer(divisor float64) Transformer {
	if divisor == 0 {
		log.Fatal("cannot use 0 as a divisor!")
	}
	return func(v float64) float64 {
		return v / divisor
	}
}

func NewMultiplierTransformer(multiplier float64) Transformer {
	return func(v float64) float64 {
		return v * multiplier
	}
}
