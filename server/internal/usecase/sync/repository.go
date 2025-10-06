package sync

import (
	"context"
)

type Repository interface {
	GetChanges(ctx context.Context, username string, cursor int64, limit int32) (*PullResponse, error)
}
