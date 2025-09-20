package crypto

type AuthKeyAlgorithm string

const (
	AuthKeyAlgorithmHMACSHA256 AuthKeyAlgorithm = "HMAC_SHA256"
	AuthKeyAlgorithmHMACSHA512 AuthKeyAlgorithm = "HMAC_SHA512"
)
