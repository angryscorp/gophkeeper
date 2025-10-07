package domain

import (
	"gophkeeper/pkg/crypto"

	"github.com/google/uuid"
)

// User represents an application account stored in the domain layer.
// It includes key derivation parameters and cryptographic material
// required for authentication and data encryption.
type User struct {
	ID               uuid.UUID
	Username         string
	KDFParameters    crypto.KDFParameters
	EncryptedDataKey []byte
	AuthKeyAlgorithm crypto.AuthKeyAlgorithm
	AuthKey          []byte
}
