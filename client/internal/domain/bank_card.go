package domain

import (
	"fmt"

	"github.com/google/uuid"
)

// BankCard represents sensitive payment card information
// stored securely in the vault.
type BankCard struct {
	ID         uuid.UUID
	Owner      string
	Number     string
	CVV        string
	ExpireDate string
	Note       string
}

// ToRecord converts a BankCard into a generic Record representation.
func (b BankCard) ToRecord() Record {
	return Record{
		ID:            b.ID,
		Title:         b.Note,
		Kind:          UserDataKindBankCard,
		SensitiveInfo: fmt.Sprintf("Owner: %s | Number: %s | CVV: %s | ExpireDate: %s", b.Owner, b.Number, b.CVV, b.ExpireDate),
	}
}
