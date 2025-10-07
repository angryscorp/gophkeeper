package tokens

import (
	"crypto/ed25519"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

// Verifier verifies JWT access tokens signed with Ed25519 and validates audience and expiry.
type Verifier struct {
	publicKey ed25519.PublicKey
	aud       string
}

// NewVerifier creates a new Verifier with the given Ed25519 public key and expected audience.
func NewVerifier(publicKey ed25519.PublicKey, aud string) *Verifier {
	return &Verifier{publicKey: publicKey, aud: aud}
}

// Verify parses and validates the given JWT string.
// It checks signature, algorithm, audience, issued-at, and expiration claims.
// Returns the subject (user identifier) if valid.
func (v *Verifier) Verify(tokenStr string) (string, error) {
	ac := &AccessClaims{}
	token, err := jwt.ParseWithClaims(
		tokenStr,
		ac,
		func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodEdDSA {
				return nil, fmt.Errorf("unexpected alg: %v", t.Header["alg"])
			}
			return v.publicKey, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}),
		jwt.WithAudience(v.aud),
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	return token.Claims.GetSubject()
}
