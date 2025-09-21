package domain

import "github.com/google/uuid"

type Credentials struct {
	ID       uuid.UUID
	Login    string
	Password string
	Note     string
}
