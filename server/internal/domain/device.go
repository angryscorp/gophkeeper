package domain

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID        uuid.UUID
	Name      string
	CreatedAt time.Time
	LastSeen  time.Time
}
