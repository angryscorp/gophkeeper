package tokens

import (
	"crypto/ed25519"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Signer issues signed JWT access tokens using an Ed25519 private key.
// Each token includes standard claims (subject, audience, issued-at, expiry)
// and a custom DeviceID claim.
type Signer struct {
	privateKey ed25519.PrivateKey
	aud        string
	ttl        time.Duration
}

// NewSigner creates a new Signer with the given Ed25519 private key,
// target audience, and token time-to-live (TTL).
func NewSigner(privateKey ed25519.PrivateKey, aud string, ttl time.Duration) *Signer {
	return &Signer{privateKey: privateKey, aud: aud, ttl: ttl}
}

// IssueAccess creates and signs a new JWT access token for the given user and device.
// The token subject is set to userID, and a custom DeviceID claim is added.
// The token is valid from now until now+ttl.
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
