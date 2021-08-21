package handler

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseSpreadsheetIDFromURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		url                   string
		expectedSpreadsheetID string
		withErr               string
	}{
		{
			url:                   "https://docs.google.com/spreadsheets/d/1AYtRsdEcoEjmh-OtribxJ9et7qvCf6Z_UkkYNnKqqZY/edit#gid=1266119125",
			expectedSpreadsheetID: "1AYtRsdEcoEjmh-OtribxJ9et7qvCf6Z_UkkYNnKqqZY",
		},
		{
			url:                   "https://docs.google.com/spreadsheets/d/1bIsm2i28hdbolVmn14gLMYNvBV6G2dJ99rqsrf0FSqg/edit",
			expectedSpreadsheetID: "1bIsm2i28hdbolVmn14gLMYNvBV6G2dJ99rqsrf0FSqg",
		},
		{
			url:                   "https://docs.google.com/spreadsheets/d/1QDankgrs6mhfecdAzn0Xp0kE3PhhqISCyecI9rGFKHo/edit",
			expectedSpreadsheetID: "1QDankgrs6mhfecdAzn0Xp0kE3PhhqISCyecI9rGFKHo",
		},
	}

	for i, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("case_%v", i), func(t *testing.T) {
			t.Parallel()

			spreadsheetID, err := parseSpreadsheetIDFromURL(tt.url)
			switch {
			case tt.withErr == "":
				require.NoError(t, err)
				assert.Equal(t, tt.expectedSpreadsheetID, spreadsheetID)
			default:
				require.Error(t, err)
				assert.Equal(t, "", spreadsheetID)
			}
		})
	}
}
