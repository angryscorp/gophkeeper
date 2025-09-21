package domain

import "gophkeeper/pkg/crypto"

type LoginPayload struct {
	DeviceId         string
	KDFParameters    crypto.KDFParameters
	EncryptedDataKey []byte
	AuthKeyAlgorithm crypto.AuthKeyAlgorithm
	Challenge        []byte
}
