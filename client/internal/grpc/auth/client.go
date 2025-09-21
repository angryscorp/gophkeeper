package auth

import (
	"context"
	"gophkeeper/pkg/crypto"
	"gophkeeper/pkg/grpc/auth"

	"google.golang.org/grpc"
)

type Client struct {
	client auth.AuthServiceClient
}

func New(conn *grpc.ClientConn) *Client {
	return &Client{client: auth.NewAuthServiceClient(conn)}
}

func (c Client) Register(ctx context.Context, username string, kdf crypto.KDFParameters, edKey, authKey []byte, algorithm crypto.AuthKeyAlgorithm) error {
	req := &auth.RegisterRequest{
		Username:         username,
		Kdf:              mapKDFToGRPC(kdf),
		EncryptedDataKey: edKey,
		AuthKey:          authKey,
		AuthKeyAlg:       mapAuthAlgoToRPC(algorithm),
	}

	_, err := c.client.Register(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
