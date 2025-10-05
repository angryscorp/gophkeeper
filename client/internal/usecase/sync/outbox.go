package sync

import (
	"context"
	"gophkeeper/client/internal/domain"

	"github.com/google/uuid"
)

type OutboxRepository interface {
	GetBatch(ctx context.Context, limit int64) ([]domain.OutboxMessage, error)
	DeleteBatch(ctx context.Context, ids []uuid.UUID) error
}
