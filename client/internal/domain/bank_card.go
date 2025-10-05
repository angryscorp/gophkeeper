package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type BankCard struct {
	ID         uuid.UUID
	Owner      string
	Number     string
	CVV        string
	ExpireDate string
	Note       string
}

func (b BankCard) ToRecord() Record {
	return Record{
		ID:            b.ID,
		Title:         b.Note,
		Kind:          UserDataKindBankCard,
		SensitiveInfo: fmt.Sprintf("Owner: %s | Number: %s | CVV: %s | ExpireDate: %s", b.Owner, b.Number, b.CVV, b.ExpireDate),
	}
}
