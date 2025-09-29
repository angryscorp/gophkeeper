package tokens

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AccessClaims struct {
	jwt.RegisteredClaims
	DeviceID string `json:"device_id"`
}

type Claims struct {
	Sub      string // user_id (UUID as string)
	DeviceID string
	Exp      time.Time
}
