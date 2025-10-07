package auth

import (
	"context"
	"gophkeeper/server/internal/domain"
)

// Users defines access to user persistence.
type Users interface {
	Get(ctx context.Context, username string) (*domain.User, error)
	Add(ctx context.Context, user domain.User) error
}
