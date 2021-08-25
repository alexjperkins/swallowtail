package encryption

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"swallowtail/libraries/util"
)

func TestAESCipherEncryption_16Bytes(t *testing.T) {
	t.Parallel()

	var (
		input    = util.RandString(16, util.AlphaNumeric)
		password = util.RandString(16, util.AlphaNumeric)
	)

	ciphertext, err := EncryptWithAES([]byte(input), password)
	require.NoError(t, err)
	require.NotEqual(t, input, ciphertext)
	require.NotEqual(t, password, ciphertext)

	plaintext, err := DecryptWithAES(ciphertext, password)
	require.NoError(t, err)
	require.NotEqual(t, password, plaintext)

	assert.Equal(t, input, plaintext)
}

func TestAESCipherEncryption_64Bytes(t *testing.T) {
	t.Parallel()

	var (
		input    = util.RandString(64, util.AlphaNumeric)
		password = util.RandString(64, util.AlphaNumeric)
	)

	ciphertext, err := EncryptWithAES([]byte(input), password)
	require.NoError(t, err)
	require.NotEqual(t, input, ciphertext)
	require.NotEqual(t, password, ciphertext)

	plaintext, err := DecryptWithAES(ciphertext, password)
	require.NoError(t, err)
	require.NotEqual(t, password, plaintext)

	assert.Equal(t, input, plaintext)
}
