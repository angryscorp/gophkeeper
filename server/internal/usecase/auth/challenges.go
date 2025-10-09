package auth

import (
	"context"
	"time"

	"gophkeeper/pkg/crypto"

	"github.com/google/uuid"
)

// ChallengeInfo holds details about a stored login challenge.
// It is passed to the validator callback in order to check correctness.
type ChallengeInfo struct {
	Challenge        []byte
	Attempts         int32
	AuthKey          []byte
	AuthKeyAlgorithm crypto.AuthKeyAlgorithm
}

// Challenges defines the persistence and validation interface for login challenges.
// Implementations are responsible for storing, updating and atomically verifying challenges.
type Challenges interface {
	Add(ctx context.Context, userId uuid.UUID, deviceName string, challenge []byte, expiresAt time.Time) error
	GetForUpdate(ctx context.Context, username, deviceName string, challengerValidator func(ChallengeInfo) bool) error
}
