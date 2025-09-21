package users

import (
	"context"
	"gophkeeper/server/internal/domain"
)

const (
	UsernameIsAlreadyTaken = "username is already taken"
)

type Users interface {
	Get(context.Context) ([]domain.User, error)
	Add(context.Context, domain.User) error
}
