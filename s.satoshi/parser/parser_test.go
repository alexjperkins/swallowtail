package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanContent(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                   string
		content                string
		expectedCleanedContent string
	}{
		{
			name: "update_new_and_old_removed",
			content: `---OLD---
long btc 44.4  dca til 43.9 sl 43599 high risk
---NEW---
long btc 44.4 @​everyone dca til 43.9 sl 43399 high risk edit1 updated sl
			`,
			expectedCleanedContent: ` long btc 44.4 @​everyone dca til 43.9 sl 43399 high risk edit1 updated sl`,
		},
		{
			name: "dollar_signs",
			content: `
			XTZ $6.28 SL $6.02
			`,
			expectedCleanedContent: `xtz 6.28 sl 6.02`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cleanedContent := cleanContent(tt.content)

			assert.Equal(t, tt.expectedCleanedContent, cleanedContent)
		})
	}
}
