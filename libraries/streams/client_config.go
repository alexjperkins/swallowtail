package streams

import (
	"fmt"
	"hash/fnv"
	"io"
	"strconv"
	"strings"
	"swallowtail/libraries/util"
	"time"

	"github.com/Shopify/sarama"
)

type clientConfig struct {
	// Brokers.
	brokers []string

	// Swallowtail-specific configuration.
	consumerGroupSubscribeWithConcurrentSubscription bool
	debugLoggingEnabled                              bool

	// Sarama specific configuration.
	producerConfiguration      *sarama.Config
	consumerGroupConfiguration *sarama.Config

	consumePanicFunc func(err error)
}

func newClientConfig() *clientConfig {
	return &clientConfig{
		brokers: make([]string, 0),
		consumePanicFunc: func(err error) {
			panic(err)
		},
	}
}

// Hash generates a hash of the configuration, useful to know if config that
// requires a client to restart has changed
func (c *clientConfig) Hash() uint32 {
	hasher := fnv.New32a()
	io.WriteString(hasher, strings.Join(c.brokers, ","))

	// Hash production configuration.
	if c.producerConfiguration != nil {
		io.WriteString(hasher, fmt.Sprintf("%d", c.producerConfiguration.Producer.RequiredAcks))
		io.WriteString(hasher, fmt.Sprintf("%d", c.producerConfiguration.Producer.Compression))
		io.WriteString(hasher, fmt.Sprintf("%d", c.producerConfiguration.Producer.Flush.Frequency))
		io.WriteString(hasher, fmt.Sprintf("%d", c.producerConfiguration.Producer.Flush.MaxMessages))
		io.WriteString(hasher, fmt.Sprintf("%d", c.producerConfiguration.Producer.Flush.Messages))
	}

	// Hash consumerGroup configuration.
	if c.consumerGroupConfiguration != nil {
		io.WriteString(hasher, fmt.Sprintf("%d", c.consumerGroupConfiguration.Consumer.Offsets.Initial))
		io.WriteString(hasher, fmt.Sprintf("%d", c.consumerGroupConfiguration.Consumer.Offsets.Retention))
		io.WriteString(hasher, fmt.Sprintf("%d", c.consumerGroupConfiguration.Consumer.Group.Session.Timeout))
		io.WriteString(hasher, fmt.Sprintf("%d", c.consumerGroupConfiguration.Consumer.Group.Heartbeat.Interval))
		io.WriteString(hasher, fmt.Sprintf("%d", c.consumerGroupConfiguration.Consumer.Group.Rebalance.Strategy))
		io.WriteString(hasher, fmt.Sprintf("%v", c.consumerGroupConfiguration.Consumer.Group.Rebalance.Timeout))
		io.WriteString(hasher, fmt.Sprintf("%v", c.consumerGroupConfiguration.Consumer.Group.Rebalance.Retry.Max))
		io.WriteString(hasher, fmt.Sprintf("%v", c.consumerGroupConfiguration.Consumer.Group.Rebalance.Retry.Backoff))
		io.WriteString(hasher, fmt.Sprintf("%d", c.consumerGroupConfiguration.Consumer.MaxProcessingTime))
	}

	// Hash the non-sarama consumer configuration.
	io.WriteString(hasher, strconv.FormatBool(c.consumerGroupSubscribeWithConcurrentSubscription))
	io.WriteString(hasher, strconv.FormatBool(c.debugLoggingEnabled))

	return hasher.Sum32()
}

// String implements the Stringer interface.
func (c *clientConfig) String() string {
	return strings.Join(c.brokers, ",")
}

// LoadConfig loads the config with sensible defaults. In the future when we have a dynamic configuration
// designed we shall be loading from there instead.
func LoadConfig() *clientConfig {
	c := newClientConfig()

	// TODO: retrieve brokers from DB (at least until we have dynamic configuration)
	brokers := brokersFromEnvironment()
	if len(brokers) == 0 {
		brokers = append(brokers, "127.0.0.1:9092")
	}

	// Sensible defaults for producer config.
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.RequiredAcks = sarama.WaitForLocal
	producerConfig.Producer.Compression = sarama.CompressionNone
	producerConfig.Producer.Flush.Frequency = 200 * time.Millisecond
	producerConfig.Producer.Flush.MaxMessages = 0
	producerConfig.Producer.Flush.Messages = 0

	// Sensible defaults for consumer config.
	consumerConfig := sarama.NewConfig()
	consumerConfig.Version = sarama.V1_1_1_0
	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	consumerConfig.Consumer.Offsets.Retention = 7 * 24 * time.Hour
	consumerConfig.Consumer.Return.Errors = true
	consumerConfig.Consumer.Group.Session.Timeout = 30 * time.Second
	consumerConfig.Consumer.Group.Heartbeat.Interval = 10 * time.Second
	consumerConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	consumerConfig.Consumer.Group.Rebalance.Timeout = 1 * time.Minute
	consumerConfig.Consumer.Group.Rebalance.Retry.Max = 4
	consumerConfig.Consumer.MaxProcessingTime = 100 * time.Millisecond

	// Internal configuration.
	c.consumerGroupSubscribeWithConcurrentSubscription = false
	c.debugLoggingEnabled = false

	// Setting this to true means we block when consuming async on the Successes & Error chans.
	producerConfig.Producer.Return.Successes = false

	c.producerConfiguration = producerConfig
	c.consumerGroupConfiguration = consumerConfig

	c.brokers = brokers

	return c
}

func brokersFromEnvironment() []string {
	h := util.EnvWithDefault("SWALLOWTAIL_KAFKA_TCP_ADDRESSES", "")
	if h == "" {
		return nil
	}

	return strings.Split(h, ",")
}
