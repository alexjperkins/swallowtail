package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"swallowtail/libraries/gerrors"
)

// Sha256HMAC signs some request using the HMAC SHA256 algorithm as required by binance.
// https://binance-docs.github.io/apidocs/#endpoint-security-type
func Sha256HMAC(secret, queryString string, timestamp int64, requestBody interface{}) (string, error) {
	switch {
	case queryString == "" && requestBody == nil:
		return sha256HMACQueryString(secret, queryString, timestamp)
	case queryString == "":
		return sha256HMACBody(secret, timestamp, requestBody)
	case requestBody == nil:
		return sha256HMACQueryString(secret, queryString, timestamp)
	default:
		return sha256HMACBodyAndQueryString(secret, queryString, timestamp, requestBody)
	}
}

func sha256HMACBody(secret string, timestamp int64, requestBody interface{}) (string, error) {
	return "", gerrors.Unimplemented("binance_signer_body_unimplemented", nil)
}

func sha256HMACQueryString(secret, queryString string, timestamp int64) (string, error) {
	queryStringAndTimestamp := constructQueryStringWithTimestamp(queryString, timestamp)

	h := hmac.New(sha256.New, []byte(secret))

	_, err := h.Write([]byte(queryStringAndTimestamp))
	if err != nil {
		return "", gerrors.Augment(err, "hmac_sha256_query_string_failed", map[string]string{
			"query_string": queryString,
		})
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func sha256HMACBodyAndQueryString(secret, queryString string, timestamp int64, requestBody interface{}) (string, error) {
	return "", gerrors.Unimplemented("binance_signer_reqbody_and_query_string_unimplemented", nil)
}

func constructQueryStringWithTimestamp(queryString string, timestamp int64) string {
	switch {
	case queryString == "":
		return fmt.Sprintf("timestamp=%d", timestamp)
	default:
		return fmt.Sprintf("%s&timestamp=%d", queryString, timestamp)
	}
}
