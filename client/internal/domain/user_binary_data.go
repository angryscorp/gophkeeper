package domain

import (
	"fmt"

	"github.com/google/uuid"
)

// UserBinaryData represents arbitrary binary data stored by the user,
// such as files, images, or other non-text content, with an optional note.
type UserBinaryData struct {
	ID   uuid.UUID
	Data []byte
	Note string
}

// ToRecord converts UserBinaryData into a generic Record representation.
func (b UserBinaryData) ToRecord() Record {
	return Record{
		ID:            b.ID,
		Title:         b.Note,
		Kind:          UserDataKindBinaryData,
		SensitiveInfo: fmt.Sprintf("Data: %s", b.Data),
	}
}
