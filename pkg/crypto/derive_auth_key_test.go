package crypto_test

import (
	"bytes"
	"testing"

	"gophkeeper/pkg/crypto"
)

func TestDeriveAuthKey_DeterministicAndLength(t *testing.T) {
	key := []byte("32-bytes-master-key--------------") // 32 bytes
	info := []byte("auth")

	k1, err := crypto.DeriveAuthKey(key, info)
	if err != nil {
		t.Fatalf("DeriveAuthKey error: %v", err)
	}
	k2, err := crypto.DeriveAuthKey(key, info)
	if err != nil {
		t.Fatalf("DeriveAuthKey error: %v", err)
	}

	if len(k1) != 32 {
		t.Fatalf("expected auth key length 32, got %d", len(k1))
	}
	if !bytes.Equal(k1, k2) {
		t.Fatalf("expected deterministic output, keys differ")
	}

	// Change the info - the key must be different
	k3, err := crypto.DeriveAuthKey(key, []byte("different-info"))
	if err != nil {
		t.Fatalf("DeriveAuthKey error: %v", err)
	}
	if bytes.Equal(k1, k3) {
		t.Fatalf("expected different auth keys for different info")
	}
}
