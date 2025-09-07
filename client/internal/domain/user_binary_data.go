package domain

import "github.com/google/uuid"

type UserBinaryData struct {
	ID   uuid.UUID
	Data []byte
	Note string
}
