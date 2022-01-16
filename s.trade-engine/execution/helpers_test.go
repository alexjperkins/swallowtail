package execution

import (
	"swallowtail/libraries/risk"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateNumberOfDCABuys(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                 string
		accountBalance       float64
		expectedNumberOfBuys int
	}{
		{
			name:                 "small_account",
			accountBalance:       100,
			expectedNumberOfBuys: DCANumberOfBuysLowerBound,
		},
		{
			name:                 "small_account_on_boundary",
			accountBalance:       1000,
			expectedNumberOfBuys: DCANumberOfBuysLowerBound,
		},
		{
			name:                 "large_account_on_boundary",
			accountBalance:       1000.1,
			expectedNumberOfBuys: DCANumberOfBuysUpperBound,
		},
		{
			name:                 "large_account",
			accountBalance:       100000,
			expectedNumberOfBuys: DCANumberOfBuysUpperBound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := calculateNumberOfDCABuys(tt.accountBalance)

			assert.Equal(t, tt.expectedNumberOfBuys, res)
		})
	}
}

func TestCalculateTotalQuantityFromPositions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                  string
		accountBalance        float64
		totalRisk             float64
		positions             []*risk.PositionDetail
		expectedTotalQuantity float64
	}{
		{
			name:           "simple_dca_long_small_account",
			accountBalance: 100,
			totalRisk:      10,
			positions: []*risk.PositionDetail{
				{
					RiskCoefficient: 0.1,
				},
				{
					RiskCoefficient: 0.3,
				},
				{
					RiskCoefficient: 0.4,
				},
				{
					RiskCoefficient: 0.2,
				},
			},
			expectedTotalQuantity: 10, // 10% of 100
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := calculateTotalQuantityFromPositions(tt.accountBalance, tt.totalRisk, tt.positions)

			assert.Equal(t, tt.expectedTotalQuantity, res)
		})
	}
}
