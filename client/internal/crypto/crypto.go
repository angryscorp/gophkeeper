package crypto

type ProxyCryptoService struct {
	encryptor func(key, plaintext []byte) ([]byte, error)
	decryptor func(key, ciphertext []byte) ([]byte, error)
	dataKey   []byte
}

func New(
	encryptor func(key, plaintext []byte) ([]byte, error),
	decryptor func(key, ciphertext []byte) ([]byte, error),
) *ProxyCryptoService {
	return &ProxyCryptoService{
		encryptor: encryptor,
		decryptor: decryptor,
	}
}

func (s *ProxyCryptoService) SetDataKey(dataKey []byte) {
	s.dataKey = dataKey
}

func (s *ProxyCryptoService) Encrypt(plainText []byte) ([]byte, error) {
	return s.encryptor(s.dataKey, plainText)
}

func (s *ProxyCryptoService) Decrypt(cipherText []byte) ([]byte, error) {
	return s.decryptor(s.dataKey, cipherText)
}
