package pubsub

import (
	"context"
	"fmt"
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
	if cfg.promPort == 0 {
		log.Printf("Skip Prometheus")
		return nil
	}
	log.Printf("NewPrometheusService port %d", cfg.promPort)
	return &PrometheusSvc{
		ctx: ctx,
		srv: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.promPort),
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
