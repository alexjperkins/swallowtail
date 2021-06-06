package satoshi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateRisk(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                                     string
		entry, stopLoss, accountSize, percentage float64
		expectedNumContract                      float64
	}{
		{
			name:                "equal-entry-and-stop-loss",
			expectedNumContract: 0.0,
		},
		{
			name:                "10-percent-of-account",
			entry:               1.0,
			stopLoss:            0.9,
			accountSize:         100.0,
			percentage:          0.1,
			expectedNumContract: 100.0,
		},
		{
			name:                "50-percent-of-account",
			entry:               100.0,
			stopLoss:            80.0,
			accountSize:         1000.0,
			percentage:          0.5,
			expectedNumContract: 25.0,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res := calculateRisk(tt.entry, tt.stopLoss, tt.accountSize, tt.percentage)
			assert.InDelta(t, tt.expectedNumContract, res, 0.001)
		})
	}
}

func TestContains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		needle        string
		haystack      []string
		shouldContain bool
	}{
		{
			name:          "does-not-contain",
			needle:        "lol",
			haystack:      []string{"nope", "not, it", "fictional"},
			shouldContain: false,
		},
		{
			name:          "does-contain",
			needle:        "here",
			haystack:      []string{"yes", "I am", "here"},
			shouldContain: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			res := contains(tt.needle, tt.haystack)
			assert.Equal(t, tt.shouldContain, res)
		})
	}

}
