package consumers

import (
	"fmt"
	"sync"
)

var (
	registry = map[string]Consumer{}
	mu       sync.RWMutex
)

func register(id string, consumer Consumer) {
	if _, ok := registry[id]; ok {
		panic(fmt.Sprintf("Cannot register consumers with the same ID: %s", id))
	}

	registry[id] = consumer
}

// Registry returns the registry of all consumers registered to satoshi.
func Registry() map[string]Consumer {
	mu.RLock()
	defer mu.RUnlock()
	return registry
}
