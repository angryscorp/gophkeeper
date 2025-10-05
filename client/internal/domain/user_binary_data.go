package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type UserBinaryData struct {
	ID   uuid.UUID
	Data []byte
	Note string
}

func (b UserBinaryData) ToRecord() Record {
	return Record{
		ID:            b.ID,
		Title:         b.Note,
		Kind:          UserDataKindBinaryData,
		SensitiveInfo: fmt.Sprintf("Data: %s\n", b.Data),
	}
}
