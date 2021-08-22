package cache

import "swallowtail/libraries/ttlcache"

// InMemoryCache is an in memory cache that makes use of a TTL for the values assigned to given keys.
type InMemoryCache struct {
	*ttlcache.TTLCache
}

// Get ...
func (i *InMemoryCache) Get(key string) (interface{}, bool, error) {
	v, ok := i.TTLCache.Get(key)
	return v, ok, nil
}

// Set ...
func (i *InMemoryCache) Set(key string, value interface{}) error {
	i.TTLCache.Set(key, value)
	return nil
}
