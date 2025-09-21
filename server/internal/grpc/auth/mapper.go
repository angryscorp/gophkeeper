package auth

import (
	"gophkeeper/pkg/crypto"
	"gophkeeper/pkg/grpc/auth"
	"gophkeeper/server/internal/domain"

	"github.com/google/uuid"
)

func requestToDomain(req *auth.RegisterRequest) domain.User {
	return domain.User{
		ID:               uuid.New(),
		Username:         req.Username,
		KDFParameters:    kdfParametersToDomain(req.Kdf),
		EncryptedDataKey: req.EncryptedDataKey,
		AuthKeyAlgorithm: authAlgoToDomain(req.AuthKeyAlg),
		AuthKey:          req.AuthKey,
	}
}

func kdfParametersToDomain(kdf *auth.KdfParams) crypto.KDFParameters {
	return crypto.KDFParameters{
		Algorithm:   kdfAlgoToDomain(kdf.Alg),
		TimeCost:    kdf.TimeCost,
		MemoryCost:  kdf.MemoryCost,
		Parallelism: kdf.Parallelism,
		Salt:        kdf.Salt,
	}
}

func kdfAlgoToDomain(algo auth.KdfAlg) crypto.KDFAlgorithm {
	switch algo {
	case auth.KdfAlg_ARGON2ID:
		return crypto.KDFAlgorithmARGON2ID
	default:
		panic("Unknown KDF algorithm")
	}
}

func authAlgoToDomain(algo auth.AuthKeyAlg) crypto.AuthKeyAlgorithm {
	switch algo {
	case auth.AuthKeyAlg_HMAC_SHA256:
		return crypto.AuthKeyAlgorithmHMACSHA256
	case auth.AuthKeyAlg_HMAC_SHA512:
		return crypto.AuthKeyAlgorithmHMACSHA512
	default:
		panic("Unknown KDF algorithm")
	}
}
