package list

import (
	"context"
	"gophkeeper/client/internal/domain"

	"github.com/google/uuid"
)

type RawRecord struct {
	ID      uuid.UUID
	Kind    domain.UserDataKind
	Payload []byte
}

type Repository interface {
	GetAll(ctx context.Context) ([]RawRecord, error)
}
