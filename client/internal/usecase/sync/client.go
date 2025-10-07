package sync

import (
	"context"
	"gophkeeper/client/internal/domain"

	"github.com/google/uuid"
)

// PullResponse is returned by Client.Pull. It contains a batch of changes
// from the server, the new cursor position, and a flag indicating if more
// data remains to be fetched.
type PullResponse struct {
	Changes   []domain.Message
	NewCursor int64
	HasMore   bool
}

// Client defines the remote sync API used by the Sync use case.
// It abstracts the gRPC/HTTP client that talks to the server.
type Client interface {
	Ping(ctx context.Context, accessToken string) error
	Push(ctx context.Context, accessToken string, messages []domain.Message) ([]uuid.UUID, error)
	Pull(ctx context.Context, accessToken string, cursor int64) (*PullResponse, error)
}
