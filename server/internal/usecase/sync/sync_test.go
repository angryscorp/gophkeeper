package sync

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"gophkeeper/server/internal/domain"

	"github.com/google/uuid"
)

func TestNew_ReturnsService(t *testing.T) {
	svc := New(&stubRepo{})
	if svc == nil {
		t.Fatalf("New() returned nil")
	}
}

func TestPull_Success_DelegatesAndReturns(t *testing.T) {
	want := &PullResponse{
		Changes: []domain.Message{
			{ID: uuid.New(), RecordID: uuid.New(), Kind: 1, UpdatedAtUnix: 123},
		},
		NextCursor: 42,
		HasMore:    true,
	}

	stub := &stubRepo{
		getChangesFn: func(ctx context.Context, username string, cursor int64, limit int32) (*PullResponse, error) {
			if username != "alice" || cursor != 10 || limit != 50 {
				t.Fatalf("unexpected args: user=%s cursor=%d limit=%d", username, cursor, limit)
			}
			return want, nil
		},
	}

	svc := New(stub)
	got, err := svc.Pull(context.Background(), "alice", 10, 50)
	if err != nil {
		t.Fatalf("Pull() unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Pull() mismatch:\nwant: %#v\ngot:  %#v", want, got)
	}
}

func TestPull_Error_Propagates(t *testing.T) {
	sentinel := errors.New("repo down")
	stub := &stubRepo{
		getChangesFn: func(ctx context.Context, username string, cursor int64, limit int32) (*PullResponse, error) {
			return nil, sentinel
		},
	}
	svc := New(stub)
	_, err := svc.Pull(context.Background(), "u", 0, 1)
	if !errors.Is(err, sentinel) {
		t.Fatalf("Pull() want %v, got %v", sentinel, err)
	}
}

func TestPush_Success_DelegatesAndReturns(t *testing.T) {
	changes := []domain.Message{
		{ID: uuid.New(), RecordID: uuid.New(), Kind: 2, UpdatedAtUnix: 999},
	}
	wantIDs := []uuid.UUID{changes[0].RecordID}

	stub := &stubRepo{
		addChangesFn: func(ctx context.Context, username string, got []domain.Message) ([]uuid.UUID, error) {
			if username != "bob" {
				t.Fatalf("unexpected username: %s", username)
			}
			if !reflect.DeepEqual(got, changes) {
				t.Fatalf("AddChanges() got changes mismatch")
			}
			return wantIDs, nil
		},
	}

	svc := New(stub)
	got, err := svc.Push(context.Background(), "bob", changes)
	if err != nil {
		t.Fatalf("Push() unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, wantIDs) {
		t.Fatalf("Push() ids mismatch: want %v, got %v", wantIDs, got)
	}
}

func TestPush_Error_Propagates(t *testing.T) {
	sentinel := errors.New("write failed")
	stub := &stubRepo{
		addChangesFn: func(ctx context.Context, username string, changes []domain.Message) ([]uuid.UUID, error) {
			return nil, sentinel
		},
	}
	svc := New(stub)
	_, err := svc.Push(context.Background(), "u", nil)
	if !errors.Is(err, sentinel) {
		t.Fatalf("Push() want %v, got %v", sentinel, err)
	}
}
