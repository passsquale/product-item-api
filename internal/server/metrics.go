package server

import (
	"fmt"
	"github.com/passsquale/product-item-api/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func createMetricsServer(cfg *config.Config) *http.Server {
	addr := fmt.Sprintf("%s:%d", cfg.Metrics.Host, cfg.Metrics.Port)

	mux := http.DefaultServeMux
	mux.Handle(cfg.Metrics.Path, promhttp.Handler())

	metricsServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return metricsServer
}
