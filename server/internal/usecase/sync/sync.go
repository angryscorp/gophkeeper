package sync

import (
	"context"
	"gophkeeper/server/internal/domain"

	"github.com/google/uuid"
)

type Sync struct{}

func New() *Sync {
	return &Sync{}
}

type PullResponse struct {
	Changes    []domain.Message
	NextCursor int64
	HasMore    bool
}

func (s *Sync) Pull(ctx context.Context, cursor int64, limit int32) (*PullResponse, error) {
	return &PullResponse{
		Changes:    []domain.Message{},
		NextCursor: 0,
		HasMore:    false,
	}, nil
}

func (s *Sync) Push(ctx context.Context, changes []domain.Message) ([]uuid.UUID, error) {
	res := make([]uuid.UUID, len(changes))
	for i, change := range changes {
		res[i] = change.ID
	}
	return res, nil
}
