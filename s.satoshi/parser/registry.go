package parser

import (
	"fmt"
	"sync"
)

var (
	registry = map[string]TradeParser{}
	mu       sync.RWMutex
)

func register(identifier string, parser TradeParser) {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := registry[identifier]; ok {
		panic(fmt.Sprintf("Cannot register the same parser more than once; %s", identifier))
	}

	registry[identifier] = parser
}

func getParserByIdentifier(identifier string) (TradeParser, bool) {
	mu.RLock()
	defer mu.RUnlock()

	p, ok := registry[identifier]
	return p, ok
}
