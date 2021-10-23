package cache

import (
	"time"

	"github.com/ReneKroon/ttlcache/v2"
)

// InMemoryCache is an in memory cache that makes use of a TTL for the values assigned to given keys.
type InMemoryCache struct {
	c *ttlcache.Cache
}

// NewInMemoryCache is a factory method to instantiate & setup an in memory cache that satisifies the coingecko cache interface.
func NewInMemoryCache(ttl time.Duration, loader func(key string) (interface{}, time.Duration, error)) *InMemoryCache {
	c := &InMemoryCache{
		c: ttlcache.NewCache(),
	}

	c.c.SetCacheSizeLimit(1024)
	c.c.SetLoaderFunction(loader)
	c.c.SetTTL(ttl)
	c.c.SkipTTLExtensionOnHit(true)

	return c
}

// Get ...
func (i *InMemoryCache) Get(key string) (interface{}, error) {
	return i.c.Get(key)
}

// Set ...
func (i *InMemoryCache) Set(key string, value interface{}) error {
	return i.c.Set(key, value)
}

// Close ...
func (i *InMemoryCache) Close() {
	i.Close()
}
