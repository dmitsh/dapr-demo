package pubsub

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Config struct {
	pubsub      string
	topics      []string
	appPort     int
	promPort    int
	pubInterval time.Duration
	podName     string
	debug       bool
}

func (c *Config) Topics() []string {
	return c.topics
}

func ProcessCommandLine() (*Config, error) {
	cfg := &Config{
		podName: os.Getenv("POD_NAME"),
		topics:  []string{},
	}
	a := kingpin.New(filepath.Base(os.Args[0]), "PubSub test app")
	a.HelpFlag.Short('h')
	a.Flag("pubsub", "Dapr pubsub component name.").Short('p').Default("pubsub").StringVar(&cfg.pubsub)
	a.Flag("app.port", "application service address.").Short('a').Default("6100").IntVar(&cfg.appPort)
	a.Flag("prom.port", "Prometheus service address.").Short('m').Default("0").IntVar(&cfg.promPort)
	a.Flag("pub.interval", "publishing time interval.").Short('i').Default("1s").DurationVar(&cfg.pubInterval)
	a.Flag("topic", "topic name (repetitive)").Short('t').Required().StringsVar(&cfg.topics)
	a.Flag("debug", "debug flag.").Short('d').Default("false").BoolVar(&cfg.debug)

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		a.Usage(os.Args[1:])
		return nil, errors.Wrapf(err, "Error parsing commandline arguments")
	}

	return cfg, nil
}

func Exit(err error) {
	log.Printf("Error: %s", err.Error())
	os.Exit(1)
}

func WaitForDapr(ctx context.Context, port string) error {
	url := fmt.Sprintf("http://localhost:%s/v1.0/metadata", port)
	// wait for up to 30 seconds
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Minute))
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Printf("probing %s\n", url)
			if resp, err := http.Get(url); err != nil {
				fmt.Printf("dapr probing error: %v\n", err)
			} else {
				fmt.Printf("dapr probing status %s (%d)\n", resp.Status, resp.StatusCode)
				if resp.StatusCode == http.StatusOK {
					return nil
				}
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
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
