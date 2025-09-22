package challenges

import (
	"context"
	"gophkeeper/pkg/crypto"
	"time"

	"github.com/google/uuid"
)

const (
	ChallengeNotFound      = "challenge not found"
	WrongChallengeResponse = "wrong challenge response"
)

type ChallengeInfo struct {
	Challenge        []byte
	Attempts         int32
	AuthKey          []byte
	AuthKeyAlgorithm crypto.AuthKeyAlgorithm
}

type Challenges interface {
	Add(ctx context.Context, userId uuid.UUID, deviceName string, challenge []byte, expiresAt time.Time) error
	GetForUpdate(ctx context.Context, username, deviceName string, challengerValidator func(ChallengeInfo) bool) error
}
