package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundToPrecisionString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		input         float64
		precision     int
		expectedValue string
	}{
		{
			name:          "sol_example",
			input:         2.6,
			precision:     0,
			expectedValue: "3",
		},
		{
			name:          "zero_value",
			input:         0,
			precision:     100,
			expectedValue: "",
		},
		{
			name:          "4_decimal_precision",
			input:         0.12346,
			precision:     4,
			expectedValue: "0.1235",
		},
		{
			name:          "bitcoin_example",
			input:         1.267,
			precision:     2,
			expectedValue: "1.27",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := roundToPrecisionString(tt.input, tt.precision)
			assert.Equal(t, tt.expectedValue, result)
		})
	}
}
