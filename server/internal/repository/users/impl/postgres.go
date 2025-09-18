package impl

import (
	"context"
	"fmt"
	"gophkeeper/pkg/pgx"
	"gophkeeper/server/internal/domain"
	"gophkeeper/server/internal/repository/users"
	"gophkeeper/server/internal/repository/users/db"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/lo"
)

type Users struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func New(dsn string) (*Users, func(), error) {
	pool, err := pgx.CreatePGXPool(dsn)
	if err != nil {
		return nil, func() {}, fmt.Errorf("failed to create pool: %w", err)
	}

	return &Users{
		queries: db.New(pool),
		pool:    pool,
	}, pool.Close, nil
}

var _ users.Users = (*Users)(nil)

func (u Users) Get(ctx context.Context) ([]domain.User, error) {
	usernames, err := u.queries.Get(ctx)
	if err != nil {
		return nil, err
	}

	return lo.Map(usernames, func(item db.User, index int) domain.User {
		return domain.User{ID: item.ID, Username: item.Username}
	}), nil
}

func (u Users) Add(ctx context.Context, user domain.User) error {
	return u.queries.Add(ctx, db.AddParams{
		ID:       user.ID,
		Username: user.Username,
	})
}
