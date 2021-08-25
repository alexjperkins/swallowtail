package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"

	"github.com/monzo/terrors"

	"swallowtail/libraries/gerrors"
)

// EncryptWithAES encrypts the data with a passphrase using the AES cipher.
func EncryptWithAES(dest []byte, passphrase string) (string, error) {
	hashed := sha256.Sum256([]byte(passphrase))

	block, err := aes.NewCipher(hashed[:])
	if err != nil {
		return "", terrors.Augment(err, "Failed to encrypt with AES cipher; failed to create cipher block", nil)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", terrors.Augment(err, "Failed to encrypt with AES cipher; failed to create cipher gcm", nil)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", terrors.Augment(err, "Failed to encrypt with AES cipher; failed to create nonce", nil)
	}

	return hex.EncodeToString(gcm.Seal(nonce, nonce, dest, nil)), nil
}

// DecryptWithAES decrypts the data with a passphrase using the AES cipher.
func DecryptWithAES(encryptedText string, passphrase string) (string, error) {
	d, err := hex.DecodeString(encryptedText)
	if err != nil {
		return "", gerrors.Augment(err, "aes_decryption_failed", nil)
	}
	hashed := sha256.Sum256([]byte(passphrase))

	block, err := aes.NewCipher(hashed[:])
	if err != nil {
		return "", terrors.Augment(err, "Failed to decrypt with AES cipher; failed to create cipher block", nil)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", terrors.Augment(err, "Failed to decrypt with AES cipher; failed to create cipher gcm", nil)
	}

	nonce, ciphertext := d[:gcm.NonceSize()], d[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", terrors.Augment(err, "Failed to decrypt with AES cipher; failed to decrypt", nil)
	}

	return string(plaintext), nil
}
