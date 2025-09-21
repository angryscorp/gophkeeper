package crypto

import (
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/hkdf"
)

// DeriveKey derives a key from a password using the specified KDF parameters.
func DeriveKey(password string, params KDFParameters) ([]byte, error) {
	switch params.Algorithm {
	case KDFAlgorithmARGON2ID:
		return argon2.IDKey(
			[]byte(password),
			params.Salt,
			params.TimeCost,
			params.MemoryCost,
			uint8(params.Parallelism),
			32,
		), nil
	default:
		return nil, fmt.Errorf("unsupported KDF algorithm: %v", params.Algorithm)
	}
}

// DeriveAuthKey derives an authentication key from a key and an info string.
func DeriveAuthKey(key, info []byte) ([]byte, error) {
	h := hkdf.New(sha256.New, key, nil, info)

	authKey := make([]byte, 32)
	if _, err := io.ReadFull(h, authKey); err != nil {
		return nil, fmt.Errorf("failed to derive auth key: %w", err)
	}

	return authKey, nil
}
