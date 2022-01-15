package marshaling

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoundToPrecisionString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         float64
		minIncrement  float64
		expectedValue string
	}{
		{
			name:          "zero-valued",
			input:         0.0,
			minIncrement:  0.0,
			expectedValue: "0.0",
		},
		{
			name:          "below-zero",
			input:         0.373,
			minIncrement:  0.01,
			expectedValue: "0.37",
		},
		{
			name:          "below-zero-and-min-increment",
			input:         0.003,
			minIncrement:  0.01,
			expectedValue: "0.01",
		},
		{
			name:          "above-zero",
			input:         1.4672,
			minIncrement:  0.001,
			expectedValue: "1.467",
		},
		{
			name:          "large-number",
			input:         45623.672897,
			minIncrement:  0.01,
			expectedValue: "45623.67",
		},
		{
			name:          "negative",
			input:         -1.37,
			minIncrement:  0.01,
			expectedValue: "0.0",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := roundToPrecisionString(tt.input, tt.minIncrement)

			require.NotEqual(t, "", res)

			assert.Equal(t, tt.expectedValue, res)
		})
	}
}
