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

type Repository struct {
	queries   *db.Queries
	conn      *sql.DB
	dbFactory func() (*sql.DB, error)
}

func New(
	dbFactory func() (*sql.DB, error),
) *Repository {
	return &Repository{
		dbFactory: dbFactory,
	}
}

var _ sync.Repository = (*Repository)(nil)

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
