package encrypting

import (
	"crypto/aes"
	"crypto/cipher"
)

// Encrypt Encryption function.
func Encrypt(stringToEncrypt string, key []byte) ([]byte, error) {
	plaintext := []byte(stringToEncrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce, err := GenerateRandom(aesGCM.NonceSize())
	if err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}
