package tokens

import (
	"crypto/ed25519"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Signer struct {
	privateKey ed25519.PrivateKey
	aud        string
	ttl        time.Duration
}

func NewSigner(privateKey ed25519.PrivateKey, aud string, ttl time.Duration) *Signer {
	return &Signer{privateKey: privateKey, aud: aud, ttl: ttl}
}

func (s *Signer) IssueAccess(userID, deviceID string) (string, error) {
	now := time.Now().UTC()
	claims := AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Audience:  []string{s.aud},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.ttl)),
		},
		DeviceID: deviceID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	return token.SignedString(s.privateKey)
}
