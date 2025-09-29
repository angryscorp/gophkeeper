package crypto

// AuthKeyAlgorithm represents the type of algorithm used to derive an authentication key, such as HMAC-SHA256 or HMAC-SHA512.
type AuthKeyAlgorithm string

const (
	AuthKeyAlgorithmHMACSHA256 AuthKeyAlgorithm = "HMAC_SHA256"
	AuthKeyAlgorithmHMACSHA512 AuthKeyAlgorithm = "HMAC_SHA512"
)

// DefaultAuthKeyAlgorithm returns the default algorithm used to derive an authentication key.
func DefaultAuthKeyAlgorithm() AuthKeyAlgorithm {
	return AuthKeyAlgorithmHMACSHA256
}
