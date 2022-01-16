package risk

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	tradeengineproto "swallowtail/s.trade-engine/proto"
)

func TestCalculateRiskCoefficient(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		accountBalance    float64
		totalRisk         float64
		entry             float64
		stopLoss          float64
		expectedTotalSize float64
	}{
		{
			name:              "basic_long",
			accountBalance:    1000,
			totalRisk:         10,
			entry:             100,
			stopLoss:          90,
			expectedTotalSize: 1000,
		},
		{
			name:              "short_basic",
			accountBalance:    20000,
			totalRisk:         5,
			entry:             9000,
			stopLoss:          10000,
			expectedTotalSize: 9000,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := CalculateRiskCoefficient(tt.entry, tt.stopLoss)

			totalQuantity := res * tt.accountBalance * tt.totalRisk
			assert.Equal(t, tt.expectedTotalSize, totalQuantity*tt.entry)
		})
	}
}

func TestCalculateRiskPositionsByRisk(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		entries  []float64
		stopLoss float64
		howMany  int
		side     tradeengineproto.TRADE_SIDE
		strategy tradeengineproto.DCA_EXECUTION_STRATEGY
	}{
		{
			name:     "5_entries",
			entries:  []float64{100, 200},
			stopLoss: 80,
			howMany:  5,
			side:     tradeengineproto.TRADE_SIDE_BUY,
			strategy: tradeengineproto.DCA_EXECUTION_STRATEGY_LINEAR,
		},
		{
			name:     "7_entries",
			entries:  []float64{10, 12},
			stopLoss: 8,
			howMany:  7,
			side:     tradeengineproto.TRADE_SIDE_BUY,
			strategy: tradeengineproto.DCA_EXECUTION_STRATEGY_LINEAR,
		},
		{
			name:     "5_entries_real_example",
			entries:  []float64{3200, 3550},
			stopLoss: 2500,
			howMany:  5,
			side:     tradeengineproto.TRADE_SIDE_BUY,
			strategy: tradeengineproto.DCA_EXECUTION_STRATEGY_LINEAR,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			positions, err := CalculatePositionsByRisk(tt.entries, tt.stopLoss, tt.howMany, tt.side, tt.strategy)
			require.NoError(t, err)

			d := diff(1.0/100, sumRisk(positions))
			assert.True(t, d < 0.1, "Got: %f, expecting: %f", d, 1.0/100)
		})
	}
}

func sum(vs []float64) float64 {
	if len(vs) == 0 {
		return 0
	}

	return vs[0] + sum(vs[1:])
}

func sumRisk(ps []*PositionDetail) float64 {
	if len(ps) == 0 {
		return 0
	}

	return ps[0].RiskCoefficient + sumRisk(ps[1:])
}

func sumContracts(ps []*PositionDetail) float64 {
	if len(ps) == 0 {
		return 0
	}

	return ps[0].RiskCoefficient*ps[0].Price + sumRisk(ps[1:])
}

func diff(a, b float64) float64 {
	return math.Abs(a - b)
}
