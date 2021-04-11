package ttlcache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTTLCache_GetAndSet(t *testing.T) {
	t.Parallel()

	key := "test-key"
	initialInputValue := 1

	// ttlcache with long ttl
	tc := New(time.Hour)
	tc.Set(key, initialInputValue)

	v, expired := tc.Get(key)
	assert.False(t, expired)

	vf, _ := v.(int)
	assert.Equal(t, initialInputValue, vf)

	// ttlcache with short ttl
	ttl := time.Duration(time.Millisecond)
	tc = New(ttl)
	tc.Set(key, initialInputValue)

	time.Sleep(time.Millisecond * 100)

	v, expired = tc.Get(key)
	assert.True(t, expired)
}
