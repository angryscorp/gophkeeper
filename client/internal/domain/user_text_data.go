package domain

import (
	"fmt"

	"github.com/google/uuid"
)

// UserTextData represents a piece of plain text stored by the user,
// with an optional note for context or description.
type UserTextData struct {
	ID   uuid.UUID
	Data string
	Note string
}

// ToRecord converts UserTextData into a generic Record representation.
func (u UserTextData) ToRecord() Record {
	return Record{
		ID:            u.ID,
		Title:         u.Note,
		Kind:          UserDataKindTextData,
		SensitiveInfo: fmt.Sprintf("Data: %s", u.Data),
	}
}
