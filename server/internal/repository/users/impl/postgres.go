package impl

import (
	"context"
	"fmt"
	"gophkeeper/pkg/crypto"
	"gophkeeper/pkg/pgx"
	"gophkeeper/server/internal/domain"
	"gophkeeper/server/internal/repository/users"
	"gophkeeper/server/internal/repository/users/db"

	"github.com/jackc/pgx/v5/pgxpool"
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

func (repo Users) Get(ctx context.Context, username string) (domain.User, error) {
	row, err := repo.queries.Get(ctx, username)
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		ID:       row.ID,
		Username: row.Username,
		KDFParameters: crypto.KDFParameters{
			Algorithm:   crypto.KDFAlgorithm(row.KdfAlgorithm),
			TimeCost:    uint32(row.KdfTimeCost),
			MemoryCost:  uint32(row.KdfMemoryCost),
			Parallelism: uint32(row.KdfParallelism),
			Salt:        row.KdfSalt,
		},
		EncryptedDataKey: row.EncryptedDataKey,
		AuthKeyAlgorithm: crypto.AuthKeyAlgorithm(row.AuthKeyAlgorithm),
		AuthKey:          row.AuthKey,
	}, nil
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
