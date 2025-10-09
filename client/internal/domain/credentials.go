package domain

import (
	"fmt"

	"github.com/google/uuid"
)

// Credentials represents a pair of login and password
// with an optional note.
type Credentials struct {
	ID       uuid.UUID
	Login    string
	Password string
	Note     string
}

// ToRecord converts Credentials into a generic Record representation.
func (c Credentials) ToRecord() Record {
	return Record{
		ID:            c.ID,
		Title:         c.Note,
		Kind:          UserDataKindCredentials,
		SensitiveInfo: fmt.Sprintf("Login: %s | Password: %s", c.Login, c.Password),
	}
}
