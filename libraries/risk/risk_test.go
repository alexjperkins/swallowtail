package risk

import (
	"math"
	"testing"

	tradeengineproto "swallowtail/s.trade-engine/proto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalculateRiskPositionsByRisk(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		entries   []float64
		totalRisk float64
		stopLoss  float64
		howMany   int
		side      tradeengineproto.TRADE_SIDE
		strategy  tradeengineproto.DCA_EXECUTION_STRATEGY
	}{
		{
			name:      "5_entries",
			entries:   []float64{100, 200},
			totalRisk: 10,
			stopLoss:  80,
			howMany:   5,
			side:      tradeengineproto.TRADE_SIDE_BUY,
			strategy:  tradeengineproto.DCA_EXECUTION_STRATEGY_LINEAR,
		},
		{
			name:      "7_entries",
			entries:   []float64{10, 12},
			totalRisk: 5,
			stopLoss:  8,
			howMany:   7,
			side:      tradeengineproto.TRADE_SIDE_BUY,
			strategy:  tradeengineproto.DCA_EXECUTION_STRATEGY_LINEAR,
		},
		{
			name:      "5_entries_real_example",
			entries:   []float64{3200, 3550},
			totalRisk: 5,
			stopLoss:  2500,
			howMany:   5,
			side:      tradeengineproto.TRADE_SIDE_BUY,
			strategy:  tradeengineproto.DCA_EXECUTION_STRATEGY_LINEAR,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			positions, err := CalculatePositionsByRisk(tt.entries, tt.stopLoss, tt.totalRisk, tt.howMany, tt.side, tt.strategy)
			require.NoError(t, err)

			d := diff(tt.totalRisk/100, sumRisk(positions))
			assert.True(t, d < 0.1, "Got: %f, expecting: %f", d, tt.totalRisk/100)
		})
	}
}

func TestSummedLinspace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		howMany int
		total   float64
	}{
		{
			name:    "non_zero_odd",
			howMany: 5,
			total:   100,
		},
		{
			name:    "non_zero_odd_3",
			howMany: 5,
			total:   1,
		},
		{
			name:    "non_zero_odd_2",
			howMany: 3,
			total:   50,
		},
		{
			name:    "non_zero_odd_3",
			howMany: 5,
			total:   1,
		},
		{
			name:    "non_zero_even",
			howMany: 4,
			total:   1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := summedLinspace(tt.howMany, tt.total)

			assert.Len(t, result, tt.howMany)
			// Since we're using floats we can assert this is correct if the difference
			// between the result and expected value is 1.
			assert.True(t, diff(tt.total, sum(result)) < 0.00001)
		})
	}
}

func sum(vs []float64) float64 {
	if len(vs) == 0 {
		return 0
	}

	return vs[0] + sum(vs[1:])
}

func sumRisk(ps []*Position) float64 {
	if len(ps) == 0 {
		return 0
	}

	return ps[0].Risk + sumRisk(ps[1:])
}

func sumContracts(ps []*Position) float64 {
	if len(ps) == 0 {
		return 0
	}

	return ps[0].Risk*ps[0].Price + sumRisk(ps[1:])
}

func diff(a, b float64) float64 {
	return math.Abs(a - b)
}
