package crypto

type LoginPayload struct {
	DeviceId         string
	KDFParameters    KDFParameters
	EncryptedDataKey []byte
	AuthKeyAlgorithm AuthKeyAlgorithm
	Challenge        []byte
}
