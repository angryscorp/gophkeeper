package crypto_test

import (
	"bytes"
	"testing"

	"gophkeeper/pkg/crypto"
)

func TestRandBytes_LengthAndRandomness(t *testing.T) {
	const n = 32

	b1 := crypto.RandBytes(n)
	b2 := crypto.RandBytes(n)

	if len(b1) != n {
		t.Errorf("expected length %d, got %d", n, len(b1))
	}
	if len(b2) != n {
		t.Errorf("expected length %d, got %d", n, len(b2))
	}
	if bytes.Equal(b1, b2) {
		t.Errorf("expected different random values, got identical slices: %x", b1)
	}
}
