package sync

import (
	"fmt"
	"sync"
)

var (
	registry = map[string]Syncer{}
	mu       sync.RWMutex
)

func register(id string, syncer Syncer) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := registry[id]; ok {
		panic(fmt.Sprintf("Cannot register the same syncer twice: %v", id))
	}

	registry[id] = syncer
}
