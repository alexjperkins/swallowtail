package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildQueryStringFromFuturesPerpetualTrade(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		input          *ExecutePerpetualFuturesTradeRequest
		expectedOutput string
	}{
		{
			name: "sol_example_market",
			input: &ExecutePerpetualFuturesTradeRequest{
				Symbol:           "SOLUSDT",
				Side:             "BUY",
				Type:             "MARKET",
				Quantity:         "3",
				NewOrderRespType: "ACK",
			},
			expectedOutput: "symbol=SOLUSDT&side=BUY&type=MARKET&quantity=3&newOrderRespType=ACK",
		},
		{
			name: "sol_example_limit",
			input: &ExecutePerpetualFuturesTradeRequest{
				Symbol:           "SOLUSDT",
				Side:             "BUY",
				Type:             "LIMIT",
				Quantity:         "3",
				NewOrderRespType: "ACK",
				Price:            "150",
				TimeInForce:      "GTC",
			},
			expectedOutput: "symbol=SOLUSDT&side=BUY&type=LIMIT&timeInForce=GTC&quantity=3&price=150&newOrderRespType=ACK",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			output := buildQueryStringFromFuturesPerpetualTrade(tt.input)
			assert.Equal(t, tt.expectedOutput, output)

		})
	}

}
