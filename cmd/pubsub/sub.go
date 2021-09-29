package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
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
	s := daprd.NewService(":8080")

	if err := s.AddTopicEventHandler(subRed, eventHandler); err != nil {
		log.Fatalf("error adding topic subscription: %v", err)
	}

	if err := s.AddTopicEventHandler(subBlue, eventHandler); err != nil {
		log.Fatalf("error adding topic subscription: %v", err)
	}

	var g run.Group
	{
		// Termination handler.
		term := make(chan os.Signal, 1)
		signal.Notify(term, os.Interrupt, syscall.SIGTERM)
		cancel := make(chan struct{})
		g.Add(
			func() error {
				select {
				case <-term:
					fmt.Println("Received SIGTERM, exiting gracefully...")
				case <-cancel:
					fmt.Println("Canceled")
				}
				return nil
			},
			func(err error) {
				fmt.Println("Stopping term")
				close(cancel)
			},
		)
		// Application service
		g.Add(
			func() error {
				return s.Start()
			},
			func(err error) {
				fmt.Println("Stopping application service")
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
				fmt.Println("Stopping Prometheus")
			},
		)
	}

	if err := g.Run(); err != nil {
		fmt.Printf("ERROR: %v\n", err)
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
