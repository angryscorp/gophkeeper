package auth

import (
	"context"
	"gophkeeper/server/internal/domain"
	"gophkeeper/server/internal/repository/users"

	"github.com/google/uuid"
)

type Auth struct {
	repo users.Users
}

func New(repo users.Users) *Auth {
	return &Auth{
		repo: repo,
	}
}

func (a *Auth) Register(ctx context.Context, username string) error {
	user := domain.User{
		ID:       uuid.New(),
		Username: username,
	}
	return a.repo.Add(ctx, user)
}
