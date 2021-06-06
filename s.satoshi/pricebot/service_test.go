package pricebot

import (
	"context"
	"fmt"
	"testing"

	coingecko "swallowtail/s.coingecko/clients"
	coingeckomock "swallowtail/s.coingecko/clients/mocks"
)

func TestService_GetPricesAsFormattedString(t *testing.T) {
	// Not running in parallel since we mutate our factory function for creating
	// a coingecko client.
	var (
		ctx = context.Background()
	)
	tests := []struct {
		symbols []string
	}{
		{
			symbols: []string{},
		},
		{
			symbols: []string{
				"BTC",
				"ETH",
				"LTC",
				"LINK",
				"SOL",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("%v_symbols", len(tt.symbols)), func(t *testing.T) {
			// Create mock.
			m := &coingeckomock.CoinGeckoClient{}

			// Replace the default factory function.
			orig := defaultCoingeckoClient
			defaultCoingeckoClient = func(ctx context.Context) coingecko.CoinGeckoClient {
				return m
			}

			// Cleanup after ourselves.
			t.Cleanup(func() {
				defaultCoingeckoClient = orig
			})

			for _, s := range tt.symbols {
				m.On("GetCurrentPriceFromSymbol", ctx, s, "usd").Return(0.0, nil)
			}

			// We don't particular care for the returned string format; that is tested elsewhere.
			// Here we just want to check our mock has been called.
			svc := NewService(ctx)
			svc.GetPricesAsFormattedString(ctx, tt.symbols, false)

			// Lets check we have the correct calls.
			for _, s := range tt.symbols {
				m.AssertCalled(t, "GetCurrentPriceFromSymbol", ctx, s, "usd")
			}
		})
	}
}
