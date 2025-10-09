package tokens

import (
	"github.com/golang-jwt/jwt/v5"
)

// AccessClaims defines the JWT claims used in access tokens.
// It embeds standard registered claims and adds a custom DeviceID field.
type AccessClaims struct {
	jwt.RegisteredClaims
	DeviceID string `json:"device_id"`
}
