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
	log.Printf("Starting publisher on DAPR_GRPC_PORT %s", os.Getenv("DAPR_GRPC_PORT"))

	ctx := context.Background()

	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	var g run.Group
	// Signal handler
	g.Add(run.SignalHandler(ctx, os.Interrupt, syscall.SIGTERM))
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
