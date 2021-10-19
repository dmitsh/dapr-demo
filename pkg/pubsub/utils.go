package pubsub

import (
	"flag"
	"os"
	"time"
)

const (
	TopicRed   = "red"
	TopicBlue  = "blue"
	TopicGreen = "green"
)

type Config struct {
	pubsub      string
	podName     string
	appAddr     string
	promAddr    string
	pubInterval time.Duration
	debug       bool
}

func ProcessCommandLine() *Config {
	cfg := &Config{
		podName: os.Getenv("POD_NAME"),
	}
	flag.StringVar(&cfg.pubsub, "p", "pubsub", "Dapr pubsub component name")
	flag.StringVar(&cfg.appAddr, "a", ":6100", "application service address")
	flag.StringVar(&cfg.promAddr, "m", "", "prometheus service address")
	flag.DurationVar(&cfg.pubInterval, "t", time.Second, "publishing time interval")
	flag.BoolVar(&cfg.debug, "d", false, "debug flag")
	flag.Parse()
	return cfg
}

/*
func WaitForDapr(ctx context.Context, cfg *Config) error {
	if len(cfg.healthz) == 0 {
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			resp, err := http.Get(cfg.healthz)
			if err == nil {
				switch resp.StatusCode {
				case http.StatusOK, http.StatusNoContent:
					fmt.Printf("HEALTHZ OK: %#v\n", resp)
					return nil
				}
				fmt.Printf("HEALTHZ BAD STATUS: %#v\n", resp)
			}
			fmt.Printf("HEALTHZ ERR: %v\n", err)
		}
	}
}
*/
