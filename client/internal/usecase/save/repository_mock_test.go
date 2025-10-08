package save

import (
	"context"

	"gophkeeper/client/internal/domain"

	"github.com/google/uuid"
)

type mockRepo struct {
	called bool
	kind   domain.UserDataKind
	id     uuid.UUID
	data   []byte
	err    error
}

func (m *mockRepo) Save(ctx context.Context, kind domain.UserDataKind, id uuid.UUID, data []byte) error {
	m.called, m.kind, m.id, m.data = true, kind, id, data
	return m.err
}
