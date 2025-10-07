package domain

import "github.com/google/uuid"

// Record describes a stored item in the vault.
// It can be a credential, note, bank card, or file.
type Record struct {
	ID            uuid.UUID
	Title         string
	Kind          UserDataKind
	SensitiveInfo string
}
