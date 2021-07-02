package ttlcache

import (
	"sync"
	"time"
)

type cacheResult struct {
	expirationDate time.Time
	value          interface{}
}

func (ttlr *cacheResult) HasExpired() bool {
	return ttlr.expirationDate.Before(time.Now())
}

// TTLCache provides a simple cache with a time-to-live eviction policy policy.
// old expired keys aren't removed, rather replaced; so this isn't intented for uses
// with a large possible set of keys.
// It works well for a small finite set of known keys, that are both hot & liable to refresh
// in near term.
type TTLCache struct {
	ttl   time.Duration
	cache map[string]*cacheResult
	mu    sync.RWMutex
}

// TTLCacheNull ...
type TTLCacheNull struct{}

// New creates a new TTLCache, with the given time-to-live duration.
func New(ttl time.Duration) *TTLCache {
	return &TTLCache{
		cache: map[string]*cacheResult{},
		ttl:   ttl,
	}
}

// Get retreives the value for the given key and whether it is expired or not
func (tc *TTLCache) Get(key string) (interface{}, bool) {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	r, ok := tc.cache[key]
	if !ok {
		return nil, false
	}
	return r.value, r.HasExpired()
}

// Set sets the value to the given key. It is thread-safe
func (tc *TTLCache) Set(key string, value interface{}) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	r := &cacheResult{
		expirationDate: time.Now().Add(tc.ttl),
		value:          value,
	}
	tc.cache[key] = r
}

func (tc *TTLCache) Exists(key string) bool {
	v, hasExpired := tc.Get(key)
	if v == nil {
		return false
	}
	return !hasExpired
}

// SetNull is for when the caller doesn't care about the value. But wants to use the ttlcache
// to determine some key has been accessed in some time.
func (tc *TTLCache) SetNull(key string) {
	tc.Set(key, TTLCacheNull{})
}

//GetAndRefreshExpiry will get the key, if expired it will update the expiry data to now plus the default ttl
// the first return value is the value stored for that key, and the second is whether it is expired or not.
func (tc *TTLCache) GetAndRefreshExpiry(key string) (interface{}, bool) {
	v, expired := tc.Get(key)
	if !expired {
		return v, false
	}
	tc.Set(key, v)
	return v, true
}
