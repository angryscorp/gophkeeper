package sync

import (
	"errors"
	"testing"

	"gophkeeper/client/internal/domain"

	"github.com/google/uuid"
)

func TestPing_Success(t *testing.T) {
	tokens := &mockTokens{ready: true, token: "tok"}
	client := &mockClient{}
	repo := &mockRepo{}
	s := New(client, tokens, repo)

	if err := s.Ping(); err != nil {
		t.Fatalf("Ping() unexpected error: %v", err)
	}
}

func TestPing_NotReady(t *testing.T) {
	tokens := &mockTokens{ready: false}
	client := &mockClient{}
	repo := &mockRepo{}
	s := New(client, tokens, repo)

	if err := s.Ping(); err == nil {
		t.Fatalf("Ping() expected error when not ready")
	}
}

func TestPing_TokenErr(t *testing.T) {
	tokens := &mockTokens{ready: true, err: errors.New("no token")}
	client := &mockClient{}
	repo := &mockRepo{}
	s := New(client, tokens, repo)

	if err := s.Ping(); err == nil {
		t.Fatalf("Ping() expected token error")
	}
}

func TestPing_ClientErr(t *testing.T) {
	tokens := &mockTokens{ready: true, token: "tok"}
	client := &mockClient{pingErr: errors.New("ping fail")}
	repo := &mockRepo{}
	s := New(client, tokens, repo)

	if err := s.Ping(); err == nil {
		t.Fatalf("Ping() expected client error")
	}
}

func TestPull_MultiPages_Success(t *testing.T) {
	tokens := &mockTokens{ready: true, token: "tok"}
	// two pages: first hasMore=true, second false
	msg1 := domain.Message{ID: uuid.New(), UpdatedAtUnix: 1}
	msg2 := domain.Message{ID: uuid.New(), UpdatedAtUnix: 2}
	client := &mockClient{
		pullSeq: []*PullResponse{
			{Changes: []domain.Message{msg1}, NewCursor: 10, HasMore: true},
			{Changes: []domain.Message{msg2}, NewCursor: 20, HasMore: false},
		},
	}
	repo := &mockRepo{cursor: 0}
	s := New(client, tokens, repo)

	if err := s.Pull(); err != nil {
		t.Fatalf("Pull() unexpected error: %v", err)
	}

	if len(repo.saveCalls) != 2 {
		t.Fatalf("expected 2 SaveChanges calls, got %d", len(repo.saveCalls))
	}
	if repo.saveCalls[0].newCursor != 10 || repo.saveCalls[1].newCursor != 20 {
		t.Errorf("unexpected cursors saved: %+v", repo.saveCalls)
	}
	if client.pullCalls != 2 {
		t.Errorf("expected 2 pulls, got %d", client.pullCalls)
	}
}

func TestPull_Error_GetCursor(t *testing.T) {
	tokens := &mockTokens{ready: true, token: "tok"}
	client := &mockClient{}
	repo := &mockRepo{getCursorErr: errors.New("cursor fail")}
	s := New(client, tokens, repo)

	if err := s.Pull(); err == nil {
		t.Fatalf("expected error from GetCursor")
	}
}

func TestPull_Error_ClientPull(t *testing.T) {
	tokens := &mockTokens{ready: true, token: "tok"}
	client := &mockClient{
		pullSeq:   []*PullResponse{{Changes: nil, NewCursor: 1, HasMore: false}},
		pullErrAt: 1, // fail on first call
	}
	repo := &mockRepo{}
	s := New(client, tokens, repo)

	if err := s.Pull(); err == nil {
		t.Fatalf("expected error from client.Pull")
	}
}

func TestPull_Error_SaveChanges(t *testing.T) {
	tokens := &mockTokens{ready: true, token: "tok"}
	client := &mockClient{
		pullSeq: []*PullResponse{
			{Changes: []domain.Message{{ID: uuid.New()}}, NewCursor: 5, HasMore: false},
		},
	}
	repo := &mockRepo{saveErrAt: 1}
	s := New(client, tokens, repo)

	if err := s.Pull(); err == nil {
		t.Fatalf("expected error from SaveChanges")
	}
}

func TestPush_Batches_Success(t *testing.T) {
	tokens := &mockTokens{ready: true, token: "tok"}

	// two batches, then empty
	m1 := domain.Message{ID: uuid.New()}
	m2 := domain.Message{ID: uuid.New()}
	m3 := domain.Message{ID: uuid.New()}
	repo := &mockRepo{
		batches: [][]domain.Message{
			{m1, m2},
			{m3},
		},
	}

	client := &mockClient{
		returnIDsFromMsgs: true,
	}

	s := New(client, tokens, repo)

	if err := s.Push(); err != nil {
		t.Fatalf("Push() unexpected error: %v", err)
	}

	// After last batch, DeleteBatch must be called with the IDs returned by server
	if len(repo.lastDeleted) != 1 || repo.lastDeleted[0] != m3.ID {
		t.Fatalf("expected last delete to contain [%s], got %v", m3.ID, repo.lastDeleted)
	}
}

func TestPush_Error_ClientPush(t *testing.T) {
	tokens := &mockTokens{ready: true, token: "tok"}
	m1 := domain.Message{ID: uuid.New()}
	repo := &mockRepo{
		batches: [][]domain.Message{
			{m1},
		},
	}
	client := &mockClient{
		pushErr: errors.New("push fail"),
	}
	s := New(client, tokens, repo)

	if err := s.Push(); err == nil {
		t.Fatalf("expected error from client.Push")
	}
}

func TestPush_Error_DeleteBatch(t *testing.T) {
	tokens := &mockTokens{ready: true, token: "tok"}
	m1 := domain.Message{ID: uuid.New()}
	repo := &mockRepo{
		batches:   [][]domain.Message{{m1}},
		deleteErr: errors.New("delete fail"),
	}
	client := &mockClient{returnIDsFromMsgs: true}
	s := New(client, tokens, repo)

	if err := s.Push(); err == nil {
		t.Fatalf("expected error from DeleteBatch")
	}
}

func TestPush_NotReady(t *testing.T) {
	tokens := &mockTokens{ready: false}
	repo := &mockRepo{}
	client := &mockClient{}
	s := New(client, tokens, repo)

	if err := s.Push(); err == nil {
		t.Fatalf("expected error when not ready")
	}
}

func TestPull_NotReady(t *testing.T) {
	tokens := &mockTokens{ready: false}
	repo := &mockRepo{}
	client := &mockClient{}
	s := New(client, tokens, repo)

	if err := s.Pull(); err == nil {
		t.Fatalf("expected error when not ready")
	}
}
