package records

import (
	"context"
	"database/sql"
	"gophkeeper/client/internal/domain"
	"gophkeeper/client/internal/repository/records/db"
	"gophkeeper/client/internal/usecase/list"
	"gophkeeper/client/internal/usecase/save"
	"time"

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

var _ save.Repository = (*Repository)(nil)

func (r *Repository) Save(ctx context.Context, kind domain.UserDataKind, id uuid.UUID, payload []byte) error {
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
	updatedAt := time.Now().UTC().UnixMilli()

	err = qtx.AddRecord(ctx, db.AddRecordParams{
		ID:            id.String(),
		Kind:          int64(kind),
		UpdatedAtUnix: updatedAt,
		Payload:       payload,
	})

	err = qtx.EnqueueOutbox(ctx, db.EnqueueOutboxParams{
		OperationID:   uuid.New().String(),
		RecordID:      id.String(),
		Kind:          int64(kind),
		UpdatedAtUnix: updatedAt,
		Payload:       payload,
		CreatedAtUnix: time.Now().UTC().UnixMilli(),
	})

	if err != nil {
		return err
	}

	return tx.Commit()
}

var _ list.Repository = (*Repository)(nil)

func (r *Repository) GetAll(ctx context.Context) ([]list.RawRecord, error) {
	if r.queries == nil {
		conn, err := r.dbFactory()
		if err != nil {
			return nil, err
		}
		r.queries = db.New(conn)
	}

	rows, err := r.queries.GetRecords(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]list.RawRecord, 0)
	for _, row := range rows {
		id, err := uuid.Parse(row.ID)
		if err != nil {
			return nil, err
		}

		record := list.RawRecord{
			ID:      id,
			Kind:    domain.UserDataKind(row.Kind),
			Payload: row.Payload,
		}

		result = append(result, record)
	}

	return result, nil
}
