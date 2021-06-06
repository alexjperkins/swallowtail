package satoshi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		a              float64
		expectedResult float64
	}{
		{
			name:           "above_zero",
			a:              10.0,
			expectedResult: 10.0,
		},
		{
			name:           "below_zero",
			a:              -10.0,
			expectedResult: 10.0,
		},
		{
			name:           "zero",
			a:              0.0,
			expectedResult: 0.0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := abs(tt.a)
			assert.Equal(t, tt.expectedResult, res)
		})
	}
}
