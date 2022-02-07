package streams

import (
	"context"
	"sync"
)

var (
	defaultClient   StreamsClient
	defaultClientMu sync.Mutex
)

// Handler defines the interface of an ordered event handler: it accepts an event & returns a result.
type Handler func(event Event) Result

// StreamsClient ...
type StreamsClient interface {
	Subscribe(ctx context.Context, topic, group string)
}

func Client() StreamsClient {
	defaultClientMu.Lock()
	defer defaultClientMu.Unlock()

	return defaultClient
}
