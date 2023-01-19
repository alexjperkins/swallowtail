package cache

import (
	"swallowtail/libraries/gerrors"
	"time"
)

// InMemoryCache is an in memory cache that makes use of a TTL for the values assigned to given keys.
type InMemoryCache struct{}

// NewInMemoryCache is a factory method to instantiate & setup an in memory cache that satisifies the coingecko cache interface.
func NewInMemoryCache(ttl time.Duration, loader func(key string) (interface{}, time.Duration, error)) *InMemoryCache {
	return &InMemoryCache{}
}

// Get ...
func (i *InMemoryCache) Get(key string) (interface{}, error) {
	return nil, gerrors.Unimplemented("get inmemory cache", nil)
}

// Set ...
func (i *InMemoryCache) Set(key string, value interface{}) error {
	return gerrors.Unimplemented("setinmemory cache", nil)
}

// Close ...
func (i *InMemoryCache) Close() {
}
