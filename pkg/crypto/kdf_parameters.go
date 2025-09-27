package crypto

// KDFParameters represents the parameters used to derive a key from a password.
type KDFParameters struct {
	Algorithm   KDFAlgorithm
	TimeCost    uint32
	MemoryCost  uint32
	Parallelism uint32
	Salt        []byte
}

// KDFAlgorithm represents the KDF algorithm used to derive a key from a password.
type KDFAlgorithm string

const KDFAlgorithmARGON2ID KDFAlgorithm = "ARGON2ID"

// DefaultKDFParameters returns the default KDF parameters.
func DefaultKDFParameters() KDFParameters {
	return KDFParameters{
		Algorithm:   KDFAlgorithmARGON2ID,
		TimeCost:    3,
		MemoryCost:  64 * 1024,
		Parallelism: 4,
		Salt:        RandBytes(32),
	}
}
