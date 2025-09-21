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

func mapKDFToDomain(kdf *auth.KdfParams) crypto.KDFParameters {
	var algorithm crypto.KDFAlgorithm
	switch kdf.Alg {
	case auth.KdfAlg_ARGON2ID:
		algorithm = crypto.KDFAlgorithmARGON2ID
	default:
		panic("Unknown KDF algorithm")
	}
	return crypto.KDFParameters{
		Algorithm:   algorithm,
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
	case crypto.AuthKeyAlgorithmHMACSHA512:
		authAlgorithm = auth.AuthKeyAlg_HMAC_SHA512
	}
	return authAlgorithm
}
func mapAuthAlgoToDomain(algorithm auth.AuthKeyAlg) crypto.AuthKeyAlgorithm {
	var authAlgorithm crypto.AuthKeyAlgorithm
	switch algorithm {
	case auth.AuthKeyAlg_HMAC_SHA256:
		authAlgorithm = crypto.AuthKeyAlgorithmHMACSHA256
	case auth.AuthKeyAlg_HMAC_SHA512:
		authAlgorithm = crypto.AuthKeyAlgorithmHMACSHA512
	}
	return authAlgorithm
}
