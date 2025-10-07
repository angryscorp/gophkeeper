package auth

import (
	"context"
	"fmt"
	"time"

	"gophkeeper/pkg/crypto"
	"gophkeeper/pkg/device"
)

const (
	authContext = "auth"
	ctxTimeout  = 5 * time.Second
)

// Auth coordinates client/server authentication and secure
// storage of keys and tokens. It handles registration and login.
type Auth struct {
	client        Client
	repo          Tokens
	dataKeySetter func(dataKey []byte)
}

// New creates a new Auth service with a client, token repo,
// and a setter function for the derived data key.
func New(
	client Client,
	repo Tokens,
	dataKeySetter func(dataKey []byte),
) *Auth {
	return &Auth{
		client:        client,
		repo:          repo,
		dataKeySetter: dataKeySetter,
	}
}

// Register creates a new user account. It derives a master key
// from the password, generates and encrypts a data key, derives
// an auth key, and sends the parameters to the server.
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

// Login authenticates an existing user. It performs the
// challenge/response handshake with the server, decrypts
// the stored data key, unlocks the local repo, and saves
// the access token.
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

	// Password is correct, unlock DB
	err = auth.repo.Unlock(dataKey)
	if err != nil {
		return err
	}

	// Save an access token
	err = auth.repo.SaveAccessToken(ctx, token)
	if err != nil {
		return err
	}

	auth.dataKeySetter(dataKey)

	return nil
}
