package crypto_test

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"testing"

	"gophkeeper/pkg/crypto"
)

func TestSignChallenge_HMACSHA256(t *testing.T) {
	authKey := []byte("32-bytes-auth-key---------------") // 32 bytes
	ch := []byte("challenge")

	got := crypto.SignChallenge(authKey, ch, crypto.AuthKeyAlgorithmHMACSHA256)

	// Standard library reference
	mac := hmac.New(sha256.New, authKey)
	mac.Write(ch)
	want := mac.Sum(nil)

	if !hmac.Equal(got, want) {
		t.Fatalf("HMAC-SHA256 mismatch")
	}
	if len(got) != 32 {
		t.Fatalf("expected 32-byte tag for SHA256, got %d", len(got))
	}
}

func TestSignChallenge_HMACSHA512(t *testing.T) {
	authKey := []byte("32-bytes-auth-key---------------")
	ch := []byte("challenge")

	got := crypto.SignChallenge(authKey, ch, crypto.AuthKeyAlgorithmHMACSHA512)

	mac := hmac.New(sha512.New, authKey)
	mac.Write(ch)
	want := mac.Sum(nil)

	if !bytes.Equal(got, want) {
		t.Fatalf("HMAC-SHA512 mismatch")
	}
	if len(got) != 64 {
		t.Fatalf("expected 64-byte tag for SHA512, got %d", len(got))
	}
}

func TestSignChallenge_UnknownAlgoPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic on unknown algorithm, got none")
		}
	}()
	_ = crypto.SignChallenge([]byte("k"), []byte("c"), crypto.AuthKeyAlgorithm("UNKNOWN"))
}
