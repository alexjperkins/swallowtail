package util

import (
	"crypto/sha256"
	"encoding/hex"
)

// Sha256Hash hashes the input string with the sha256 algorithm, returning a string.
// It's worth noting this function returns a string of 64 characters.
func Sha256Hash(c string) string {
	b := []byte(c)
	hashed := sha256.Sum256(b)

	return hex.EncodeToString(hashed[:])
}
