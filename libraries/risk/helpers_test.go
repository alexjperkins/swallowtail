package risk

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSummedLinspace(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		howMany int
	}{
		{
			name:    "non_zero_odd",
			howMany: 5,
		},
		{
			name:    "non_zero_odd_3",
			howMany: 5,
		},
		{
			name:    "non_zero_odd_2",
			howMany: 3,
		},
		{
			name:    "non_zero_odd_3",
			howMany: 5,
		},
		{
			name:    "non_zero_even",
			howMany: 4,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := summedLinspace(tt.howMany)

			assert.Len(t, result, tt.howMany)

			// Since we're using floats we can assert this is correct if the difference
			// between the result and expected value is 1.
			assert.True(t, diff(1.0, sum(result)) < 0.00001)
		})
	}
}
