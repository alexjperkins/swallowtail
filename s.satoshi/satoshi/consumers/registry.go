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

// List lists all consumers registered to satoshi.
func List() []Consumer {
	mu.RLock()
	defer mu.RUnlock()

	commands := []Consumer{}
	for _, c := range registry {
		commands = append(commands, c)
	}

	return commands
}
