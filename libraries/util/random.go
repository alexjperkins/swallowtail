package util

import (
	"math/rand"
	"strings"
)

const (
	AlphaNumeric Charset = iota
)

type Charset int

func (c Charset) String() string {
	switch c {
	case AlphaNumeric:
		return "abcdefghijklmnopqrstuvwxyz0123456789"
	}
	return ""
}

// RandString returns a pseudo-random string how length `k`, and from Charset `charset`
func RandString(k int, charset Charset) string {
	s := charset.String()
	var output strings.Builder

	for i := 0; i < k; i++ {
		r := rand.Intn(len(s))
		output.WriteString(string(s[r]))
	}

	return output.String()
}
