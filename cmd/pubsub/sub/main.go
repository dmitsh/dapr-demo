package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"syscall"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
	"github.com/dmitsh/dapr-demo/pkg/pubsub"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	counter_total = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pubsub_total",
			Help: "Total number of requests.",
		},
		[]string{"topic"})

	counter_error = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "pubsub_errors",
			Help: "Number of errors.",
		},
		[]string{"topic"})

	subRed = &common.Subscription{
		PubsubName: pubsub.PubSub,
		Topic:      pubsub.TopicRed,
		Route:      "/red",
	}
	subBlue = &common.Subscription{
		PubsubName: pubsub.PubSub,
		Topic:      pubsub.TopicBlue,
		Route:      "/blue",
	}
)

func main() {
	log.Printf("Starting subscriber on DAPR_GRPC_PORT %s", os.Getenv("DAPR_GRPC_PORT"))

	ctx := context.Background()

	s := daprd.NewService(":8080")

	if err := s.AddTopicEventHandler(subRed, eventHandler); err != nil {
		log.Fatalf("error adding topic subscription: %v", err)
	}

	if err := s.AddTopicEventHandler(subBlue, eventHandler); err != nil {
		log.Fatalf("error adding topic subscription: %v", err)
	}

	var g run.Group
	{
		// Signal handler
		g.Add(run.SignalHandler(ctx, os.Interrupt, syscall.SIGTERM))
		// Application service
		g.Add(
			func() error {
				return s.Start()
			},
			func(err error) {
				log.Printf("Stopping application service")
				s.Stop()
			},
		)
		// Prometheus
		g.Add(
			func() error {
				http.Handle("/metrics", promhttp.Handler())
				return http.ListenAndServe(":8181", nil)
			},
			func(err error) {
				log.Printf("Stopping Prometheus")
			},
		)
	}

	if err := g.Run(); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
}

func eventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	log.Printf("event - Pubsub: %s, Topic: %s, ID: %s, Data: %s", e.PubsubName, e.Topic, e.ID, e.Data)

	if strings.Compare(e.Topic, e.Data.(string)) != 0 {
		log.Printf("ERROR in topic %q : %v", e.Topic, e.Data)
		counter_error.WithLabelValues(e.Topic).Inc()
	}
	counter_total.WithLabelValues(e.Topic).Inc()

	return false, nil
}
