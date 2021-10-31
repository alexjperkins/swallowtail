package ratelimit

import (
	"context"
	"sync"
	"time"

	"github.com/monzo/slog"
)

// NewLinearBackpressureRateLimiter creates, starts pressure relief goroutine & returns back the caller
// a `NewLinearBackpressureRateLimiter` to manage rate limits sysmatically & linearly.
func NewLinearBackpressureRateLimiter(ctx context.Context, periodicity time.Duration, maxRequestPerPeriod int) *LinearBackpressureRateLimiter {
	l := &LinearBackpressureRateLimiter{
		periodicity:          periodicity,
		maxRequestsPerPeriod: maxRequestPerPeriod,
		requests:             make(chan struct{}, maxRequestPerPeriod),
	}

	// Start pressure relief goroutine.
	go l.pressureRelief(ctx)

	return l
}

// LinearBackpressureRateLimiter provides a mechanism of rate limiting via backpressure.
// It doesn't allow a callee to throttle to progress if the max number of requests have already been
// reached in the periodicity provided.
// Pressure is reliefed completely after every period of time, as defined by the periodicity.
type LinearBackpressureRateLimiter struct {
	periodicity          time.Duration
	maxRequestsPerPeriod int

	count    int
	mu       sync.RWMutex
	requests chan struct{}
}

// Throttle ...
func (r *LinearBackpressureRateLimiter) Throttle() {
	r.requests <- struct{}{}
}

// ThrottleWithOptions ...
func (r *LinearBackpressureRateLimiter) ThrottleWithOptions(_ *RateLimiterOpts) {
	r.requests <- struct{}{}
}

func (r *LinearBackpressureRateLimiter) pressureRelief(ctx context.Context) {
	t := time.NewTicker(r.periodicity)
	for {
		select {
		case <-t.C:
			for i := 0; i <= r.maxRequestsPerPeriod; i++ {
				select {
				case <-r.requests:
				default:
					continue
				}
			}
		case <-ctx.Done():
			slog.Info(ctx, "Rate limiter gracefully shutting down; context cancelled.")
			return
		}
	}
}
