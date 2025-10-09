package sync

import (
	"context"
	"errors"

	"gophkeeper/client/internal/domain"

	"github.com/google/uuid"
)

type mockClient struct {
	// Ping
	pingErr error

	// Pull
	pullSeq   []*PullResponse // returned sequentially per call
	pullErrAt int             // if >0, return error on that call index (1-based)
	pullCalls int

	// Push
	pushErr           error
	lastPushToken     string
	lastPushMsgs      []domain.Message
	returnIDsFromMsgs bool // if true, returns IDs extracted from messages
	returnIDsFixed    []uuid.UUID
}

func (m *mockClient) Ping(ctx context.Context, accessToken string) error {
	return m.pingErr
}

func (m *mockClient) Pull(ctx context.Context, accessToken string, cursor int64) (*PullResponse, error) {
	m.pullCalls++
	if m.pullErrAt > 0 && m.pullCalls == m.pullErrAt {
		return nil, errors.New("pull err")
	}
	// defensive: if out of range, return empty page
	if m.pullCalls-1 >= len(m.pullSeq) {
		return &PullResponse{Changes: nil, NewCursor: cursor, HasMore: false}, nil
	}
	return m.pullSeq[m.pullCalls-1], nil
}

func (m *mockClient) Push(ctx context.Context, accessToken string, messages []domain.Message) ([]uuid.UUID, error) {
	m.lastPushToken = accessToken
	m.lastPushMsgs = append([]domain.Message(nil), messages...)
	if m.pushErr != nil {
		return nil, m.pushErr
	}
	if m.returnIDsFromMsgs {
		out := make([]uuid.UUID, 0, len(messages))
		for _, msg := range messages {
			out = append(out, msg.ID)
		}
		return out, nil
	}
	return append([]uuid.UUID(nil), m.returnIDsFixed...), nil
}
