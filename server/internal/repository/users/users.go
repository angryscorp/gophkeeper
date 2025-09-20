package users

import (
	"context"
	"gophkeeper/server/internal/domain"
)

type Users interface {
	Get(context.Context) ([]domain.User, error)
	Add(context.Context, domain.User) error
}
