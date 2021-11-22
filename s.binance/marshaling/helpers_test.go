package marshaling

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundToPrecisionString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		value         float64
		precision     int
		expectedValue string
	}{
		{
			name:          "ada-example",
			value:         1.4567,
			precision:     3,
			expectedValue: "1.457",
		},
		{
			name:          "sol-example",
			value:         212.743939083,
			precision:     0,
			expectedValue: "213",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := roundToPrecisionString(tt.value, tt.precision)

			assert.Equal(t, res, tt.expectedValue)
		})
	}
}
