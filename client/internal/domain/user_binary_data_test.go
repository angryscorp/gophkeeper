package domain

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestUserBinaryData_ToRecord_Basic(t *testing.T) {
	id := uuid.New()
	data := []byte("hello world")
	bin := UserBinaryData{
		ID:   id,
		Data: data,
		Note: "test file",
	}

	got := bin.ToRecord()

	if got.ID != id {
		t.Fatalf("ID mismatch: got %v want %v", got.ID, id)
	}
	if got.Title != bin.Note {
		t.Fatalf("Title mismatch: got %q want %q", got.Title, bin.Note)
	}
	if got.Kind != UserDataKindBinaryData {
		t.Fatalf("Kind mismatch: got %v want %v", got.Kind, UserDataKindBinaryData)
	}

	wantSensitive := fmt.Sprintf("Data: %s", data)
	if got.SensitiveInfo != wantSensitive {
		t.Fatalf("SensitiveInfo mismatch:\n got:  %q\n want: %q", got.SensitiveInfo, wantSensitive)
	}
}

func TestUserBinaryData_ToRecord_EmptyFields(t *testing.T) {
	id := uuid.New()
	bin := UserBinaryData{
		ID:   id,
		Data: []byte{},
		Note: "",
	}

	got := bin.ToRecord()

	if got.ID != id {
		t.Fatalf("ID mismatch: got %v want %v", got.ID, id)
	}
	if got.Title != "" {
		t.Fatalf("Title mismatch: got %q want empty", got.Title)
	}

	wantSensitive := "Data: "
	if got.SensitiveInfo != wantSensitive {
		t.Fatalf("SensitiveInfo mismatch:\n got: %q\n want: %q", got.SensitiveInfo, wantSensitive)
	}
}

func TestUserBinaryData_ToRecord_WithBytes(t *testing.T) {
	id := uuid.New()
	bin := UserBinaryData{
		ID:   id,
		Data: []byte("hello"),
		Note: "binary sample",
	}

	got := bin.ToRecord()

	wantSensitive := "Data: hello"
	if got.SensitiveInfo != wantSensitive {
		t.Fatalf("SensitiveInfo mismatch:\n got: %q\n want: %q", got.SensitiveInfo, wantSensitive)
	}
}
