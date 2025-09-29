package impl

import (
	"context"
	"fmt"
	"gophkeeper/pkg/crypto"
	pkgpgx "gophkeeper/pkg/pgx"
	"gophkeeper/server/internal/repository/challenges"
	"gophkeeper/server/internal/repository/challenges/db"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Challenges struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

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

var _ challenges.Challenges = (*Challenges)(nil)

func (c Challenges) Add(ctx context.Context, userId uuid.UUID, deviceName string, challenge []byte, expiresAt time.Time) error {
	return c.queries.Add(ctx, db.AddParams{
		UserID:     userId,
		ID:         uuid.New(),
		DeviceName: deviceName,
		Challenge:  challenge,
		ExpiresAt:  expiresAt,
	})
}

func (c Challenges) GetForUpdate(ctx context.Context, username, deviceName string, challengerValidator func(challenges.ChallengeInfo) bool) error {
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
		return fmt.Errorf("failed to get challenge: %w", err)
	}

	info := challenges.ChallengeInfo{
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

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil

}
