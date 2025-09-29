package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
)

// Encrypt encrypts the plaintext using the provided 32-byte key and returns the ciphertext or an error.
func Encrypt(key, plaintext []byte) ([]byte, error) {
	// Create the cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create the GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generate a random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// Encrypt the plaintext
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

// Decrypt decrypts the ciphertext using the provided 32-byte key and returns the plaintext or an error.
func Decrypt(key, ciphertext []byte) ([]byte, error) {
	// Create the cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create the GCM
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Decrypt the ciphertext
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// Decrypt the ciphertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
