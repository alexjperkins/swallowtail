package streams

import (
	"context"
	"fmt"
	"swallowtail/libraries/gerrors"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/monzo/slog"
)

type kafkaProducer struct {
	sarama.SyncProducer
}

type kafkaSubscription interface {
	connect(brokers []string, config *clientConfig)
}

type KafkaClient struct {
	subscriptionsMu sync.Mutex
	subscriptions   map[string]kafkaSubscription

	producersMu sync.Mutex
	producers   map[string]kafkaProducer

	clientConfig            *clientConfig
	subscriptionsConfigHash uint32
}

func NewKafkaClient(opts ...func(*KafkaClient) *KafkaClient) {
	sarama.PanicHandler = func(v interface{}) {
		var err error

		switch v := v.(type) {
		case error:
			err = gerrors.Augment(err, "streams client: panic", nil)
		default:
			err = fmt.Errorf("streams client: panic: %+v", v)
		}

		slog.Critical(context.Background(), "streams recovered panic; %+v", err)
	}
}
