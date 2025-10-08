package domain

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestBankCard_ToRecord_Basic(t *testing.T) {
	id := uuid.New()
	bc := BankCard{
		ID:         id,
		Owner:      "John Doe",
		Number:     "4111 1111 1111 1111",
		CVV:        "123",
		ExpireDate: "12/34",
		Note:       "My primary card",
	}

	got := bc.ToRecord()

	if got.ID != id {
		t.Fatalf("ID mismatch: got %v want %v", got.ID, id)
	}
	if got.Title != bc.Note {
		t.Fatalf("Title mismatch: got %q want %q", got.Title, bc.Note)
	}
	if got.Kind != UserDataKindBankCard {
		t.Fatalf("Kind mismatch: got %v want %v", got.Kind, UserDataKindBankCard)
	}

	wantSensitive := fmt.Sprintf(
		"Owner: %s | Number: %s | CVV: %s | ExpireDate: %s",
		bc.Owner, bc.Number, bc.CVV, bc.ExpireDate,
	)
	if got.SensitiveInfo != wantSensitive {
		t.Fatalf("SensitiveInfo mismatch:\n got:  %q\n want: %q", got.SensitiveInfo, wantSensitive)
	}
}

func TestBankCard_ToRecord_EmptyFields(t *testing.T) {
	id := uuid.New()
	bc := BankCard{
		ID:         id,
		Owner:      "",
		Number:     "",
		CVV:        "",
		ExpireDate: "",
		Note:       "",
	}

	got := bc.ToRecord()

	if got.ID != id {
		t.Fatalf("ID mismatch: got %v want %v", got.ID, id)
	}
	if got.Title != "" {
		t.Fatalf("Title mismatch: got %q want empty", got.Title)
	}

	wantSensitive := "Owner:  | Number:  | CVV:  | ExpireDate: "
	if got.SensitiveInfo != wantSensitive {
		t.Fatalf("SensitiveInfo mismatch:\n got:  %q\n want: %q", got.SensitiveInfo, wantSensitive)
	}
}

func TestBankCard_ToRecord_Unicode(t *testing.T) {
	bc := BankCard{
		ID:         uuid.New(),
		Owner:      "Jack Smith",
		Number:     "4000 0000 0000 0002",
		CVV:        "234",
		ExpireDate: "03/29",
		Note:       "💳 Alfa-Bank",
	}

	got := bc.ToRecord()

	if got.Title != "💳 Alfa-Bank" {
		t.Fatalf("unicode Title mismatch: got %q", got.Title)
	}

	wantSensitive := "Owner: Jack Smith | Number: 4000 0000 0000 0002 | CVV: 234 | ExpireDate: 03/29"
	if got.SensitiveInfo != wantSensitive {
		t.Fatalf("unicode SensitiveInfo mismatch:\n got:  %q\n want: %q", got.SensitiveInfo, wantSensitive)
	}
}
