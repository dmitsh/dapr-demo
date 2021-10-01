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
	ctx := context.Background()

	var g run.Group
	// Signal handler
	g.Add(run.SignalHandler(ctx, os.Interrupt, syscall.SIGTERM))
	// Subscriber service
	{
		sub, err := pubsub.NewSubscriberService()
		if err != nil {
			log.Fatalf("ERROR: %v", err)
		}
		g.Add(
			sub.Start,
			func(err error) {
				sub.Stop()
				log.Printf("Stopped subscriber service")
			},
		)
	}
	// Prometheus service
	{
		prom := pubsub.NewPrometheusService(ctx, ":8181")
		g.Add(prom.Start, prom.Stop)
	}

	if err := g.Run(); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
}
