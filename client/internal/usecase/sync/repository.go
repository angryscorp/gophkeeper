package sync

import (
	"context"

	"gophkeeper/client/internal/domain"

	"github.com/google/uuid"
)

// Repository defines the local persistence needed for synchronization.
// It manages the outbox (pending changes), the cursor for pull progress,
// and the application of server-side changes.
type Repository interface {
	GetBatch(ctx context.Context, limit int64) ([]domain.Message, error)
	DeleteBatch(ctx context.Context, ids []uuid.UUID) error
	GetCursor(ctx context.Context) (int64, error)
	SaveChanges(ctx context.Context, changes []domain.Message, newCursor int64) error
}
