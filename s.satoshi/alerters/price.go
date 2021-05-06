package alerters

import (
	"context"
	"fmt"
	"swallowtail/libraries/structures/window"
	"time"

	"github.com/monzo/slog"
)

type PriceAlerter struct {
	symbol          string
	targetPrice     float64
	c               ExchangeClient
	m               MessageClient
	interval        time.Duration
	withJitter      bool
	done            chan struct{}
	flagApproaching bool
	deltaPercentage float64
}

type ExchangeClient interface {
	Ping() bool
	GetPrice(ctx context.Context, symbol string) (float64, error)
	ID() string
}

type MessageClient interface {
	Send(message string) error
}

func NewPriceAlerter(symbol string, exchangeClient ExchangeClient) *PriceAlerter {
	return &PriceAlerter{
		c: exchangeClient,
	}
}

func (pa *PriceAlerter) Run(ctx context.Context) {
	t := time.NewTicker(pa.interval)
	cache := window.NewMovingWindow(16)
	defer slog.Info(ctx, "Closing down price alerter", map[string]string{
		"symbol": pa.symbol,
	})
	if pa.withJitter {
		time.Sleep(time.Minute)
	}
	for {
		select {
		case <-t.C:
			price, err := pa.c.GetPrice(ctx, pa.symbol)
			if err != nil {
				slog.Error(ctx, "Failed to fetch price", map[string]string{
					"exchange_client_id": pa.c.ID(),
					"error":              err.Error(),
				})
			}
			defer cache.Push(float32(price))
			if priceCrossedTarget(cache, price) {
				pa.m.Send(fmt.Sprintf("%s has crossed target price %.4f: %.4f", pa.symbol, pa.targetPrice, price))
				continue
			}
			if !pa.flagApproaching {
				continue
			}
			if isApproaching(cache, price, pa.deltaPercentage) {
				pa.m.Send(fmt.Sprintf("%s is approaching price %.4f: %.4f", pa.symbol, pa.targetPrice, price))
				continue
			}

		case <-pa.done:
			return
		case <-ctx.Done():
			return
		}
	}
}

func isApproaching(w *window.MovingWindow, p, d float64) bool {
	return false
}

func priceCrossedTarget(w *window.MovingWindow, p float64) bool {
	return false
}
