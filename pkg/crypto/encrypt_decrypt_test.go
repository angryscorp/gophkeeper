package crypto_test

import (
	"bytes"
	"testing"

	"gophkeeper/pkg/crypto"
)

func TestEncryptDecrypt_Roundtrip(t *testing.T) {
	key := make([]byte, 32) // The null key is also valid AES-256
	plaintext := []byte("hello secret world")

	ciphertext, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt error: %v", err)
	}

	decrypted, err := crypto.Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("Decrypt error: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Fatalf("roundtrip failed: got %q, want %q", decrypted, plaintext)
	}
}

func TestDecrypt_InvalidCiphertext(t *testing.T) {
	key := make([]byte, 32)
	// Ciphertext too short
	_, err := crypto.Decrypt(key, []byte("short"))
	if err == nil {
		t.Fatalf("expected error on short ciphertext, got nil")
	}

	// Corrupted ciphertext
	plaintext := []byte("important")
	ciphertext, _ := crypto.Encrypt(key, plaintext)
	ciphertext[len(ciphertext)-1] ^= 0xFF // breaking the last byte

	_, err = crypto.Decrypt(key, ciphertext)
	if err == nil {
		t.Fatalf("expected error on tampered ciphertext, got nil")
	}
}
