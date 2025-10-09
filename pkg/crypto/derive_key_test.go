package crypto_test

import (
	"bytes"
	"testing"

	"gophkeeper/pkg/crypto"
)

func TestDeriveKey_Argon2ID_DeterministicAndLength(t *testing.T) {
	params := crypto.KDFParameters{
		Algorithm:   crypto.KDFAlgorithmARGON2ID,
		TimeCost:    2,
		MemoryCost:  64 * 1024,
		Parallelism: 2,
		Salt:        []byte("0123456789abcdef0123456789abcdef"),
	}

	k1, err := crypto.DeriveKey("password", params)
	if err != nil {
		t.Fatalf("DeriveKey error: %v", err)
	}
	k2, err := crypto.DeriveKey("password", params)
	if err != nil {
		t.Fatalf("DeriveKey error: %v", err)
	}

	if len(k1) != 32 {
		t.Fatalf("expected key length 32, got %d", len(k1))
	}
	if !bytes.Equal(k1, k2) {
		t.Fatalf("expected deterministic output, keys differ")
	}

	// Change the salt - the key should be different
	params.Salt = []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
	k3, err := crypto.DeriveKey("password", params)
	if err != nil {
		t.Fatalf("DeriveKey error: %v", err)
	}
	if bytes.Equal(k1, k3) {
		t.Fatalf("expected different keys for different salt")
	}
}
