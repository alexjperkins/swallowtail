package util

import "strings"

// MaskKey will mask any key, losing all information expect for the last k-digits.
func MaskKey(s string, k int) string {
	l := len(s)
	if l < k {
		return s
	}

	return strings.Repeat("*", l-k) + s[l-k:]
}
