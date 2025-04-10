package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Create a non-global registry.
	reg := prometheus.NewRegistry()

	// Add collectors for Go runtime stats and process stats.
	// TODO: check if this can be done more easily
	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	// Expose metrics using the custom registry.
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
