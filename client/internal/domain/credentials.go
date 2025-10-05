package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type Credentials struct {
	ID       uuid.UUID
	Login    string
	Password string
	Note     string
}

func (c Credentials) ToRecord() Record {
	return Record{
		ID:            c.ID,
		Title:         c.Note,
		Kind:          UserDataKindCredentials,
		SensitiveInfo: fmt.Sprintf("Login: %s | Password: %s", c.Login, c.Password),
	}
}
