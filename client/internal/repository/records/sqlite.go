package records

import (
	"context"
	"database/sql"
	"gophkeeper/client/internal/domain"
	"gophkeeper/client/internal/repository/records/db"
	"gophkeeper/client/internal/usecase/save"
	"time"

	"github.com/google/uuid"
)

type Repository struct {
	queries   *db.Queries
	dbFactory func() (*sql.DB, error)
	encryptor func(plaintext []byte) ([]byte, error)
}

func New(
	dbFactory func() (*sql.DB, error),
	encryptor func(plaintext []byte) ([]byte, error),
) *Repository {
	return &Repository{
		dbFactory: dbFactory,
		encryptor: encryptor,
	}
}

var _ save.Repository = (*Repository)(nil)

func (r Repository) Save(ctx context.Context, kind domain.UserDataKind, id uuid.UUID, data []byte) error {
	if r.queries == nil {
		conn, err := r.dbFactory()
		if err != nil {
			return err
		}
		r.queries = db.New(conn)
	}

	payload, err := r.encryptor(data)
	if err != nil {
		return err
	}

	err = r.queries.Add(ctx, db.AddParams{
		ID:            id.String(),
		Kind:          int64(kind),
		UpdatedAtUnix: time.Now().UnixMilli(),
		Payload:       payload,
	})

	if err != nil {
		return err
	}

	return nil
}
