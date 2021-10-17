package client

import (
	"context"
	"sync"
	"time"

	"github.com/monzo/slog"
)

const (
	maxRequestsPerMinute = 50
)

// RateLimiter is a bucket rate limiter that flushes every minute.
type RateLimiter struct {
	count    int
	mu       sync.RWMutex
	requests chan struct{}
}

// NewRateLimiter is a factory func that setups the rate limiter.
func NewRateLimiter(ctx context.Context) *RateLimiter {
	r := &RateLimiter{
		requests: make(chan struct{}, maxRequestsPerMinute-1),
	}

	// Start request monitor loop async.
	go r.requestMonitor(ctx)

	return r
}

// Throttle blocks until we can proceed.
func (r *RateLimiter) Throttle() {
	r.requests <- struct{}{}
}

func (r *RateLimiter) requestMonitor(ctx context.Context) {
	t := time.NewTicker(59 * time.Second)
	for {
		select {
		case <-t.C:
			for i := 0; i <= maxRequestsPerMinute; i++ {
				select {
				case <-r.requests:
				default:
					continue
				}
			}
		case <-ctx.Done():
			slog.Info(ctx, "Coingecko rate limiter gracefully shutting down; context cancelled.")
			return
		}
	}
}
