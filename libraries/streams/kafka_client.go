package streams

import (
	"context"
	"fmt"
	"swallowtail/libraries/gerrors"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/monzo/slog"
	"google.golang.org/grpc"

	streamsconsumerproto "swallowtail/s.streams-consumer/proto"
)

type kafkaProducer struct {
	sarama.SyncProducer
	wg sync.WaitGroup
}

type kafkaSubscription interface {
	// connect is called to start a subscription & contains message consuming logic.
	connect(brokers []string, config *clientConfig) error

	// reconnect is called to restart a connection whenever a configuration is changed.
	// it should be equivalent to `stop()` & then `connect()`.
	reconnect(brokers []string, config *clientConfig) error

	// stop gracefully shutdowns a subscription.
	stop()

	// Topic returns the topic for which the subscription is for.
	Topic() string

	// Group returns the consumer group ID the subscription is for.
	Group() string
}

// KafkaClient is a streams client for Kafka: a distributed commit log.
type KafkaClient struct {
	subscriptionsMu sync.Mutex
	subscriptions   map[string]kafkaSubscription

	producersMu sync.Mutex
	producers   map[string]kafkaProducer

	clientConfig *clientConfig

	closed chan struct{}

	streamsConsumerClient streamsconsumerproto.StreamsconsumerClient
}

// WithStreamsConsumerClient returns a closure over a Kafka Client in order to set the internal Streams Consumer Client to the one passed.
func WithStreamsConsumerClient(streamsConsumerClient streamsconsumerproto.StreamsconsumerClient) func(*KafkaClient) {
	return func(client *KafkaClient) {
		client.streamsConsumerClient = streamsConsumerClient
	}
}

// NewKafkaClient is factory method that generates a new Kafka Client.
func NewKafkaClient(opts ...func(*KafkaClient)) *KafkaClient {
	// Set internal sarama panic handler.
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

	subscriptions := map[string]kafkaSubscription{}
	k := &KafkaClient{
		closed:                make(chan struct{}),
		clientConfig:          newClientConfig(),
		subscriptions:         subscriptions,
		streamsConsumerClient: streamsconsumerproto.NewStreamsconsumerClient(&grpc.ClientConn{}),
	}

	// Apply options to Kafka client.
	for _, o := range opts {
		o(k)
	}

	return k
}

// Subscribe ...
func (k *KafkaClient) Subscribe(ctx context.Context, topic, group string, handler Handler) error {
	k.subscriptionsMu.Lock()
	defer k.subscriptionsMu.Unlock()

	select {
	case <-k.closed:
		err := gerrors.InternalService("kafka_client_closed", nil)
		slog.Critical(ctx, "Error subscribing to topic: %s %v", topic, err)
		return err
	default:
	}

	// Validate subscription.
	validation := validateTopicName(
		func(topic, group string) error {
			if len(topic) > 249 {
				return gerrors.InternalService(
					"illegal kafka topic name: length cannot exceed than 249 chars", nil,
				)
			}
			return nil
		},
	)

	if err := validation(topic, group); err != nil {
		errParams := map[string]string{
			"topic": topic,
			"group": group,
		}
		slog.Critical(ctx, "Failed to subscribe to topic: %v %+v", err, errParams)
		return gerrors.Augment(err, "kafka client failed validation", errParams)
	}

	return nil
}
