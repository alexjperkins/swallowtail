package client

import (
	"context"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/monzo/slog"
)

const (
	maxRequestsPerMinute = 50
)

type RateLimiter struct {
	count    int
	mu       sync.RWMutex
	requests chan struct{}
}

func NewRateLimiter(ctx context.Context) *RateLimiter {
	r := &RateLimiter{
		requests: make(chan struct{}, 1),
	}
	go r.requestMonitor(ctx)
	return r
}

// Throttle blocks until we can proceed.
func (r *RateLimiter) Throttle() {
	// Register the request first.
	r.registerRequest()

	boff := backoff.NewExponentialBackOff()
	for {
		shouldThrottle := func() bool {
			r.mu.RLock()
			defer r.mu.RUnlock()
			return r.count >= maxRequestsPerMinute
		}()
		if !shouldThrottle {
			return
		}

		<-time.After(boff.NextBackOff())
	}
}

// RegisterRequest adds a request to the internal count. Nonblocking.
func (r *RateLimiter) registerRequest() {
	select {
	case r.requests <- struct{}{}:
	default:
	}
}

func (r *RateLimiter) requestMonitor(ctx context.Context) {
	t := time.NewTicker(59 * time.Second)
	for {
		select {
		case <-t.C:
			r.mu.Lock()
			defer r.mu.Unlock()
			r.count = 0
		case <-r.requests:
			r.mu.Lock()
			defer r.mu.Unlock()
			r.count++
		case <-ctx.Done():
			slog.Info(ctx, "Coingecko rate limiter gracefully shutting down; context cancelled.")
			return
		}
	}
}
