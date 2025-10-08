package domain

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestCredentials_ToRecord_Basic(t *testing.T) {
	id := uuid.New()
	creds := Credentials{
		ID:       id,
		Login:    "alice",
		Password: "s3cr3t!",
		Note:     "work account",
	}

	got := creds.ToRecord()

	if got.ID != id {
		t.Fatalf("ID mismatch: got %v want %v", got.ID, id)
	}
	if got.Title != creds.Note {
		t.Fatalf("Title mismatch: got %q want %q", got.Title, creds.Note)
	}
	if got.Kind != UserDataKindCredentials {
		t.Fatalf("Kind mismatch: got %v want %v", got.Kind, UserDataKindCredentials)
	}

	wantSensitive := fmt.Sprintf("Login: %s | Password: %s", creds.Login, creds.Password)
	if got.SensitiveInfo != wantSensitive {
		t.Fatalf("SensitiveInfo mismatch:\n got:  %q\n want: %q", got.SensitiveInfo, wantSensitive)
	}
}

func TestCredentials_ToRecord_EmptyFields(t *testing.T) {
	id := uuid.New()
	creds := Credentials{
		ID:       id,
		Login:    "",
		Password: "",
		Note:     "",
	}

	got := creds.ToRecord()

	if got.ID != id {
		t.Fatalf("ID mismatch: got %v want %v", got.ID, id)
	}
	if got.Title != "" {
		t.Fatalf("Title mismatch: got %q want empty", got.Title)
	}

	wantSensitive := "Login:  | Password: "
	if got.SensitiveInfo != wantSensitive {
		t.Fatalf("SensitiveInfo mismatch:\n got:  %q\n want: %q", got.SensitiveInfo, wantSensitive)
	}
}

func TestCredentials_ToRecord_WithUnicode(t *testing.T) {
	creds := Credentials{
		ID:       uuid.New(),
		Login:    "user@example.com",
		Password: "p@ss🔒",
		Note:     "personal",
	}

	got := creds.ToRecord()

	if got.Title != "personal" {
		t.Fatalf("Title mismatch: got %q want %q", got.Title, "personal")
	}

	wantSensitive := "Login: user@example.com | Password: p@ss🔒"
	if got.SensitiveInfo != wantSensitive {
		t.Fatalf("SensitiveInfo mismatch:\n got:  %q\n want: %q", got.SensitiveInfo, wantSensitive)
	}
}
