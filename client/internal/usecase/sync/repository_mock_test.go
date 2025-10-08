package sync

import (
	"context"
	"errors"

	"gophkeeper/client/internal/domain"

	"github.com/google/uuid"
)

type mockRepo struct {
	// cursor
	cursor       int64
	getCursorErr error

	// SaveChanges
	saveCalls []struct {
		changes   []domain.Message
		newCursor int64
	}
	saveErrAt int // if >0, fail on that call index (1-based)

	// Outbox batches (for Push)
	batches   [][]domain.Message // returned sequentially, then empty
	batchCall int

	// Delete
	lastDeleted []uuid.UUID
	deleteErr   error
}

func (m *mockRepo) GetCursor(ctx context.Context) (int64, error) {
	if m.getCursorErr != nil {
		return 0, m.getCursorErr
	}
	// emulate advancing cursor by last saved call if present
	if n := len(m.saveCalls); n > 0 {
		return m.saveCalls[n-1].newCursor, nil
	}
	return m.cursor, nil
}

func (m *mockRepo) SaveChanges(ctx context.Context, changes []domain.Message, newCursor int64) error {
	callNo := len(m.saveCalls) + 1
	if m.saveErrAt > 0 && callNo == m.saveErrAt {
		return errors.New("save err")
	}
	cp := append([]domain.Message(nil), changes...)
	m.saveCalls = append(m.saveCalls, struct {
		changes   []domain.Message
		newCursor int64
	}{changes: cp, newCursor: newCursor})
	return nil
}

func (m *mockRepo) GetBatch(ctx context.Context, limit int64) ([]domain.Message, error) {
	if m.batchCall >= len(m.batches) {
		return nil, nil
	}
	b := m.batches[m.batchCall]
	m.batchCall++
	cp := append([]domain.Message(nil), b...)
	return cp, nil
}

func (m *mockRepo) DeleteBatch(ctx context.Context, ids []uuid.UUID) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	m.lastDeleted = append([]uuid.UUID(nil), ids...)
	return nil
}
