package domain

import "github.com/google/uuid"

// Message represents a change (event) that is synchronized between
// client and server. It is used in the sync journal / outbox.
type Message struct {
	ID            uuid.UUID
	RecordID      uuid.UUID
	Kind          int32
	UpdatedAtUnix int64
	Payload       []byte
}
