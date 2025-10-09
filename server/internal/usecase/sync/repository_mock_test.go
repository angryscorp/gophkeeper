package sync

import (
	"context"

	"gophkeeper/server/internal/domain"

	"github.com/google/uuid"
)

type stubRepo struct {
	getChangesFn func(ctx context.Context, username string, cursor int64, limit int32) (*PullResponse, error)
	addChangesFn func(ctx context.Context, username string, changes []domain.Message) ([]uuid.UUID, error)

	// captured args (to assert)
	lastUser   string
	lastCursor int64
	lastLimit  int32
	lastPush   []domain.Message
}

func (s *stubRepo) GetChanges(ctx context.Context, username string, cursor int64, limit int32) (*PullResponse, error) {
	s.lastUser, s.lastCursor, s.lastLimit = username, cursor, limit
	if s.getChangesFn != nil {
		return s.getChangesFn(ctx, username, cursor, limit)
	}
	return &PullResponse{}, nil
}

func (s *stubRepo) AddChanges(ctx context.Context, username string, changes []domain.Message) ([]uuid.UUID, error) {
	s.lastUser = username
	s.lastPush = changes
	if s.addChangesFn != nil {
		return s.addChangesFn(ctx, username, changes)
	}
	return nil, nil
}
