package auth

import (
	"context"
	"crypto/rand"
	"gophkeeper/pkg/grpc/auth"

	"google.golang.org/grpc"
)

type Client struct {
	client auth.AuthServiceClient
}

func New(conn *grpc.ClientConn) *Client {
	return &Client{client: auth.NewAuthServiceClient(conn)}
}

func (c Client) Register(ctx context.Context, username string) (*auth.RegisterResponse, error) {
	kdf := &auth.KdfParams{
		Alg:         auth.KdfAlg_ARGON2ID,
		TimeCost:    3,
		MemoryCost:  64 * 1024,
		Parallelism: 1,
		Salt:        randBytes(16),
	}

	regReq := &auth.RegisterRequest{
		Username:         username,
		Kdf:              kdf,
		EncryptedDataKey: []byte("EncryptedDataKey"),
		AuthKey:          []byte("AuthKey"),
		AuthKeyAlg:       auth.AuthKeyAlg_HMAC_SHA256,
	}

	return c.client.Register(ctx, regReq)
}

func randBytes(n int) []byte {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return b
}
