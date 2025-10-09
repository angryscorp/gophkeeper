package sync

import (
	"context"

	"gophkeeper/server/internal/domain"

	"github.com/google/uuid"
)

// Sync provides the high-level synchronization use cases
// by delegating persistence operations to the Repository.
type Sync struct {
	repo Repository
}

// New constructs a Sync service backed by the given repository.
func New(repo Repository) *Sync {
	return &Sync{repo: repo}
}

// Pull retrieves a batch of changes for the given user since the provided cursor.
func (s *Sync) Pull(ctx context.Context, username string, cursor int64, limit int32) (*PullResponse, error) {
	return s.repo.GetChanges(ctx, username, cursor, limit)
}

// Push persists a batch of client-side changes for the given user.
// Returns the list of successfully applied record IDs.
func (s *Sync) Push(ctx context.Context, username string, changes []domain.Message) ([]uuid.UUID, error) {
	return s.repo.AddChanges(ctx, username, changes)
}
