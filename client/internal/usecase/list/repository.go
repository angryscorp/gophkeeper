package list

import (
	"context"
	"gophkeeper/client/internal/domain"

	"github.com/google/uuid"
)

// RawRecord is the low-level encrypted record stored in the repository.
// Payload contains encrypted JSON for a specific UserDataKind.
type RawRecord struct {
	ID      uuid.UUID
	Kind    domain.UserDataKind
	Payload []byte
}

// Repository defines access to stored raw records (encrypted).
// It abstracts the local storage (e.g. SQLite) behind a simple API.
type Repository interface {
	GetAll(ctx context.Context) ([]RawRecord, error)
}
