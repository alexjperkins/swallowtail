package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSha256Hash(t *testing.T) {
	t.Parallel()

	resA := Sha256Hash("hello")
	resB := Sha256Hash("hello")
	resC := Sha256Hash("bad answer")

	assert.Equal(t, 32, len(resA))

	assert.Equal(t, resA, resB)
	assert.NotEqual(t, resA, resC)
}
