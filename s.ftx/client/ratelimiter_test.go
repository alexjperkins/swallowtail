package client

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFTXRatelimiter(t *testing.T) {
	if !testing.Short() {
		t.Skip("Skipping rate limiter test due to length taken")
	}

	t.Parallel()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	t.Cleanup(func() {
		cancel()
	})

	rl := newFTXRateLimiter(ctx)

	// Assert wait negliable if number of requests doesn't exceed the limit.
	then := time.Now()
	rl.Wait()
	elasped := time.Since(then)

	assert.True(t, elasped < (periodicity/maxRequestsPerPeriod-1)) // minus one to give us a buffer.

	// Sleep small period to refresh bucket tokens.
	time.Sleep(periodicity / maxRequestsPerPeriod)

	// Assert that we do actually get rate limited if number of requests does exceed limit.
	then = time.Now()
	for i := 0; i < maxRequestsPerPeriod+1; i++ {
		rl.Wait()
	}
	elasped = time.Since(then)

	assert.True(t, elasped > periodicity/maxRequestsPerPeriod)
}
