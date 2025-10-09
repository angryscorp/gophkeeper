package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type stubChallenges struct {
	addFn          func(ctx context.Context, userId uuid.UUID, deviceName string, challenge []byte, expiresAt time.Time) error
	getForUpdateFn func(ctx context.Context, username, deviceName string, validator func(ChallengeInfo) bool) error
}

func (s stubChallenges) Add(ctx context.Context, userId uuid.UUID, deviceName string, challenge []byte, expiresAt time.Time) error {
	if s.addFn != nil {
		return s.addFn(ctx, userId, deviceName, challenge, expiresAt)
	}
	return nil
}

func (s stubChallenges) GetForUpdate(ctx context.Context, username, deviceName string, validator func(ChallengeInfo) bool) error {
	if s.getForUpdateFn != nil {
		return s.getForUpdateFn(ctx, username, deviceName, validator)
	}
	return nil
}
