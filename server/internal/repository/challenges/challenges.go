package challenges

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	ChallengeNotFound      = "challenge not found"
	WrongChallengeResponse = "wrong challenge response"
)

type Challenges interface {
	Add(ctx context.Context, userId uuid.UUID, deviceName string, challenge []byte, expiresAt time.Time) error
}
