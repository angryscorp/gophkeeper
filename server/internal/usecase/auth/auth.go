package auth

import (
	"context"
	"crypto/hmac"
	"errors"
	"gophkeeper/pkg/crypto"
	"gophkeeper/server/internal/domain"
	"gophkeeper/server/internal/repository/users"
	"log"
)

type Auth struct {
	repo users.Users
}

func New(repo users.Users) *Auth {
	return &Auth{
		repo: repo,
	}
}

func (auth *Auth) Register(ctx context.Context, user domain.User) error {
	log.Printf("Registering user: %s\n", user.Username)
	return auth.repo.Add(ctx, user)
}

func (auth *Auth) LoginStart(ctx context.Context, username, deviceId string) (crypto.LoginPayload, error) {
	log.Printf("Starting login for user: %s\n", username)
	resp, err := auth.repo.Get(ctx, username)
	if err != nil {
		return crypto.LoginPayload{}, err
	}

	return crypto.LoginPayload{
		DeviceId:         deviceId,
		KDFParameters:    resp.KDFParameters,
		EncryptedDataKey: resp.EncryptedDataKey,
		AuthKeyAlgorithm: resp.AuthKeyAlgorithm,
		Challenge:        []byte("12345"),
	}, nil
}

func (auth *Auth) LoginFinish(ctx context.Context, deviceName string, challenge []byte) error {
	log.Printf("Finishing login: %s\n", deviceName)
	resp, err := auth.repo.Get(ctx, "test_user")
	if err != nil {
		return err
	}
	expected := crypto.SignChallenge(resp.AuthKey, []byte("12345"), resp.AuthKeyAlgorithm)
	if !hmac.Equal(expected, challenge) {
		return errors.New("invalid challenge")
	}
	return nil
}
