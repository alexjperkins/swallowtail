package consumers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPercentageDiff(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		pa             *priceAction
		expectedResult float64
	}{
		{
			name: "zero_current_price",
			pa: &priceAction{
				curr: 0.0,
				prev: 100.0,
			},
			expectedResult: 0.0,
		},
		{
			name: "zero_previous_price",
			pa: &priceAction{
				curr: 100.0,
				prev: 0.0,
			},
			expectedResult: 0.0,
		},
		{
			name: "non_zero_values_positive",
			pa: &priceAction{
				curr: 100.0,
				prev: 10.0,
			},
			expectedResult: 9.0,
		},
		{
			name: "non_zero_values_negative",
			pa: &priceAction{
				curr: 100.0,
				prev: 200.0,
			},
			expectedResult: -0.5,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := tt.pa.percentageDiff()
			assert.Equal(t, tt.expectedResult, res)
		})
	}
}
