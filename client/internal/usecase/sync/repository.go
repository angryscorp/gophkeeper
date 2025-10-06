package sync

import (
	"context"
	"gophkeeper/client/internal/domain"

	"github.com/google/uuid"
)

type Repository interface {
	GetBatch(ctx context.Context, limit int64) ([]domain.Message, error)
	DeleteBatch(ctx context.Context, ids []uuid.UUID) error
	GetCursor(ctx context.Context) (int64, error)
	SaveChanges(ctx context.Context, changes []domain.Message, newCursor int64) error
}
