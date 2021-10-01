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
	ctx := context.Background()

	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

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
			func() error {
				return sub.Start()
			},
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
	// Publish red
	g.Add(pubsub.PublishHandler(ctx, client, pubsub.TopicRed))
	// Publish blue
	g.Add(pubsub.PublishHandler(ctx, client, pubsub.TopicBlue))
	// Publish green
	g.Add(pubsub.PublishHandler(ctx, client, pubsub.TopicGreen))

	if err := g.Run(); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
}
