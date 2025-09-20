package crypto

import (
	"time"

	"github.com/google/uuid"
)

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
