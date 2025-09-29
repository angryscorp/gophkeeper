package mapper

import (
	"gophkeeper/pkg/crypto"
	"gophkeeper/pkg/grpc/auth"
)

func KdfParametersToDomain(kdf *auth.KdfParams) crypto.KDFParameters {
	return crypto.KDFParameters{
		Algorithm:   kdfAlgoToDomain(kdf.Alg),
		TimeCost:    kdf.TimeCost,
		MemoryCost:  kdf.MemoryCost,
		Parallelism: kdf.Parallelism,
		Salt:        kdf.Salt,
	}
}

func KdfParametersToGRPC(kdf crypto.KDFParameters) *auth.KdfParams {
	return &auth.KdfParams{
		Alg:         kdfAlgoToGRPC(kdf.Algorithm),
		TimeCost:    kdf.TimeCost,
		MemoryCost:  kdf.MemoryCost,
		Parallelism: kdf.Parallelism,
		Salt:        kdf.Salt,
	}
}

func kdfAlgoToDomain(algorithm auth.KdfAlg) crypto.KDFAlgorithm {
	switch algorithm {
	case auth.KdfAlg_ARGON2ID:
		return crypto.KDFAlgorithmARGON2ID
	default:
		panic("Unknown KDF algorithm")
	}
}

func kdfAlgoToGRPC(algorithm crypto.KDFAlgorithm) auth.KdfAlg {
	switch algorithm {
	case crypto.KDFAlgorithmARGON2ID:
		return auth.KdfAlg_ARGON2ID
	default:
		panic("Unknown KDF algorithm")
	}
}

func AuthAlgoToDomain(algorithm auth.AuthKeyAlg) crypto.AuthKeyAlgorithm {
	switch algorithm {
	case auth.AuthKeyAlg_HMAC_SHA256:
		return crypto.AuthKeyAlgorithmHMACSHA256
	case auth.AuthKeyAlg_HMAC_SHA512:
		return crypto.AuthKeyAlgorithmHMACSHA512
	default:
		panic("Unknown AuthKey algorithm")
	}
}

func AuthAlgoToGRPC(algorithm crypto.AuthKeyAlgorithm) auth.AuthKeyAlg {
	switch algorithm {
	case crypto.AuthKeyAlgorithmHMACSHA256:
		return auth.AuthKeyAlg_HMAC_SHA256
	case crypto.AuthKeyAlgorithmHMACSHA512:
		return auth.AuthKeyAlg_HMAC_SHA512
	default:
		panic("Unknown AuthKey algorithm")
	}
}
