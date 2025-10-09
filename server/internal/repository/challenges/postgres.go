package challenges

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"gophkeeper/pkg/crypto"
	pkgpgx "gophkeeper/pkg/pgx"
	"gophkeeper/server/internal/domain"
	"gophkeeper/server/internal/repository/challenges/db"
	"gophkeeper/server/internal/usecase/auth"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Challenges provides access to the login challenge storage (PostgreSQL) and
// exposes methods to create and atomically validate challenges.
type Challenges struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

// New initializes a Challenges repository backed by a pgx connection pool.
func New(dsn string) (*Challenges, func(), error) {
	pool, err := pkgpgx.CreatePGXPool(dsn)
	if err != nil {
		return nil, func() {}, fmt.Errorf("failed to create pool: %w", err)
	}

	return &Challenges{
		queries: db.New(pool),
		pool:    pool,
	}, pool.Close, nil
}

var _ auth.Challenges = (*Challenges)(nil)

// Add persists a new challenge for a given user/device with an expiration time.
// The challenge bytes are stored as-is; no validation is performed here.
func (c Challenges) Add(ctx context.Context, userId uuid.UUID, deviceName string, challenge []byte, expiresAt time.Time) error {
	return c.queries.Add(ctx, db.AddParams{
		UserID:     userId,
		ID:         uuid.New(),
		DeviceName: deviceName,
		Challenge:  challenge,
		ExpiresAt:  expiresAt,
	})
}

// GetForUpdate loads the latest challenge for the given username/device in a
// transaction with row-level locking, invokes challengerValidator to verify it,
// and atomically marks the attempt as success or failure. On success it commits
// the transaction; on failure it returns an appropriate error.
func (c Challenges) GetForUpdate(ctx context.Context, username, deviceName string, challengerValidator func(auth.ChallengeInfo) bool) error {
	tx, err := c.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	qtx := c.queries.WithTx(tx)

	resp, err := qtx.GetForUpdate(ctx, db.GetForUpdateParams{Username: username, DeviceName: deviceName})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrChallengeNotFound
		}
		return fmt.Errorf("failed to get challenge: %w", err)
	}

	info := auth.ChallengeInfo{
		Challenge:        resp.Challenge,
		Attempts:         resp.Attempts,
		AuthKey:          resp.AuthKey,
		AuthKeyAlgorithm: crypto.AuthKeyAlgorithm(resp.AuthKeyAlgorithm),
	}

	var updater func(context.Context, uuid.UUID) error
	if challengerValidator(info) {
		updater = qtx.UpdateWithSuccess
	} else {
		updater = qtx.UpdateWithFailure
	}

	err = updater(ctx, resp.ChallengeID)
	if err != nil {
		return fmt.Errorf("failed to update challenge: %w", err)
	}

	return tx.Commit(ctx)
}
