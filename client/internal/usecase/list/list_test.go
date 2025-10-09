package list

import (
	"encoding/json"
	"errors"
	"testing"

	"gophkeeper/client/internal/domain"

	"github.com/google/uuid"
)

func TestGetAllRecords_Success_MultipleKinds(t *testing.T) {
	// prepare 4 kinds
	card := domain.BankCard{ID: uuid.New(), Owner: "John", Number: "4111", CVV: "123", ExpireDate: "12/30", Note: "My Visa"}
	creds := domain.Credentials{ID: uuid.New(), Login: "alice", Password: "secret", Note: "GitHub"}
	text := domain.UserTextData{ID: uuid.New(), Data: "hello", Note: "Note A"}
	bin := domain.UserBinaryData{ID: uuid.New(), Data: []byte{1, 2, 3}, Note: "Photo"}

	repo := mockRepo{
		rows: []RawRecord{
			{ID: card.ID, Kind: domain.UserDataKindBankCard, Payload: mustJSON(t, card)},
			{ID: creds.ID, Kind: domain.UserDataKindCredentials, Payload: mustJSON(t, creds)},
			{ID: text.ID, Kind: domain.UserDataKindTextData, Payload: mustJSON(t, text)},
			{ID: bin.ID, Kind: domain.UserDataKindBinaryData, Payload: mustJSON(t, bin)},
		},
	}
	svc := New(repo, passthrough)

	got, err := svc.GetAllRecords()
	if err != nil {
		t.Fatalf("GetAllRecords() error = %v", err)
	}
	if len(got) != 4 {
		t.Fatalf("expected 4 records, got %d", len(got))
	}

	// spot-check kinds and titles
	tests := []struct {
		i    int
		kind domain.UserDataKind
		want string
	}{
		{0, domain.UserDataKindBankCard, card.Note},
		{1, domain.UserDataKindCredentials, creds.Note},
		{2, domain.UserDataKindTextData, text.Note},
		{3, domain.UserDataKindBinaryData, bin.Note},
	}
	for _, tc := range tests {
		if got[tc.i].Kind != tc.kind {
			t.Errorf("record %d kind = %v, want %v", tc.i, got[tc.i].Kind, tc.kind)
		}
		if got[tc.i].Title != tc.want {
			t.Errorf("record %d title = %q, want %q", tc.i, got[tc.i].Title, tc.want)
		}
	}
}

func TestGetAllRecords_RepoError(t *testing.T) {
	repo := mockRepo{err: errors.New("boom")}
	svc := New(repo, passthrough)

	_, err := svc.GetAllRecords()
	if err == nil {
		t.Fatal("expected error from repo, got nil")
	}
}

func TestGetAllRecords_DecryptError(t *testing.T) {
	row := RawRecord{
		ID:      uuid.New(),
		Kind:    domain.UserDataKindTextData,
		Payload: []byte("anything"),
	}
	repo := mockRepo{rows: []RawRecord{row}}
	decryptErr := errors.New("decrypt failed")
	svc := New(repo, func([]byte) ([]byte, error) { return nil, decryptErr })

	_, err := svc.GetAllRecords()
	if err == nil || !errors.Is(err, decryptErr) {
		t.Fatalf("expected decrypt error, got %v", err)
	}
}

func TestGetAllRecords_BadJSON(t *testing.T) {
	// invalid JSON after "decryption"
	row := RawRecord{
		ID:      uuid.New(),
		Kind:    domain.UserDataKindCredentials,
		Payload: []byte("{not-json"),
	}
	repo := mockRepo{rows: []RawRecord{row}}
	svc := New(repo, passthrough)

	_, err := svc.GetAllRecords()
	if err == nil {
		t.Fatal("expected JSON unmarshal error, got nil")
	}
}

func TestDecryptRecord_UnknownKind(t *testing.T) {
	repo := mockRepo{}
	svc := New(repo, passthrough)

	_, err := svc.decryptRecord(RawRecord{
		ID:      uuid.New(),
		Kind:    999, // unknown
		Payload: []byte("{}"),
	})
	if err == nil {
		t.Fatal("expected error for unknown data kind, got nil")
	}
}

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}

// identity "decryptor" for happy paths
func passthrough(b []byte) ([]byte, error) { return b, nil }
