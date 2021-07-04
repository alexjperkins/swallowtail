package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatPrice(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		input          string
		expectedOutput string
	}{
		{
			name:           "empty-string",
			input:          "",
			expectedOutput: "",
		},
		{
			name:           "integer",
			input:          "999",
			expectedOutput: "999.000",
		},
		{
			name:           "float",
			input:          "999.99",
			expectedOutput: "999.990",
		},
		{
			name:           "bitcoin",
			input:          "58765.437",
			expectedOutput: "58765.437",
		},
		{
			name:           "bitcoin-round",
			input:          "58700",
			expectedOutput: "58700.000",
		},
		{
			name:           "float-under-one",
			input:          "0.0012305",
			expectedOutput: "0.00123",
		},
		{
			name:           "another-float-under-one",
			input:          "0.0654389",
			expectedOutput: "0.0654",
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r, err := FormatPriceFromString(tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, r)
		})
	}
}

func TestFormatPriceAsString(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		input          float64
		expectedOutput string
	}{
		{
			name:           "zero",
			input:          0.0,
			expectedOutput: "0.00",
		},
		{
			name:           "above-one",
			input:          10.7649,
			expectedOutput: "10.765",
		},
		{
			name:           "under-one",
			input:          0.058339,
			expectedOutput: "0.0583",
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			res, err := FormatPriceAsString(tc.input)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedOutput, res)
		})
	}
}
