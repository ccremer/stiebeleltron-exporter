package cfg

import (
	"fmt"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"net/http"
	"os"
	"strings"
	"time"
)

// ParseConfig overrides internal config defaults with an optional YAML file, then environment variables and lastly CLI flags.
// Ensures basic validation.
func ParseConfig(version, commit, date string, fs *flag.FlagSet, args []string) *Configuration {
	config := NewDefaultConfig()

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s (version %s, %s, %s):\n", os.Args[0], version, commit, date)
		fs.PrintDefaults()
	}
	fs.String("bindAddr", config.BindAddr, "IP Address to bind to listen for Prometheus scrapes")
	fs.String("log.level", config.Log.Level, "Logging level")
	fs.BoolP("log.verbose", "v", config.Log.Verbose, "Shortcut for --log.level=debug")
	fs.StringSlice("isg.header", []string{},
		"List of \"key: value\" headers to append to the requests going to Stiebel Eltron ISG")
	fs.StringP("isg.url", "u", config.ISG.URL, "Target URL of Stiebel Eltron ISG device")
	fs.Int64("isg.timeout", int64(config.ISG.Timeout.Seconds()),
		"Timeout in seconds when collecting metrics from Stiebel Eltron ISG. Should not be larger than the scrape interval")
	fs.String("config", "", "Configuration file that may hold translations of metric names. Accepts full and relative path to a .yaml file")

	if err := fs.Parse(args); err != nil {
		log.WithError(err).Fatal("Could not parse flags")
	}

	k := koanf.New(".")
	path, _ := fs.GetString("config")
	if path != "" {
		log.WithFields(log.Fields{
			"path": path,
		}).Info("Loading configuration")
		err := k.Load(file.Provider(path), yaml.Parser())
		if err != nil {
			log.WithError(err).Fatal("Could not load config file")
		}
	}

	err := k.Load(env.Provider("", ".", func(s string) string {
		return strings.Replace(strings.ToLower(s), "_", ".", -1)
	}), nil)
	if err != nil {
		log.WithError(err).Fatal("Could not load environment variables")
	}

	err = k.Load(posflag.Provider(fs, ".", k), nil)
	if err != nil {
		log.WithError(err).Fatal("Could not load CLI flags")
	}

	if err := k.Unmarshal("", config); err != nil {
		log.WithError(err).Fatal("Could not read config")
	}

	config.ISG.Timeout *= time.Second
	if config.Log.Verbose {
		config.Log.Level = "debug"
	}
	level, err := log.ParseLevel(config.Log.Level)
	if err != nil {
		log.WithError(err).Warn("Could not parse log level, fallback to info level")
		config.Log.Level = "info"
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(level)
	}
	log.WithField("config", *config).Debug("Parsed config")
	return config
}

// ConvertHeaders takes a list of `key=value` headers and adds those trimmed to the specified header struct. It ignores
// any malformed entries.
func ConvertHeaders(headers []string, header *http.Header) {
	for _, hd := range headers {
		arr := strings.SplitN(hd, "=", 2)
		if len(arr) < 2 {
			log.WithFields(log.Fields{
				"arg":   hd,
				"error": "cannot split: missing equal sign",
			}).Warn("Could not parse header, ignoring")
			continue
		}
		key := strings.TrimSpace(arr[0])
		value := strings.TrimSpace(arr[1])
		log.WithFields(log.Fields{
			"key":   key,
			"value": value,
		}).Debug("Using header")
		header.Set(key, value)
	}
}
