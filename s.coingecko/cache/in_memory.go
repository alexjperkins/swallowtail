package cache

import "swallowtail/libraries/ttlcache"

// InMemoryCache is an in memory cache that makes use of a TTL for the values assigned to given keys.
type InMemoryCache struct {
	ttl *ttlcache.TTLCache
}

// Get ...
func (i *InMemoryCache) Get(key string) (interface{}, bool, error) {
	v, ok := i.ttl.Get(key)
	return v, ok, nil
}

// Set ...
func (i *InMemoryCache) Set(key string, value interface{}) error {
	i.ttl.Set(key, value)
	return nil
}
