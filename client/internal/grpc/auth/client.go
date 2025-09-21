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

func (c Client) LoginStart(ctx context.Context, username string, deviceName string) (crypto.LoginPayload, error) {
	resp, err := c.client.LoginStart(ctx, &auth.LoginStartRequest{Username: username, DeviceName: deviceName})
	if err != nil {
		return crypto.LoginPayload{}, err
	}
	return crypto.LoginPayload{
		DeviceId:         resp.DeviceId,
		KDFParameters:    mapKDFToDomain(resp.Kdf),
		EncryptedDataKey: resp.EncryptedDataKey,
		AuthKeyAlgorithm: mapAuthAlgoToDomain(resp.AuthKeyAlg),
		Challenge:        resp.Challenge,
	}, nil
}

func (c Client) LoginFinish(ctx context.Context, deviceName string, challenge []byte) error {
	resp, err := c.client.LoginFinish(ctx, &auth.LoginFinishRequest{
		DeviceId: deviceName,
		Response: challenge,
	})
	if err != nil {
		return err
	}

	_ = resp

	return nil
}
