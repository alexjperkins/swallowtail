package client

import (
	"context"
	"fmt"
	"swallowtail/libraries/gerrors"
	"swallowtail/s.binance/client/auth"
	"time"
)

func (c *binanceClient) do(ctx context.Context, method, endpoint, queryString string, reqBody, rspBody interface{}, credentials *Credentials) error {
	formattedEndpoint := fmt.Sprintf("%s?%s", endpoint, queryString)
	if credentials == nil {
		return c.http.Do(ctx, method, formattedEndpoint, reqBody, rspBody)
	}

	return c.http.DoWithEphemeralHeaders(ctx, method, formattedEndpoint, reqBody, rspBody, credentials.AsHeaders())
}

func (c *binanceClient) doWithSignature(ctx context.Context, method, endpoint, queryString string, reqBody, rspBody interface{}, credentials *Credentials) error {
	errParams := map[string]string{
		"method":   method,
		"endpoint": endpoint,
	}

	// First check that credentials have indeed been passed correctly.
	switch {
	case credentials == nil:
		return gerrors.FailedPrecondition("cannot_sign_binance_request.nil_credentials", errParams)
	case credentials.SecretKey == "":
		return gerrors.FailedPrecondition("cannot_sign_binance_request.empty_secret_key", errParams)
	}

	// Sign our request with the secret key passed.
	signedRequest, err := c.signRequest(credentials.SecretKey, queryString, reqBody)
	if err != nil {
		return gerrors.Augment(err, "failed_do_request.signature_failure", errParams)
	}

	formattedEndpoint := fmt.Sprintf("%s?%s", endpoint, signedRequest)

	return c.http.DoWithEphemeralHeaders(ctx, method, formattedEndpoint, reqBody, rspBody, credentials.AsHeaders())
}

func (c *binanceClient) signRequest(secret, queryString string, reqBody interface{}) (string, error) {
	// converts to unix nano time to that of millisecond precision; this is all that we need.
	now := time.Now().UnixNano() / 1_000_000

	// sign the request
	hmac, err := auth.Sha256HMAC(secret, queryString, now, reqBody)
	if err != nil {
		return "", gerrors.Augment(err, "failed_to_sign_request", nil)
	}

	// Return the new converted querystring with timestamp & signature appended.
	switch {
	case queryString == "":
		return fmt.Sprintf("timestamp=%d&signature=%s", now, hmac), nil
	default:
		return fmt.Sprintf("%s&timestamp=%d&signature=%s", queryString, now, hmac), nil
	}
}
