package pricebot

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildDiscordMessageFromPrices(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		prices []*PriceBotPrice
	}{
		{
			name:   "empty_prices",
			prices: []*PriceBotPrice{},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := buildMessage(tt.prices, false)
			switch {
			case len(tt.prices) == 0:
				assert.Equal(t, "", res)
			}
		})
	}
}

func TestBuildLine(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		price        *PriceBotPrice
		expectedLine string
	}{
		{
			name: "zero-price",
			price: &PriceBotPrice{
				Price:  0.0,
				Symbol: "BTC",
			},
			expectedLine: "[BTC]: [N/A]",
		},
		{
			name: "fetched-price",
			price: &PriceBotPrice{
				Price:  50000.000,
				Symbol: "ETH",
			},
			expectedLine: "[ETH]: [50000.000]",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			line := buildLine(tt.price)
			assert.Equal(t, tt.expectedLine, line)
		})
	}
}
