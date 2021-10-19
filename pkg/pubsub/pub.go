package pubsub

import (
	"context"
	"fmt"
	"log"
	"time"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var pub_total = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "pub_total",
		Help: "Total number of requests.",
	},
	[]string{"pubsub", "topic", "error"})

func PublishHandler(ctx context.Context, client dapr.Client, topic string, cfg *Config) (execute func() error, interrupt func(error)) {
	ctx, cancel := context.WithCancel(ctx)
	return func() error {
			return publish(ctx, client, topic, cfg)
		},
		func(err error) {
			cancel()
			log.Printf("Stopped publisher %s", topic)
		}
}

func publish(ctx context.Context, client dapr.Client, topic string, cfg *Config) error {
	ticker := time.NewTicker(cfg.pubInterval)
	defer ticker.Stop()
	prom := len(cfg.promAddr) != 0
	var cnt int64
	for {
		select {
		case <-ctx.Done():
			log.Printf("<PUB> Topic %s done", topic)
			return nil
		case <-ticker.C:
			cnt++
			data := fmt.Sprintf("%s-%s-%s-%d", cfg.pubsub, topic, cfg.podName, cnt)
			if cfg.debug {
				log.Printf("<PUB> Publish %s", data)
			}
			var errMsg string
			if err := client.PublishEvent(ctx, cfg.pubsub, topic, []byte(data)); err != nil {
				log.Printf("<PUB> Error topic %s: %v", topic, err)
				errMsg = err.Error()
			}
			if prom {
				pub_total.WithLabelValues(cfg.pubsub, topic, errMsg).Inc()
			}
		}
	}
}
