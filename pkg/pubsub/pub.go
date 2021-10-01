package pubsub

import (
	"context"
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
	[]string{"topic", "error"})

func PublishHandler(ctx context.Context, client dapr.Client, topic string) (execute func() error, interrupt func(error)) {
	ctx, cancel := context.WithCancel(ctx)
	return func() error {
			return publish(ctx, client, topic)
		},
		func(err error) {
			cancel()
			log.Printf("Stopped publisher %s", topic)
		}
}

func publish(ctx context.Context, client dapr.Client, topic string) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("<PUB> Topic %s done", topic)
			return nil
		case <-ticker.C:
			//log.Printf("<PUB> Publish %s", topic)
			var msg string
			if err := client.PublishEvent(ctx, PubSub, topic, []byte(topic)); err != nil {
				log.Printf("<PUB> Error topic %s: %v", topic, err)
				msg = err.Error()
			}
			pub_total.WithLabelValues(topic, msg).Inc()
		}
	}
}
