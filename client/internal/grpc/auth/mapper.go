package auth

import (
	"gophkeeper/pkg/crypto"
	"gophkeeper/pkg/grpc/auth"
)

func mapKDFToGRPC(kdf crypto.KDFParameters) *auth.KdfParams {
	algorithm := auth.KdfAlg_KDF_ALG_UNSPECIFIED
	switch kdf.Algorithm {
	case crypto.KDFAlgorithmARGON2ID:
		algorithm = auth.KdfAlg_ARGON2ID
	}
	return &auth.KdfParams{
		Alg:         algorithm,
		TimeCost:    kdf.TimeCost,
		MemoryCost:  kdf.MemoryCost,
		Parallelism: kdf.Parallelism,
		Salt:        kdf.Salt,
	}
}

func mapAuthAlgoToRPC(algorithm crypto.AuthKeyAlgorithm) auth.AuthKeyAlg {
	authAlgorithm := auth.AuthKeyAlg_AUTH_KEY_ALG_UNSPECIFIED
	switch algorithm {
	case crypto.AuthKeyAlgorithmHMACSHA256:
		authAlgorithm = auth.AuthKeyAlg_HMAC_SHA256
	}
	return authAlgorithm
}
