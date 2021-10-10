package parser

import (
	"fmt"
	"sync"

	"github.com/monzo/slog"
)

var (
	registry = make(map[string][]TradeParser)
	mu       sync.RWMutex
)

func register(identifier string, parser []TradeParser) {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := registry[identifier]; ok {
		panic(fmt.Sprintf("Cannot register the same parser more than once; %s", identifier))
	}

	registry[identifier] = parser

	slog.Info(nil, "Registered parser for: %s", identifier)
}

func getParsersByIdentifier(identifier string) ([]TradeParser, bool) {
	mu.RLock()
	defer mu.RUnlock()

	p, ok := registry[identifier]
	return p, ok
}
