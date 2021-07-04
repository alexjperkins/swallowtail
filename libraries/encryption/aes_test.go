package encryption

import (
	"swallowtail/libraries/util"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAESCipherEncryption(t *testing.T) {
	t.Parallel()

	var (
		input    = util.RandString(16, util.AlphaNumbmeric)
		password = util.RandString(16, util.AlphaNumbmeric)
	)

	ciphertext, err := EncryptWithAES([]byte(input), password)
	require.NoError(t, err)
	require.NotEqual(t, input, ciphertext)
	require.NotEqual(t, password, ciphertext)

	plaintext, err := DecryptWithAES([]byte(ciphertext), password)
	require.NoError(t, err)
	require.NotEqual(t, password, plaintext)

	assert.Equal(t, input, plaintext)
}
