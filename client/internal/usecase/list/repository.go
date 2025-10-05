package list

import (
	"context"
	"gophkeeper/client/internal/domain"

	"github.com/google/uuid"
)

type RawRecord struct {
	Kind domain.UserDataKind
	ID   uuid.UUID
	Data []byte
}

type Repository interface {
	GetAll(ctx context.Context) ([]RawRecord, error)
}
