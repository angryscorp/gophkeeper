package crypto

type ProxyCryptoService struct {
	encryptor func(key, plaintext []byte) ([]byte, error)
	dataKey   []byte
}

func New(encryptor func(key, plaintext []byte) ([]byte, error)) *ProxyCryptoService {
	return &ProxyCryptoService{
		encryptor: encryptor,
	}
}

func (s *ProxyCryptoService) SetDataKey(dataKey []byte) {
	s.dataKey = dataKey
}

func (s *ProxyCryptoService) Encrypt(plainText []byte) ([]byte, error) {
	return s.encryptor(s.dataKey, plainText)
}
