package main

import (
	"flag"
	"net/http"
	"os"
	"temperature_exporter/collector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var (
	listenPort       int
	addr             = flag.String("listen-address", "9001", "The address to listen on for HTTP requests.")
	metricsPath      = flag.String("web.telemetry-path", "/metrics", "A path under which to expose metrics.")
	metricsNamespace = flag.String("metrics.namespace", "raspberry", "Prometheus metrics namespace, as the prefix of metrics name")
)

func init() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	flag.Parse()
	log.Info("Init Temperature Exporter Settings")

	metrics := collector.NewMetrics(*metricsNamespace)
	registry := prometheus.NewRegistry()
	registry.MustRegister(metrics)

	http.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>A Prometheus Exporter</title></head>
			<body>
			<h1>A Prometheus Exporter</h1>
			<p><a href='/metrics'>Metrics</a></p>
			</body>
			</html>`))
	})
	log.Info("Start to listen on Port:", *addr)
	err := http.ListenAndServe(":"+*addr, nil)
	if err != nil {
		log.Panic(err)
		os.Exit(1)
	}
}
