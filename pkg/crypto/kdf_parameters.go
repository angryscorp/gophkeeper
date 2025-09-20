package crypto

type KDFParameters struct {
	Algorithm   KDFAlgorithm
	TimeCost    uint32
	MemoryCost  uint32
	Parallelism uint32
	Salt        []byte
}

type KDFAlgorithm string

const KDFAlgorithmARGON2ID KDFAlgorithm = "ARGON2ID"
