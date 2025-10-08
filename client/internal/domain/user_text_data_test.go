package domain

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestUserTextData_ToRecord_Basic(t *testing.T) {
	id := uuid.New()
	utd := UserTextData{
		ID:   id,
		Data: "hello world",
		Note: "sample text",
	}

	got := utd.ToRecord()

	if got.ID != id {
		t.Fatalf("ID mismatch: got %v want %v", got.ID, id)
	}
	if got.Title != utd.Note {
		t.Fatalf("Title mismatch: got %q want %q", got.Title, utd.Note)
	}
	if got.Kind != UserDataKindTextData {
		t.Fatalf("Kind mismatch: got %v want %v", got.Kind, UserDataKindTextData)
	}

	wantSensitive := fmt.Sprintf("Data: %s", utd.Data)
	if got.SensitiveInfo != wantSensitive {
		t.Fatalf("SensitiveInfo mismatch:\n got: %q\n want: %q", got.SensitiveInfo, wantSensitive)
	}
}

func TestUserTextData_ToRecord_EmptyFields(t *testing.T) {
	id := uuid.New()
	utd := UserTextData{
		ID:   id,
		Data: "",
		Note: "",
	}

	got := utd.ToRecord()

	if got.ID != id {
		t.Fatalf("ID mismatch: got %v want %v", got.ID, id)
	}
	if got.Title != "" {
		t.Fatalf("Title mismatch: got %q want empty", got.Title)
	}
	if got.SensitiveInfo != "Data: " {
		t.Fatalf("SensitiveInfo mismatch: got %q want %q", got.SensitiveInfo, "Data: ")
	}
}

func TestUserTextData_ToRecord_WithUnicode(t *testing.T) {
	utd := UserTextData{
		ID:   uuid.New(),
		Data: "🚀 rocket launch",
		Note: "unicode test",
	}

	got := utd.ToRecord()

	if got.Title != "unicode test" {
		t.Fatalf("Title mismatch: got %q want %q", got.Title, "unicode test")
	}

	wantSensitive := "Data: 🚀 rocket launch"
	if got.SensitiveInfo != wantSensitive {
		t.Fatalf("SensitiveInfo mismatch:\n got: %q\n want: %q", got.SensitiveInfo, wantSensitive)
	}
}
