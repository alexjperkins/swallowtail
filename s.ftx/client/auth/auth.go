package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"swallowtail/libraries/gerrors"
)

// Sign signs some signature payload with the given secret as required for FTX.
// https://docs.ftx.com/#rest-api
func SignRequest(endpoint, method, timestamp string, req interface{}, credentials *Credentials) (string, error) {
	var rawReq []byte
	if req != nil {
		var err error

		rawReq, err = json.Marshal(req)
		if err != nil {
			return "", gerrors.Augment(err, "failed_to_sign_request.marshal_request_to_bytes", nil)
		}
	}

	signaturePayload := fmt.Sprintf("%s%s%s%s", timestamp, method, endpoint, rawReq)

	h := hmac.New(sha256.New, []byte(credentials.SecretKey))
	if _, err := h.Write([]byte(signaturePayload)); err != nil {
		return "", gerrors.Augment(err, "failed_to_sign_with_hmac_sha256", nil)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// Credentials holds the credentials required for FTX.
type Credentials struct {
	APIKey     string
	SecretKey  string
	Subaccount string
	URL        string
	WSURL      string
}

// AsHeaders converts the credentials struct into the headers required to verify the user.
// It uses the request body & the timestamp to sign the request.
func (c *Credentials) AsHeaders(signature, timestamp string) map[string]string {
	if c == nil {
		return map[string]string{}
	}

	m := map[string]string{
		"Content-Type": "application/json",
		"FTX-KEY":      c.APIKey,
		"FTX-SIGN":     signature,
		"FTX-TS":       timestamp,
	}

	if c.Subaccount != "" {
		m["FTX-SUBACCOUNT"] = c.Subaccount
	}

	return m
}

// SubaccountAsHeaders returns only the subaccount as headers; if it is not null.
func (c *Credentials) SubaccountAsHeaders() map[string]string {
	if c == nil {
		return map[string]string{}
	}

	m := map[string]string{
		"Content-Type": "application/json",
	}
	if c.Subaccount != "" {
		m["FTX-SUBACCOUNT"] = c.Subaccount
	}

	return m
}
