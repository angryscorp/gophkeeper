package auth

import (
	"context"
	"gophkeeper/pkg/crypto"
	"time"
)

type GRPCClient interface {
	Register(ctx context.Context, username string, kdf crypto.KDFParameters, edKey, authKey []byte, algorithm crypto.AuthKeyAlgorithm) error
}

type Auth struct {
	client GRPCClient
}

func NewAuth(client GRPCClient) *Auth {
	return &Auth{
		client: client,
	}
}

func (a *Auth) Register(username string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	kdf := crypto.KDFParameters{
		Algorithm:   crypto.KDFAlgorithmARGON2ID,
		TimeCost:    3,
		MemoryCost:  64 * 1024,
		Parallelism: 1,
		Salt:        crypto.RandBytes(16),
	}

	return a.client.Register(
		ctx,
		username,
		kdf,
		[]byte("ed-key-stub"),
		[]byte("auth-stub"),
		crypto.AuthKeyAlgorithmHMACSHA256,
	)
}
