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
