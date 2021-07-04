package util

import (
	"crypto/sha256"
	"fmt"
)

// Sha256Hash hashes the input string with the sha256 algorithm, returning a string.
func Sha256Hash(c string) (string, error) {
	var b []byte
	hasher := sha256.New()
	_, err := hasher.Write([]byte(c))
	if err != nil {
		fmt.Println("ERROR", err)

	}

	return string(hasher.Sum(b)), nil
}
