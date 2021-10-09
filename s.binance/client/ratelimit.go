package client

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/monzo/slog"
)

const (
	maxNumberOfOrdersPer10Seconds = 10
)

type binanceRateLimiter struct {
	toWait time.Duration
}

func (b *binanceRateLimiter) Wait() {
	// We ideally want metrics here.
	time.Sleep(b.toWait)
}

func (b *binanceRateLimiter) RefreshWait(header http.Header, statusCode int) {
	v := header.Get("X-MBX-ORDER-COUNT-10S")
	if statusCode == 429 {
		slog.Warn(context.Background(), "Binance http client has been rate limited, sleeping for 3 seconds: %v", v)
		b.toWait = 3 * time.Second
		return
	}

	howManyLeft, err := strconv.ParseInt(v, 64, 2)
	if err != nil {
		slog.Error(context.Background(), "Faield to parse binance rate limiter order count: %v", v)
		return
	}

	overhead := float64(maxNumberOfOrdersPer10Seconds-howManyLeft) / maxNumberOfOrdersPer10Seconds
	b.toWait = time.Duration(overhead) * time.Second
}
