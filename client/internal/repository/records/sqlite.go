package records

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"gophkeeper/client/internal/domain"
	"gophkeeper/client/internal/repository/records/db"
	"gophkeeper/client/internal/usecase/list"
	"gophkeeper/client/internal/usecase/save"

	"github.com/google/uuid"
)

// Repository provides access to user records stored in the local database.
// It implements both save.Repository and list.Repository use cases.
// Internally, it uses sqlc-generated queries and maintains a dbFactory
// for creating new DB connections.
type Repository struct {
	queries   *db.Queries
	conn      *sql.DB
	dbFactory func() (*sql.DB, error)
	mu        sync.Mutex
}

// New creates a new Repository using the given dbFactory function.
// The actual connection and sqlc queries are initialized lazily
// on the first operation.
func New(
	dbFactory func() (*sql.DB, error),
) *Repository {
	return &Repository{
		dbFactory: dbFactory,
	}
}

var _ save.Repository = (*Repository)(nil)

// Save inserts or updates a record in the local database and enqueues
// it into the outbox for later synchronization with the server.
// Both operations are executed within a single transaction.
func (r *Repository) Save(ctx context.Context, kind domain.UserDataKind, id uuid.UUID, payload []byte) error {
	if err := r.ensure(); err != nil {
		return err
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

// GetAll returns all stored records from the local database,
// mapped into RawRecord domain objects for further processing or display.
func (r *Repository) GetAll(ctx context.Context) ([]list.RawRecord, error) {
	if err := r.ensure(); err != nil {
		return nil, err
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

func (r *Repository) ensure() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.queries != nil {
		return nil
	}

	conn, err := r.dbFactory()
	if err != nil {
		return err
	}

	r.conn = conn
	r.queries = db.New(conn)
	return nil
}
