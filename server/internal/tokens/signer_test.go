package tokens

import (
	"crypto/ed25519"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestSigner_IssueAccess(t *testing.T) {
	pubKey, privKey, _ := ed25519.GenerateKey(nil)
	signer := NewSigner(privKey, "test-aud", time.Hour)

	tests := []struct {
		name     string
		userID   string
		deviceID string
		wantErr  bool
	}{
		{"valid token", "user123", "device456", false},
		{"empty userID", "", "device456", false},
		{"empty deviceID", "user123", "", false},
		{"both empty", "", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenStr, err := signer.IssueAccess(tt.userID, tt.deviceID)

			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr %v, got error: %v", tt.wantErr, err)
			}

			if err == nil {
				// Parse and verify the token
				claims := &AccessClaims{}
				token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
					return pubKey, nil
				})
				if err != nil {
					t.Fatalf("failed to parse token: %v", err)
				}
				if !token.Valid {
					t.Error("token is not valid")
				}

				// Verify claims
				if claims.Subject != tt.userID {
					t.Errorf("got subject %q, want %q", claims.Subject, tt.userID)
				}
				if claims.DeviceID != tt.deviceID {
					t.Errorf("got deviceID %q, want %q", claims.DeviceID, tt.deviceID)
				}
				if len(claims.Audience) != 1 || claims.Audience[0] != "test-aud" {
					t.Errorf("got audience %v, want [test-aud]", claims.Audience)
				}
				if claims.IssuedAt == nil {
					t.Error("IssuedAt is nil")
				}
				if claims.ExpiresAt == nil {
					t.Error("ExpiresAt is nil")
				}
			}
		})
	}
}
