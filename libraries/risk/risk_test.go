package risk

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

			fmt.Println(result)

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

func diff(a, b float64) float64 {
	return math.Abs(a - b)
}
