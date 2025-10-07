package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gophkeeper/pkg/crypto"
	"gophkeeper/pkg/pgx"
	"gophkeeper/server/internal/domain"
	"gophkeeper/server/internal/repository/users/db"
	"gophkeeper/server/internal/usecase/auth"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Users is a repository that loads and stores user records in PostgreSQL.
type Users struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

// New creates a Users repository backed by a pgx pool created from DSN.
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

var _ auth.Users = (*Users)(nil)

// Get fetches a user by username. It returns domain.ErrUsernameNotFound
// if no user exists with the given username.
func (repo Users) Get(ctx context.Context, username string) (*domain.User, error) {
	row, err := repo.queries.Get(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUsernameNotFound
		}
		return nil, err
	}

	return &domain.User{
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

// Add inserts a new user. If the username already exists, it returns
// domain.ErrUsernameTaken (mapped from unique constraint violation).
func (repo Users) Add(ctx context.Context, user domain.User) error {
	err := repo.queries.Add(ctx, db.AddParams{
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

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// duplicate key value violates unique constraint
			if pgErr.Code == "23505" && pgErr.ConstraintName == "users_username_key" {
				return domain.ErrUsernameTaken
			}
		}
		return err
	}

	return nil
}
