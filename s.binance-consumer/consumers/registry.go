package consumers

import (
	"swallowtail/libraries/multiplexing"
	"sync"
)

var (
	mu        = sync.RWMutex{}
	consumers = map[string][]*multiplexing.MultiplexConsumer{}
	bufSize   = 16
)

func register(endpoint string, consumer chan interface{}, filter multiplexing.MuliplexFilter, metadata map[string]string) {
	mu.Lock()
	defer mu.Unlock()
	c := multiplexing.NewMultiplexConsumer(bufSize, filter, nil)
	consumers[endpoint] = append(consumers[endpoint], c)
}

func GetAllBinanceConsumers() map[string][]*multiplexing.MultiplexConsumer {
	return consumers
}
