package domain

import "github.com/google/uuid"

type OutboxMessage struct {
	ID            uuid.UUID
	RecordID      uuid.UUID
	Kind          int32
	UpdatedAtUnix int64
	Payload       []byte
}
