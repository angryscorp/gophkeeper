package sync

import (
	"context"
	"fmt"
	commonpgx "gophkeeper/pkg/pgx"
	"gophkeeper/server/internal/domain"
	"gophkeeper/server/internal/repository/sync/db"
	"gophkeeper/server/internal/usecase/sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Sync is a repository implementation for synchronizing changes
// between clients and server. It persists and fetches records
// using PostgreSQL.
type Sync struct {
	queries *db.Queries
	pool    *pgxpool.Pool
}

// New initializes the Sync repository with a pgx pool from the given DSN.
func New(dsn string) (*Sync, func(), error) {
	pool, err := commonpgx.CreatePGXPool(dsn)
	if err != nil {
		return nil, func() {}, fmt.Errorf("failed to create pool: %w", err)
	}

	return &Sync{
		queries: db.New(pool),
		pool:    pool,
	}, pool.Close, nil
}

var _ sync.Repository = (*Sync)(nil)

// GetChanges returns a batch of changes for a user after the given cursor.
// It enforces a page size (limit) and indicates if more data is available.
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

// AddChanges stores a batch of changes for a user in a transaction.
// It returns the operation IDs of successfully persisted changes.
func (s Sync) AddChanges(ctx context.Context, username string, changes []domain.Message) ([]uuid.UUID, error) {
	tx, err := s.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	qtx := s.queries.WithTx(tx)

	res := make([]uuid.UUID, len(changes))
	for i, change := range changes {
		err := qtx.InsertChange(ctx, db.InsertChangeParams{
			Username:      username,
			ID:            change.RecordID,
			Kind:          change.Kind,
			UpdatedAtUnix: change.UpdatedAtUnix,
			Payload:       change.Payload,
			OperationID:   change.ID,
		})

		if err != nil {
			return nil, fmt.Errorf("failed to insert change: %w", err)
		}

		res[i] = change.ID
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return res, nil
}
