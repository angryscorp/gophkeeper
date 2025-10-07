package save

import (
	"context"
	"gophkeeper/client/internal/domain"

	"github.com/google/uuid"
)

// Repository defines how encrypted user records are persisted.
type Repository interface {
	Save(ctx context.Context, kind domain.UserDataKind, id uuid.UUID, data []byte) error
}
