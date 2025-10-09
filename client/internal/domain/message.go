package domain

import "github.com/google/uuid"

// Message represents a unit of data exchange during
// synchronization with the server.
type Message struct {
	ID            uuid.UUID
	RecordID      uuid.UUID
	Kind          int32
	UpdatedAtUnix int64
	Payload       []byte
}
