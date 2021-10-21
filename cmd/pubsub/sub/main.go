package main

import (
	"context"
	"log"
	"os"
	"syscall"

	"github.com/dmitsh/dapr-demo/pkg/pubsub"
	"github.com/oklog/run"
)

var ()

func main() {
	log.Printf("Starting subscriber on DAPR_GRPC_PORT %s", os.Getenv("DAPR_GRPC_PORT"))

	cfg, err := pubsub.ProcessCommandLine()
	if err != nil {
		pubsub.Exit(err)
	}

	sub, err := pubsub.NewSubscriberService(cfg)
	if err != nil {
		pubsub.Exit(err)
	}

	ctx := context.Background()
	var g run.Group
	// Signal handler
	g.Add(run.SignalHandler(ctx, os.Interrupt, syscall.SIGTERM))
	// Subscriber service
	g.Add(
		sub.Start,
		func(err error) {
			sub.Stop()
			log.Printf("Stopped subscriber service")
		},
	)
	// Prometheus service
	if prom := pubsub.NewPrometheusService(ctx, cfg); prom != nil {
		g.Add(prom.Start, prom.Stop)
	}

	if err := g.Run(); err != nil {
		pubsub.Exit(err)
	}
}
