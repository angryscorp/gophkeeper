package domain

import (
	"gophkeeper/pkg/crypto"

	"github.com/google/uuid"
)

type User struct {
	ID               uuid.UUID
	Username         string
	KDFParameters    crypto.KDFParameters
	EncryptedDataKey []byte
	AuthKeyAlgorithm crypto.AuthKeyAlgorithm
	AuthKey          []byte
}
