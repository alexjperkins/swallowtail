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

func TestTTLCache_SetNullAndExist(t *testing.T) {
	t.Parallel()
	key := "test-key"

	tc := New(time.Hour)
	exists := tc.Exists(key)
	assert.False(t, exists)

	tc.SetNull(key)
	exists = tc.Exists(key)
	assert.True(t, exists)

	tc = New(100 * time.Millisecond)
	tc.SetNull(key)

	time.Sleep(110 * time.Millisecond)

	exists = tc.Exists(key)
	assert.False(t, exists)
}
