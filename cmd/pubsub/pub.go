package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
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
	ctx := context.Background()

	log.Printf("Starting publisher. DAPR_GRPC_PORT=%s", os.Getenv("DAPR_GRPC_PORT"))
	client, err := dapr.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

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
					stop = true
				case <-cancel:
					fmt.Println("Canceled")
				}
				return nil
			},
			func(err error) {
				fmt.Println("Stopping term")
				stop = true
				close(cancel)
			},
		)
		// Publish red
		g.Add(
			func() error {
				return publish(ctx, client, pubsub.TopicRed)
			},
			func(err error) {
				fmt.Println("Stopping RED")
				stop = true
			},
		)
		// Publish blue
		g.Add(
			func() error {
				return publish(ctx, client, pubsub.TopicBlue)
			},
			func(err error) {
				fmt.Println("Stopping BLUE")
				stop = true
			},
		)
	}

	if err := g.Run(); err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	fmt.Println("main exit")
}

func publish(ctx context.Context, client dapr.Client, topic string) error {
	for !stop {
		fmt.Println("Publish", topic)
		if err := client.PublishEvent(ctx, pubsub.PubSub, topic, []byte(topic)); err != nil {
			fmt.Printf("Error topic %s: %v\n", topic, err)
			return err
		}
		time.Sleep(3 * time.Second)
	}
	return nil
}
