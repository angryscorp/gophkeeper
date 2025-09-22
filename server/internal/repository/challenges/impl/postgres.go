package impl

import (
	"context"
	"fmt"
	"gophkeeper/pkg/pgx"
	"gophkeeper/server/internal/repository/challenges"
	"gophkeeper/server/internal/repository/challenges/db"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Challenges struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func New(dsn string) (*Challenges, func(), error) {
	pool, err := pgx.CreatePGXPool(dsn)
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
