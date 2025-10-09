package sync

import (
	"context"

	"gophkeeper/server/internal/domain"

	"github.com/google/uuid"
)

// Repository defines the persistence layer for synchronization.
type Repository interface {
	GetChanges(ctx context.Context, username string, cursor int64, limit int32) (*PullResponse, error)
	AddChanges(ctx context.Context, username string, changes []domain.Message) ([]uuid.UUID, error)
}
