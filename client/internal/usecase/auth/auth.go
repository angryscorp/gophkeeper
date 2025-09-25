package auth

import (
	"context"
	"fmt"
	"gophkeeper/client/internal/repository/tokens"
	"gophkeeper/pkg/crypto"
	"gophkeeper/pkg/device"
	"time"
)

const (
	authContext = "auth"
	ctxTimeout  = 5 * time.Second
)

type Client interface {
	Register(ctx context.Context, username string, kdf crypto.KDFParameters, edKey, authKey []byte, algorithm crypto.AuthKeyAlgorithm) error
	LoginStart(ctx context.Context, username string, deviceName string) (crypto.LoginPayload, error)
	LoginFinish(ctx context.Context, username, deviceName string, challenge []byte) (string, error)
}

type Auth struct {
	client Client
	repo   tokens.Tokens
}

func New(
	client Client,
	repo tokens.Tokens,
) *Auth {
	return &Auth{
		client: client,
		repo:   repo,
	}
}

func (auth *Auth) Register(username, password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	kdfParams := crypto.DefaultKDFParameters()

	// Calculate a master key
	masterKey, err := crypto.DeriveKey(password, kdfParams)
	if err != nil {
		return fmt.Errorf("failed to derive key: %w", err)
	}

	// Generate a data key
	dataKey := crypto.RandBytes(32)

	// Encrypt the data key with the master key
	encryptedDataKey, err := crypto.Encrypt(masterKey, dataKey)
	if err != nil {
		return err
	}

	// Generate an auth key from the data key
	authKey, err := crypto.DeriveAuthKey(dataKey, []byte(authContext))
	if err != nil {
		return err
	}

	return auth.client.Register(
		ctx,
		username,
		kdfParams,
		encryptedDataKey,
		authKey,
		crypto.DefaultAuthKeyAlgorithm(),
	)

}

func (auth *Auth) Login(username, password string) error {
	deviceName := device.GenerateDeviceName()
	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	// Start login
	resp, err := auth.client.LoginStart(ctx, username, deviceName)
	if err != nil {
		return err
	}

	// Calculate a master key
	masterKey, err := crypto.DeriveKey(password, resp.KDFParameters)
	if err != nil {
		return fmt.Errorf("failed to derive key: %w", err)
	}

	// Decrypt the data key
	dataKey, err := crypto.Decrypt(masterKey, resp.EncryptedDataKey)
	if err != nil {
		return err
	}

	// Derive the auth key from the data key
	authKey, err := crypto.DeriveAuthKey(dataKey, []byte(authContext))
	if err != nil {
		return err
	}

	// Create HMAC response for the challenge
	challengeResponse := crypto.SignChallenge(authKey, resp.Challenge, resp.AuthKeyAlgorithm)

	// Finish login
	token, err := auth.client.LoginFinish(ctx, username, deviceName, challengeResponse)
	if err != nil {
		return err
	}

	// Save an access token
	err = auth.repo.SaveAccessToken(ctx, token)
	if err != nil {
		return err
	}

	return nil
}
