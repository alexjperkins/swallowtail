package auth

import (
	"fmt"
	"swallowtail/libraries/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSha256HMAC(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		signature   string
		queryString string
		timestamp   int64
		requestBody interface{}
	}{
		{
			name:        "with_empty_query_string_empty_request_body",
			signature:   util.RandString(32, util.AlphaNumeric),
			queryString: "",
			timestamp:   time.Now().Truncate(time.Millisecond).UnixNano() / 1000000,
			requestBody: nil,
		},
		{
			name:        "with_query_string_empty_request_body",
			signature:   util.RandString(32, util.AlphaNumeric),
			queryString: "symbol=BTCUSDT",
			timestamp:   time.Now().Truncate(time.Millisecond).UnixNano() / 1000000,
			requestBody: nil,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res, err := Sha256HMAC(tt.signature, tt.queryString, tt.timestamp, tt.requestBody)
			require.NoError(t, err)

			assert.NotEqual(t, "", res)
			fmt.Println(res)
		})
	}
}
