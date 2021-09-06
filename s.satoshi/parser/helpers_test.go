package parser

import (
	"testing"

	tradeengineproto "swallowtail/s.trade-engine/proto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCaptureNumbers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		content        string
		expectedFloats []float64
	}{
		{
			name:           "decimals-and-no-decimals",
			content:        "blah 100.65 100 blah 88.3579 blah",
			expectedFloats: []float64{100.65, 100, 88.3579},
		},
		{
			name:           "wwg-case-1",
			content:        "[Rego]: BTC LONG 50000 49000 TP 52000 54000 58000",
			expectedFloats: []float64{50000, 49000, 52000, 54000, 58000},
		},
		{
			name:           "wwg-case-2-with-percentages",
			content:        "[Rego]: BTC LONG for a nice 10% 50000 49000 52000 54000 58000",
			expectedFloats: []float64{50000, 49000, 52000, 54000, 58000},
		},
		{
			name:           "swings-case-1-with-rr",
			content:        "[Rego]: BTC LONG for a nice 10% 50000 49000 52000 54000 58000 4.96RR",
			expectedFloats: []float64{50000, 49000, 52000, 54000, 58000},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			floats, err := captureNumbers(tt.content)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedFloats, floats)
		})
	}
}

func TestFindSide(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		content      string
		expectedSide tradeengineproto.TRADE_SIDE
		isOk         bool
	}{
		{
			name:         "buy-side",
			content:      "[A]: LONG BTC 50900.12 45000.43",
			expectedSide: tradeengineproto.TRADE_SIDE_BUY,
			isOk:         true,
		},
		{
			name:         "sell-side",
			content:      "[T]: Short SOLUSDT 500 550 400 350 200",
			expectedSide: tradeengineproto.TRADE_SIDE_SELL,
			isOk:         true,
		},
		{
			name:         "no-side",
			content:      "[T]: SOLUSDT 500 550 400 350 200",
			expectedSide: tradeengineproto.TRADE_SIDE_BUY,
			isOk:         false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			side, ok := findSide(tt.content)
			assert.Equal(t, tt.isOk, ok)

			assert.Equal(t, tt.expectedSide, side)
		})
	}
}
