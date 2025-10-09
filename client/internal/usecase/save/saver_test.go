package save

import (
	"errors"
	"testing"

	"gophkeeper/client/internal/domain"

	"github.com/google/uuid"
)

func TestSaveCredentials_Success(t *testing.T) {
	r := &mockRepo{}
	saver := New(r, goodEncryptor)

	c := domain.Credentials{ID: uuid.New(), Login: "alice", Password: "pw"}
	if err := saver.SaveCredentials(c); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !r.called || r.kind != domain.UserDataKindCredentials {
		t.Errorf("repo not called with expected kind")
	}
	// ensure encryptor added prefix
	if string(r.data[:2]) != "X:" {
		t.Errorf("expected encrypted prefix, got %s", r.data)
	}
}

func TestSaveBankCard_ErrorFromEncryptor(t *testing.T) {
	r := &mockRepo{}
	saver := New(r, badEncryptor)

	card := domain.BankCard{ID: uuid.New(), Owner: "bob"}
	err := saver.SaveBankCard(card)
	if err == nil || err.Error() != "encrypt fail" {
		t.Fatalf("expected encrypt error, got %v", err)
	}
	if r.called {
		t.Error("repo.Save should not be called on encrypt error")
	}
}

func TestSaveUserBinaryData_ErrorFromRepo(t *testing.T) {
	wantErr := errors.New("db fail")
	r := &mockRepo{err: wantErr}
	saver := New(r, goodEncryptor)

	bin := domain.UserBinaryData{ID: uuid.New(), Data: []byte{1, 2, 3}}
	err := saver.SaveUserBinaryData(bin)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected repo error, got %v", err)
	}
}

func goodEncryptor(b []byte) ([]byte, error) { return append([]byte("X:"), b...), nil }
func badEncryptor([]byte) ([]byte, error)    { return nil, errors.New("encrypt fail") }
