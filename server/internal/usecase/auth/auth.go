package auth

import (
	"context"
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

func (a *Auth) Register(ctx context.Context, user domain.User) error {
	log.Printf("Registering user: %s\n", user.Username)
	return a.repo.Add(ctx, user)
}

func (a *Auth) LoginStart(ctx context.Context, username string) (domain.LoginPayload, error) {
	log.Printf("Starting login for user: %s\n", username)
	return domain.LoginPayload{
		DeviceId:         "device-id",
		KDFParameters:    crypto.DefaultKDFParameters(),
		EncryptedDataKey: []byte("encrypted-data-key"),
		AuthKeyAlgorithm: crypto.DefaultAuthKeyAlgorithm(),
		Challenge:        []byte("challenge"),
	}, nil
}
