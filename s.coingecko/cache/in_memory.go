package cache

import "swallowtail/libraries/ttlcache"

// InMemoryCache is an in memory cache that makes use of a TTL for the values assigned to given keys.
type InMemoryCache struct {
	*ttlcache.TTLCache
}
