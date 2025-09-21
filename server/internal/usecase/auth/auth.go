package auth

import (
	"context"
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
