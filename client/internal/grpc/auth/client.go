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

// Client is a thin wrapper around the generated gRPC AuthServiceClient.
// It provides convenient methods for user registration and login flow,
// mapping gRPC types into domain/crypto types and domain errors.
type Client struct {
	client auth.AuthServiceClient
}

// New creates a new Auth client bound to the given gRPC connection.
func New(conn *grpc.ClientConn) *Client {
	return &Client{client: auth.NewAuthServiceClient(conn)}
}

// Register calls the AuthService.Register RPC to create a new user account.
// It sends username, KDF parameters, encrypted data key and auth key.
// Returns ErrUsernameTaken if the username already exists.
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

// LoginStart calls the AuthService.LoginStart RPC to begin login.
// Returns the server-provided KDF parameters, encrypted data key,
// auth key algorithm and a random challenge that must be signed.
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

// LoginFinish calls the AuthService.LoginFinish RPC to complete login.
// The client sends the HMAC-signed challenge response. On success,
// the server issues an access token.
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

// mapError translates gRPC status codes from the Auth service
// into domain-level errors for easier handling by the client.
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
