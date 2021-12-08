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
	grpcPort := os.Getenv("DAPR_GRPC_PORT")
	httpPort := os.Getenv("DAPR_HTTP_PORT")
	log.Printf("Starting publisher on gRPC:%s / http:%s", grpcPort, httpPort)

	cfg, err := pubsub.ProcessCommandLine()
	if err != nil {
		pubsub.Exit(err)
	}

	client, err := dapr.NewClient()
	if err != nil {
		pubsub.Exit(err)
	}
	defer client.Close()

	ctx := context.Background()

	if err = pubsub.WaitForDapr(ctx, httpPort); err != nil {
		pubsub.Exit(err)
	}

	var g run.Group
	// Signal handler
	g.Add(run.SignalHandler(ctx, os.Interrupt, syscall.SIGTERM))
	// Prometheus service
	if prom := pubsub.NewPrometheusService(ctx, cfg); prom != nil {
		g.Add(prom.Start, prom.Stop)
	}
	// Publish topics
	for _, topic := range cfg.Topics() {
		g.Add(pubsub.PublishHandler(ctx, client, topic, cfg))
	}

	if err := g.Run(); err != nil {
		pubsub.Exit(err)
	}
}
