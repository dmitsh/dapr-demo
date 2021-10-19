package main

import (
	"context"
	"log"
	"os"
	"syscall"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dmitsh/dapr-demo/pkg/pubsub"
	"github.com/oklog/run"
)

func main() {
	log.Printf("Starting publisher/subscriber on DAPR_GRPC_PORT %s", os.Getenv("DAPR_GRPC_PORT"))
	cfg := pubsub.ProcessCommandLine()
	ctx := context.Background()

	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	sub, err := pubsub.NewSubscriberService(cfg)
	if err != nil {
		log.Fatalf("ERROR: %v", err)
	}

	var g run.Group
	// Signal handler
	g.Add(run.SignalHandler(ctx, os.Interrupt, syscall.SIGTERM))
	// Subscriber service
	g.Add(sub.Start,
		func(err error) {
			sub.Stop()
			log.Printf("Stopped subscriber service")
		},
	)
	// Prometheus service
	if prom := pubsub.NewPrometheusService(ctx, cfg); prom != nil {
		g.Add(prom.Start, prom.Stop)
	}
	// Publish red
	g.Add(pubsub.PublishHandler(ctx, client, pubsub.TopicRed, cfg))
	// Publish blue
	g.Add(pubsub.PublishHandler(ctx, client, pubsub.TopicBlue, cfg))
	// Publish green
	g.Add(pubsub.PublishHandler(ctx, client, pubsub.TopicGreen, cfg))

	if err := g.Run(); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
}
