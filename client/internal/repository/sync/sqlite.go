package sync

import (
	"context"
	"database/sql"
	"errors"

	"gophkeeper/client/internal/domain"
	"gophkeeper/client/internal/usecase/sync"

	"gophkeeper/client/internal/repository/sync/db"

	"github.com/google/uuid"
)

// Repository provides access to synchronization-related data
// stored in the local database (outbox, records, sync cursor).
// It implements the sync.Repository interface and is used by
// the sync use cases to stage outgoing changes, apply incoming
// changes, and track sync progress.
type Repository struct {
	queries   *db.Queries
	conn      *sql.DB
	dbFactory func() (*sql.DB, error)
}

// New creates a new Repository using the given dbFactory.
// Connections and sqlc queries are initialized lazily.
func New(
	dbFactory func() (*sql.DB, error),
) *Repository {
	return &Repository{
		dbFactory: dbFactory,
	}
}

var _ sync.Repository = (*Repository)(nil)

// GetBatch returns a batch of messages from the outbox
// up to the specified limit. These messages represent
// local changes waiting to be pushed to the server.
func (r *Repository) GetBatch(ctx context.Context, limit int64) ([]domain.Message, error) {
	if r.queries == nil {
		conn, err := r.dbFactory()
		if err != nil {
			return nil, err
		}
		r.conn = conn
		r.queries = db.New(conn)
	}

	res, err := r.queries.ListOutboxBatch(ctx, limit)
	if err != nil {
		return nil, err
	}

	messages := make([]domain.Message, len(res))
	for i, r := range res {
		operationID, err := uuid.Parse(r.OperationID)
		if err != nil {
			return nil, err
		}

		recordID, err := uuid.Parse(r.RecordID)
		if err != nil {
			return nil, err
		}

		messages[i] = domain.Message{
			ID:            operationID,
			RecordID:      recordID,
			Kind:          int32(r.Kind),
			UpdatedAtUnix: r.UpdatedAtUnix,
			Payload:       r.Payload,
		}
	}

	return messages, nil
}

// DeleteBatch removes a set of messages from the outbox
// after they have been successfully synchronized with the server.
// The operation is wrapped in a transaction.
func (r *Repository) DeleteBatch(ctx context.Context, ids []uuid.UUID) error {
	if r.queries == nil {
		conn, err := r.dbFactory()
		if err != nil {
			return err
		}
		r.conn = conn
		r.queries = db.New(conn)
	}

	tx, err := r.conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := r.queries.WithTx(tx)

	for _, id := range ids {
		err := qtx.DeleteOutbox(ctx, id.String())
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetCursor returns the current synchronization cursor,
// representing the last server_seq successfully applied locally.
func (r *Repository) GetCursor(ctx context.Context) (int64, error) {
	if r.queries == nil {
		conn, err := r.dbFactory()
		if err != nil {
			return 0, err
		}
		r.conn = conn
		r.queries = db.New(conn)
	}

	return r.queries.GetCursor(ctx)
}

// SaveChanges applies a batch of incoming changes from the server
// into the local database. For each record, the change is applied
// only if its UpdatedAtUnix is newer than the local version
// (last-write-wins conflict resolution). After applying, the
// local sync cursor is advanced to newCursor.
func (r *Repository) SaveChanges(ctx context.Context, changes []domain.Message, newCursor int64) error {
	if r.queries == nil {
		conn, err := r.dbFactory()
		if err != nil {
			return err
		}
		r.conn = conn
		r.queries = db.New(conn)
	}

	tx, err := r.conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	qtx := r.queries.WithTx(tx)

	for _, change := range changes {
		localUpdated, err := qtx.GetRecordMeta(ctx, change.RecordID.String())
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return err
		}

		if localUpdated > change.UpdatedAtUnix {
			continue
		}

		err = qtx.UpsertRecord(ctx, db.UpsertRecordParams{
			ID:            change.RecordID.String(),
			Kind:          int64(change.Kind),
			UpdatedAtUnix: change.UpdatedAtUnix,
			Payload:       change.Payload,
		})
		if err != nil {
			return err
		}
	}

	err = qtx.SetCursor(ctx, newCursor)
	if err != nil {
		return err
	}

	return tx.Commit()
}
