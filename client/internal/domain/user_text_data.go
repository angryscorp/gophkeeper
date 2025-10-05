package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type UserTextData struct {
	ID   uuid.UUID
	Data string
	Note string
}

func (u UserTextData) ToRecord() Record {
	return Record{
		ID:            u.ID,
		Title:         u.Note,
		Kind:          UserDataKindTextData,
		SensitiveInfo: fmt.Sprintf("Data: %s", u.Data),
	}
}
