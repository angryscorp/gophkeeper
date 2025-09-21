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

func (repo Users) Get(ctx context.Context) ([]domain.User, error) {
	usernames, err := repo.queries.Get(ctx)
	if err != nil {
		return nil, err
	}

	return lo.Map(usernames, func(item db.GetRow, index int) domain.User {
		return domain.User{ID: item.ID, Username: item.Username}
	}), nil
}

func (repo Users) Add(ctx context.Context, user domain.User) error {
	return repo.queries.Add(ctx, db.AddParams{
		ID:               user.ID,
		Username:         user.Username,
		KdfAlgorithm:     string(user.KDFParameters.Algorithm),
		KdfTimeCost:      int32(user.KDFParameters.TimeCost),
		KdfMemoryCost:    int32(user.KDFParameters.MemoryCost),
		KdfParallelism:   int32(user.KDFParameters.Parallelism),
		KdfSalt:          user.KDFParameters.Salt,
		EncryptedDataKey: user.EncryptedDataKey,
		AuthKey:          user.AuthKey,
		AuthKeyAlgorithm: string(user.AuthKeyAlgorithm),
	})
}
