package main

import (
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	flag "github.com/spf13/pflag"
	"stiebeleltron-exporter/cfg"
	"stiebeleltron-exporter/pkg/stiebeleltron"
)

var (
	version     = "unknown"
	commit      = "dirty"
	date        = time.Now().String()
	config      = cfg.ParseConfig(version, commit, date, flag.NewFlagSet("main", flag.ExitOnError), os.Args[1:])
	promHandler = promhttp.Handler()
)

func main() {

	log.WithFields(log.Fields{
		"version": version,
		"commit":  commit,
		"date":    date,
	}).Info("Starting exporter.")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		log.WithFields(log.Fields{
			"uri":    r.RequestURI,
			"client": r.RemoteAddr,
		}).Debug("Accessed Root endpoint")
		http.Redirect(w, r, "/metrics", http.StatusMovedPermanently)
	})
	http.HandleFunc("/liveness", func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"uri":    r.RequestURI,
			"client": r.RemoteAddr,
		}).Debug("Accessed Liveness endpoint")
		w.WriteHeader(http.StatusNoContent)
	})

	headers := http.Header{}
	cfg.ConvertHeaders(config.ISG.Headers, &headers)
	client, err := stiebeleltron.NewISGClient(stiebeleltron.ClientOptions{
		BaseURL: config.ISG.URL,
		Headers: headers,
	})
	if err != nil {
		log.Fatal(err)
	}

	props := cfg.LoadMetricDefinitions(config).MapToPrometheusMetric()

	http.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		log.WithFields(log.Fields{
			"uri":    req.RequestURI,
			"client": req.RemoteAddr,
		}).Debug("Accessed Metrics endpoint")
		ScrapeISG(client, props)
		promHandler.ServeHTTP(w, req)
	})

	log.WithField("port", config.BindAddr).Info("Listening for scrapes.")
	log.WithError(http.ListenAndServe(config.BindAddr, nil)).Fatal("Shutting down.")
}
