package client

import (
	"context"
	"net/http"
	"time"

	"github.com/monzo/slog"
)

const (
	periodicity          = 1 * time.Second
	maxRequestsPerPeriod = 5
)

// NewFTXRateLimiter is a factory method for a FTX rate limiter.
func newFTXRateLimiter(ctx context.Context) *ftxRateLimiter {
	f := &ftxRateLimiter{
		bucket: make(chan struct{}, maxRequestsPerPeriod),
		delay:  make(chan struct{}, 1),
	}

	go f.refreshBucket(ctx)

	return f
}

type ftxRateLimiter struct {
	bucket chan struct{}
	delay  chan struct{}
}

func (f *ftxRateLimiter) RefreshWait(header http.Header, statusCode int) {
	if statusCode == 429 {
		select {
		case f.delay <- struct{}{}:
		default:
		}
	}
}

func (f *ftxRateLimiter) Wait() {
	select {
	case <-f.delay:
		// Hard delay.
		time.Sleep(1 * time.Second)
	case <-f.bucket:
		// Soft delay.
	case <-time.After(3 * time.Second):
		slog.Info(context.Background(), "FTX rate limit timeout hit: breaking the glass")
	}
}

func (f *ftxRateLimiter) refreshBucket(ctx context.Context) {
	// Fill bucket on initialization.
	for i := 0; i < maxRequestsPerPeriod; i++ {
		f.bucket <- struct{}{}
	}

	// Start refresh loop.
	for {
		select {
		case <-time.After(periodicity / maxRequestsPerPeriod):
		case <-ctx.Done():
			slog.Warn(ctx, "FTX rate limiter shutting down: context cancelled")
			return
		}

		// Accrue new request token.
		select {
		case f.bucket <- struct{}{}:
		default:
		}
	}
}
