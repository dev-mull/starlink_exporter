package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

const (
	metricsPath = "/metrics"

	// DishAddress to reach Starlink dish ip:port
	DishAddress = "192.168.100.1:9200"
)

var (
	listen, address string
)

func init() {
	flag.StringVar(&listen, "listen", ":9817", "listening port to expose metrics on")
	flag.StringVar(&address, "dish-address", DishAddress, "IP address and port to reach dish")
}

func main() {
	flag.Parse()

	exporter, err := NewExporter(address)
	if err != nil {
		log.Warnf("could not talk to dishy: %s", err.Error())
	}
	defer exporter.Close()
	log.Infof("dish id: %s", exporter.DishID)

	r := prometheus.NewRegistry()
	r.MustRegister(exporter)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
             <head><title>Starlink Exporter</title></head>
             <body>
             <h1>Starlink Exporter</h1>
             <p><a href='` + metricsPath + `'>Metrics</a></p>
             <p><a href='/health'>Health (gRPC connection state to Starlink dish)</a></p>
             </body>
             </html>`))
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		state := exporter.GetState()
		var statusCode int
		switch state {
		case 0, 2:
			// Idle or Ready
			statusCode = http.StatusOK
		case 1, 3:
			// Connecting or TransientFailure
			statusCode = http.StatusServiceUnavailable
		case 4:
			// Shutdown
			statusCode = http.StatusInternalServerError
		}
		log.Infof("code: %d, state: %d", statusCode, state)
		w.WriteHeader(statusCode)
		_, _ = fmt.Fprintf(w, "%s\n", state)
	})

	http.Handle(metricsPath, promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	log.Infof("listening on %s", listen)
	log.Fatal(http.ListenAndServe(listen, nil))
}
