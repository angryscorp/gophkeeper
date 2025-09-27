package tokens

import (
	"crypto/ed25519"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Verifier struct {
	publicKey ed25519.PublicKey
	aud       string
}

func NewVerifier(publicKey ed25519.PublicKey, aud string) *Verifier {
	return &Verifier{publicKey: publicKey, aud: aud}
}

func (v *Verifier) ParseAndVerify(tokenStr string) (*Claims, error) {
	ac := &AccessClaims{}

	tok, err := jwt.ParseWithClaims(
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

	if err != nil || !tok.Valid {
		return nil, errors.New("invalid jwt")
	}

	return &Claims{
		Sub:      ac.Subject,
		DeviceID: ac.DeviceID,
		Exp:      ac.ExpiresAt.Time,
	}, nil
}
