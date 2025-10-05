package domain

import "github.com/google/uuid"

type Record struct {
	ID            uuid.UUID
	Title         string
	Kind          UserDataKind
	SensitiveInfo string
}
