package cache

import (
	"fmt"
	"time"
)

// ErrCacheUnsupported temp error whilst we migrate to go 1.19 (with generics).
var ErrCacheUnsupported = fmt.Errorf("unsupported cache")

// InMemoryCache is an in memory cache that makes use of a TTL for the values assigned to given keys.
type InMemoryCache struct {
}

// NewInMemoryCache is a factory method to instantiate & setup an in memory cache that satisifies the coingecko cache interface.
func NewInMemoryCache(_ time.Duration, _ func(key string) (interface{}, time.Duration, error)) *InMemoryCache {
	c := &InMemoryCache{}
	return c
}

// Get ...
func (i *InMemoryCache) Get(key string) (interface{}, error) {
	return nil, fmt.Errorf("set: %w", ErrCacheUnsupported)
}

// Set ...
func (i *InMemoryCache) Set(key string, value interface{}) error {
	return fmt.Errorf("set: %w", ErrCacheUnsupported)
}

// Close ...
func (i *InMemoryCache) Close() {
}
