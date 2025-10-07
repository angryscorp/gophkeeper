package crypto

// LoginPayload contains all parameters the server sends to the client
// to complete the login flow.
type LoginPayload struct {
	DeviceId         string
	KDFParameters    KDFParameters
	EncryptedDataKey []byte
	AuthKeyAlgorithm AuthKeyAlgorithm
	Challenge        []byte
}
