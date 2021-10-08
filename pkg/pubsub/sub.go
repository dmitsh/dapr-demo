package pubsub

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
)

var (
	sub_total = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "sub_total",
			Help: "Total number of requests.",
		},
		[]string{"pubsub", "topic", "miss"})
)

type Subscriber struct {
	common.Service
	debug bool
}

func NewSubscriberService(cfg *Config) (common.Service, error) {
	log.Printf("NewSubscriberService address %s", cfg.appAddr)
	s := &Subscriber{
		Service: daprd.NewService(cfg.appAddr),
		debug:   cfg.debug,
	}

	for _, topic := range []string{TopicRed, TopicBlue, TopicGreen} {
		sub := &common.Subscription{
			PubsubName: cfg.pubsub,
			Topic:      topic,
			Route:      "/" + topic,
			Metadata:   map[string]string{},
		}
		if len(cfg.consumerID) != 0 {
			sub.Metadata["consumerID"] = cfg.consumerID
		}
		if err := s.AddTopicEventHandler(sub, s.eventHandler); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (s *Subscriber) eventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	if s.debug {
		log.Printf("<SUB> Event - Pubsub: %s, Topic: %s, ID: %s, Data: %s", e.PubsubName, e.Topic, e.ID, e.Data)
	}
	var miss string
	if !strings.HasPrefix(e.Data.(string), fmt.Sprintf("%s-%s-", e.PubsubName, e.Topic)) {
		log.Printf("<SUB> ERROR in topic %s : %v", e.Topic, e.Data)
		miss = e.Data.(string)
	}
	sub_total.WithLabelValues(e.PubsubName, e.Topic, miss).Inc()
	return false, nil
}
