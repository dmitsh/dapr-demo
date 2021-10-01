package main

import (
	"context"
	"log"
	"os"
	"syscall"
	"time"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/dmitsh/dapr-demo/pkg/pubsub"
	"github.com/oklog/run"
)

var (
	stop bool
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
	{
		// Signal handler
		g.Add(run.SignalHandler(ctx, os.Interrupt, syscall.SIGTERM))
		// Publish red
		g.Add(
			func() error {
				return publish(ctx, client, pubsub.TopicRed)
			},
			func(err error) {
				log.Printf("Stopping RED")
				stop = true
			},
		)
		// Publish blue
		g.Add(
			func() error {
				return publish(ctx, client, pubsub.TopicBlue)
			},
			func(err error) {
				log.Printf("Stopping BLUE")
				stop = true
			},
		)
	}

	if err := g.Run(); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
}

func publish(ctx context.Context, client dapr.Client, topic string) error {
	for !stop {
		log.Printf("Publish %s", topic)
		if err := client.PublishEvent(ctx, pubsub.PubSub, topic, []byte(topic)); err != nil {
			log.Printf("Error topic %s: %v", topic, err)
			return err
		}
		time.Sleep(3 * time.Second)
	}
	return nil
}
