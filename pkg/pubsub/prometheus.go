package pubsub

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type PrometheusSvc struct {
	ctx context.Context
	srv *http.Server
}

func NewPrometheusService(ctx context.Context, cfg *Config) *PrometheusSvc {
	log.Printf("NewPrometheusService address %s", cfg.promAddr)
	return &PrometheusSvc{
		ctx: ctx,
		srv: &http.Server{
			Addr:    cfg.promAddr,
			Handler: promhttp.Handler(),
		},
	}
}

func (h *PrometheusSvc) Start() error {
	return h.srv.ListenAndServe()
}

func (h *PrometheusSvc) Stop(error) {
	ctx, cancel := context.WithTimeout(h.ctx, 5*time.Second)
	defer cancel()
	if err := h.srv.Shutdown(ctx); err != nil {
		log.Printf("Error in stopping Prometheus service: %v", err)
	}
	log.Printf("Stopped Prometheus")
}
