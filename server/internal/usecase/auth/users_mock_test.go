package auth

import (
	"context"
	"errors"

	"gophkeeper/server/internal/domain"
)

type stubUsers struct {
	getFn func(ctx context.Context, username string) (*domain.User, error)
	addFn func(ctx context.Context, user domain.User) error
}

func (s stubUsers) Get(ctx context.Context, username string) (*domain.User, error) {
	if s.getFn != nil {
		return s.getFn(ctx, username)
	}
	return nil, errors.New("not implemented")
}

func (s stubUsers) Add(ctx context.Context, user domain.User) error {
	if s.addFn != nil {
		return s.addFn(ctx, user)
	}
	return nil
}
