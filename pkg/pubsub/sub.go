package pubsub

import (
	"context"
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
		[]string{"topic", "miss"})

	subRed = &common.Subscription{
		PubsubName: PubSub,
		Topic:      TopicRed,
		Route:      "/red",
	}
	subBlue = &common.Subscription{
		PubsubName: PubSub,
		Topic:      TopicBlue,
		Route:      "/blue",
	}
	subGreen = &common.Subscription{
		PubsubName: PubSub,
		Topic:      TopicGreen,
		Route:      "/green",
	}
)

func NewSubscriberService() (common.Service, error) {
	s := daprd.NewService(":8080")

	if err := s.AddTopicEventHandler(subRed, eventHandler); err != nil {
		return nil, err
	}
	if err := s.AddTopicEventHandler(subBlue, eventHandler); err != nil {
		return nil, err
	}
	if err := s.AddTopicEventHandler(subGreen, eventHandler); err != nil {
		return nil, err
	}
	return s, nil
}

func eventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	//log.Printf("<SUB> Event - Pubsub: %s, Topic: %s, ID: %s, Data: %s", e.PubsubName, e.Topic, e.ID, e.Data)
	var miss string
	if strings.Compare(e.Topic, e.Data.(string)) != 0 {
		log.Printf("<SUB> ERROR in topic %q : %v", e.Topic, e.Data)
		miss = e.Data.(string)
	}
	sub_total.WithLabelValues(e.Topic, miss).Inc()
	return false, nil
}
