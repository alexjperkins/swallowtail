package signer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"swallowtail/libraries/gerrors"
)

// Sha256HMAC signs some signature payload with the given secret as required for FTX.
// https://docs.ftx.com/#rest-api
// TODO: we also need to take care of the request body.
func Sha256HMAC(secret string, signaturePayload string) (string, error) {
	h := hmac.New(sha256.New, []byte(secret))
	if _, err := h.Write([]byte(signaturePayload)); err != nil {
		return "", gerrors.Augment(err, "failed_to_sign_with_hmac_sha256", nil)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
