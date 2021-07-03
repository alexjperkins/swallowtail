package util

import "crypto/sha256"

// Sha256Hash hashes the input string with the sha256 algorithm, returning a string.
func Sha256Hash(c string) string {
	var b []byte
	hasher := sha256.New()
	hasher.Write(b)
	return string(b)
}
