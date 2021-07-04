package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"swallowtail/libraries/util"

	"github.com/monzo/terrors"
)

// EncryptWithAES encrypts the data with a passphrase using the AES cipher.
func EncryptWithAES(d []byte, passphrase string) (string, error) {
	hash, err := util.Sha256Hash(passphrase)
	if err != nil {
		return "", terrors.Augment(err, "Failed to encrypt with AES cipher; error hashing passphrase", nil)
	}

	block, err := aes.NewCipher([]byte(hash))
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

	return string(gcm.Seal(nonce, nonce, d, nil)), nil
}

// DecryptWithAES decrypts the data with a passphrase using the AES cipher.
func DecryptWithAES(d []byte, passphrase string) (string, error) {
	hash, err := util.Sha256Hash(passphrase)
	if err != nil {
		return "", terrors.Augment(err, "Failed to decrypt with AES cipher; error hashing passphrase", nil)
	}

	key := []byte(hash)
	block, err := aes.NewCipher(key)
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
