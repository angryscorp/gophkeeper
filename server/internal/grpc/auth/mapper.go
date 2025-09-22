package auth

import (
	"gophkeeper/pkg/grpc/auth"
	"gophkeeper/pkg/grpc/mapper"
	"gophkeeper/server/internal/domain"

	"github.com/google/uuid"
)

func requestToDomain(req *auth.RegisterRequest) domain.User {
	return domain.User{
		ID:               uuid.New(),
		Username:         req.Username,
		KDFParameters:    mapper.KdfParametersToDomain(req.Kdf),
		EncryptedDataKey: req.EncryptedDataKey,
		AuthKeyAlgorithm: mapper.AuthAlgoToDomain(req.AuthKeyAlg),
		AuthKey:          req.AuthKey,
	}
}
