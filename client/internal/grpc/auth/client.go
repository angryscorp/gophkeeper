package auth

import (
	"context"
	"gophkeeper/client/internal/domain"
	"gophkeeper/pkg/crypto"
	"gophkeeper/pkg/grpc/auth"
	"gophkeeper/pkg/grpc/mapper"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		Kdf:              mapper.KdfParametersToGRPC(kdf),
		EncryptedDataKey: edKey,
		AuthKey:          authKey,
		AuthKeyAlg:       mapper.AuthAlgoToGRPC(algorithm),
	}

	_, err := c.client.Register(ctx, req)
	if err != nil {
		return mapError(err)
	}

	return nil
}

func (c Client) LoginStart(ctx context.Context, username string, deviceName string) (*crypto.LoginPayload, error) {
	resp, err := c.client.LoginStart(ctx, &auth.LoginStartRequest{Username: username, DeviceName: deviceName})
	if err != nil {
		return nil, mapError(err)
	}
	return &crypto.LoginPayload{
		DeviceId:         resp.DeviceId,
		KDFParameters:    mapper.KdfParametersToDomain(resp.Kdf),
		EncryptedDataKey: resp.EncryptedDataKey,
		AuthKeyAlgorithm: mapper.AuthAlgoToDomain(resp.AuthKeyAlg),
		Challenge:        resp.Challenge,
	}, nil
}

func (c Client) LoginFinish(ctx context.Context, username, deviceName string, challenge []byte) (string, error) {
	resp, err := c.client.LoginFinish(ctx, &auth.LoginFinishRequest{
		Username: username,
		DeviceId: deviceName,
		Response: challenge,
	})
	if err != nil {
		return "", mapError(err)
	}
	return resp.AccessToken, nil
}

func mapError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.AlreadyExists:
		return domain.ErrUsernameTaken
	case codes.NotFound:
		return domain.ErrUsernameNotFound
	}

	return err
}
