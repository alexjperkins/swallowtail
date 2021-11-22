package exchangeinfo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculatePrecision(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name              string
		value             string
		expectedPrecision int
	}{
		{
			name:              "sol-example",
			value:             "1",
			expectedPrecision: 0,
		},
		{
			name:              "ada-example",
			value:             "0.01",
			expectedPrecision: 2,
		},
		{
			name:              "rsr-example",
			value:             "0.0001",
			expectedPrecision: 4,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := calculatePrecision(tt.value)

			assert.Equal(t, tt.expectedPrecision, res)
		})
	}
}
