package crypto

// ProxyCryptoService is a lightweight wrapper around encryption and
// decryption functions bound to a specific data key.
// It delegates cryptographic operations to provided encryptor/decryptor
// functions, making it easy to swap implementations (e.g. AES-GCM, ChaCha20).
type ProxyCryptoService struct {
	encryptor func(key, plaintext []byte) ([]byte, error)
	decryptor func(key, ciphertext []byte) ([]byte, error)
	dataKey   []byte
}

// New creates a new ProxyCryptoService with given encryptor and decryptor
// functions. Both functions must accept a key and input data, and return
// the processed data or an error.
func New(
	encryptor func(key, plaintext []byte) ([]byte, error),
	decryptor func(key, ciphertext []byte) ([]byte, error),
) *ProxyCryptoService {
	return &ProxyCryptoService{
		encryptor: encryptor,
		decryptor: decryptor,
	}
}

// SetDataKey sets the data key used by the service for encryption
// and decryption. Must be called before Encrypt/Decrypt.
func (s *ProxyCryptoService) SetDataKey(dataKey []byte) {
	s.dataKey = dataKey
}

// Encrypt encrypts the given plaintext using the configured data key
// and the encryptor function provided at construction time.
func (s *ProxyCryptoService) Encrypt(plainText []byte) ([]byte, error) {
	return s.encryptor(s.dataKey, plainText)
}

// Decrypt decrypts the given ciphertext using the configured data key
// and the decryptor function provided at construction time.
func (s *ProxyCryptoService) Decrypt(cipherText []byte) ([]byte, error) {
	return s.decryptor(s.dataKey, cipherText)
}
