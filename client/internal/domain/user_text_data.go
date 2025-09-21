package domain

import "github.com/google/uuid"

type UserTextData struct {
	ID   uuid.UUID
	Data string
	Note string
}
