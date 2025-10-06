package sync

import (
	"context"
	"gophkeeper/server/internal/domain"

	"github.com/google/uuid"
)

type Sync struct {
	repo Repository
}

func New(repo Repository) *Sync {
	return &Sync{repo: repo}
}

func (s *Sync) Pull(ctx context.Context, username string, cursor int64, limit int32) (*PullResponse, error) {
	return s.repo.GetChanges(ctx, username, cursor, limit)
}

func (s *Sync) Push(ctx context.Context, changes []domain.Message) ([]uuid.UUID, error) {
	res := make([]uuid.UUID, len(changes))
	for i, change := range changes {
		res[i] = change.ID
	}
	return res, nil
}
