package sync

import (
	"context"
	"fmt"
	"gophkeeper/pkg/pgx"
	"gophkeeper/server/internal/domain"
	"gophkeeper/server/internal/repository/sync/db"
	"gophkeeper/server/internal/usecase/sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Sync struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

func New(dsn string) (*Sync, func(), error) {
	pool, err := pgx.CreatePGXPool(dsn)
	if err != nil {
		return nil, func() {}, fmt.Errorf("failed to create pool: %w", err)
	}

	return &Sync{
		queries: db.New(pool),
		pool:    pool,
	}, pool.Close, nil
}

var _ sync.Repository = (*Sync)(nil)

func (s Sync) GetChanges(ctx context.Context, username string, cursor int64, limit int32) (*sync.PullResponse, error) {
	rows, err := s.queries.GetChanges(ctx, db.GetChangesParams{
		Username:  username,
		ServerSeq: cursor,
		Limit:     limit + 1,
	})
	if err != nil {
		return nil, err
	}

	hasMore := int32(len(rows)) > limit
	if hasMore {
		rows = rows[:limit]
	}

	res := make([]domain.Message, len(rows))
	nextCursor := cursor
	for i, row := range rows {
		res[i] = domain.Message{
			ID:            row.OperationID,
			RecordID:      row.ID,
			Kind:          row.Kind,
			UpdatedAtUnix: row.UpdatedAtUnix,
			Payload:       row.Payload,
		}
		nextCursor = row.ServerSeq
	}

	return &sync.PullResponse{
		Changes:    res,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}
