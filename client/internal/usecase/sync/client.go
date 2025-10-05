package sync

import (
	"context"
	"gophkeeper/client/internal/domain"

	"github.com/google/uuid"
)

type Client interface {
	Ping(ctx context.Context, accessToken string) error
	Push(ctx context.Context, accessToken string, messages []domain.OutboxMessage) ([]uuid.UUID, error)
}
