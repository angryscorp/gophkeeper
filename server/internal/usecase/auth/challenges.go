package auth

import (
	"context"
	"gophkeeper/pkg/crypto"
	"time"

	"github.com/google/uuid"
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
