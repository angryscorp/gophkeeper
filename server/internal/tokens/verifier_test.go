package tokens

import (
	"crypto/ed25519"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestVerifier_Verify(t *testing.T) {
	pubKey, privKey, _ := ed25519.GenerateKey(nil)
	verifier := NewVerifier(pubKey, "test-aud")

	validToken := func() string {
		claims := AccessClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   "user123",
				Audience:  jwt.ClaimStrings{"test-aud"},
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
		tokenStr, _ := token.SignedString(privKey)
		return tokenStr
	}

	tests := []struct {
		name    string
		token   string
		wantSub string
		wantErr bool
	}{
		{"valid token", validToken(), "user123", false},
		{"malformed token", "not.a.valid.jwt", "", true},
		{"empty token", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub, err := verifier.Verify(tt.token)

			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr %v, got error: %v", tt.wantErr, err)
			}
			if sub != tt.wantSub {
				t.Errorf("got subject %q, want %q", sub, tt.wantSub)
			}
		})
	}
}
