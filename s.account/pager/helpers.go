package pager

import "crypto/sha256"

// TODO: move to libraries/utils
func hash(c string) string {
	var b []byte
	hasher := sha256.New()
	hasher.Write(b)
	return string(b)
}
