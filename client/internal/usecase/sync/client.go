package sync

import (
	"context"
	"gophkeeper/client/internal/domain"

	"github.com/google/uuid"
)

type PullResponse struct {
	Changes   []domain.Message
	NewCursor int64
	HasMore   bool
}

type Client interface {
	Ping(ctx context.Context, accessToken string) error
	Push(ctx context.Context, accessToken string, messages []domain.Message) ([]uuid.UUID, error)
	Pull(ctx context.Context, accessToken string, cursor int64) (*PullResponse, error)
}
