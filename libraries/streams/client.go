package streams

import "sync"

var (
	defaultClient   StreamsClient
	defaultClientMu sync.Mutex
)

// StreamsClient ...
type StreamsClient interface{}

func Client() StreamsClient {
	defaultClientMu.Lock()
	defer defaultClientMu.Unlock()

	return defaultClient
}
