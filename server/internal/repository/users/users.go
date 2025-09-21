package users

import (
	"context"
	"gophkeeper/server/internal/domain"
)

const (
	UsernameIsAlreadyTaken = "username is already taken"
)

type Users interface {
	Get(ctx context.Context, username string) (domain.User, error)
	Add(ctx context.Context, user domain.User) error
}
