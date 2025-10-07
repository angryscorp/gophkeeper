package crypto

import (
	"time"

	"github.com/google/uuid"
)

// Record represents an encrypted user record stored locally
// and synchronized with the server.
type Record struct {
	ID          uuid.UUID
	Version     int64
	IsDeleted   bool
	Nonce       []byte
	Ciphertext  []byte
	Tag         []byte
	UpdatedAt   time.Time
	OperationID uuid.UUID
}
