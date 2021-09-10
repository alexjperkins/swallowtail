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

			floats, err := parseNumbersFromContent(tt.content)
			require.NoError(t, err)

			assert.Equal(t, tt.expectedFloats, floats)
		})
	}
}

func TestParseSide(t *testing.T) {
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

			side, ok := parseSide(tt.content)
			assert.Equal(t, tt.isOk, ok)

			assert.Equal(t, tt.expectedSide, side)
		})
	}
}

func TestParseOrderType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		content           string
		currentValue      float64
		entry             float64
		side              tradeengineproto.TRADE_SIDE
		expectedOrderType tradeengineproto.ORDER_TYPE
	}{
		{
			name: "limit_order",
			content: `LIMIT LONG BTC 45000 SL 49000
			`,
			currentValue:      50000,
			entry:             45000,
			side:              tradeengineproto.TRADE_SIDE_BUY,
			expectedOrderType: tradeengineproto.ORDER_TYPE_LIMIT,
		},
		{
			name: "market_order",
			content: `LONG BTC 50000 SL 49000
			`,
			currentValue:      50000,
			entry:             49000,
			side:              tradeengineproto.TRADE_SIDE_BUY,
			expectedOrderType: tradeengineproto.ORDER_TYPE_MARKET,
		},
		{
			name: "limit_order_higher_entry_buy_side",
			content: `LIMIT LONG BTC 50000 SL 49000
			`,
			currentValue:      50000,
			entry:             55000,
			side:              tradeengineproto.TRADE_SIDE_BUY,
			expectedOrderType: tradeengineproto.ORDER_TYPE_MARKET,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			orderType, _ := parseOrderType(tt.content, tt.currentValue, tt.entry, tt.side)

			assert.Equal(t, tt.expectedOrderType, orderType)
		})
	}
}

func TestWithinRange(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		value, truth      float64
		rangeAsPercentage float64
		shouldBeInRange   bool
	}{
		{
			name:              "within-range",
			value:             50000,
			truth:             50500,
			rangeAsPercentage: 15,
			shouldBeInRange:   true,
		},
		{
			name:              "outside-range-upside",
			value:             50000,
			truth:             10000,
			rangeAsPercentage: 15,
			shouldBeInRange:   false,
		},
		{
			name:              "outside-range-downside",
			value:             10000,
			truth:             50000,
			rangeAsPercentage: 15,
			shouldBeInRange:   false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			isWithinRange := withinRange(tt.value, tt.truth, tt.rangeAsPercentage)

			assert.Equal(t, tt.shouldBeInRange, isWithinRange)
		})
	}
}
