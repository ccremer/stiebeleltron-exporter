package cfg

import "time"

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
			Headers []string `mapstructure:"header"`
		}
		BindAddr string
	}
)

// NewDefaultConfig retrieves the hardcoded configs with sane defaults
func NewDefaultConfig() *Configuration {
	c := &Configuration{}
	c.Log.Level = "info"
	c.ISG.URL = "http://isg.ip.or.hostname"
	c.ISG.Timeout = 5 * time.Second
	c.BindAddr = ":8080"
	return c
}
