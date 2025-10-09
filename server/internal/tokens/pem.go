package tokens

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"os"
)

// LoadPrivateKey reads an Ed25519 private key from a PEM file in PKCS8 format.
func LoadPrivateKey(filePath string) (ed25519.PrivateKey, error) {
	pemBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("bad private pem")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pk, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, errors.New("not ed25519 private")
	}
	return pk, nil
}

// LoadPublicKey reads an Ed25519 public key from a PEM file in PKIX format.
func LoadPublicKey(filePath string) (ed25519.PublicKey, error) {
	pemBytes, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("bad public pem")
	}
	pubAny, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub, ok := pubAny.(ed25519.PublicKey)
	if !ok {
		return nil, errors.New("not ed25519 public")
	}
	return pub, nil
}
