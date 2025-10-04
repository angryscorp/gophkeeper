package records

import (
	"context"
	"database/sql"
	"encoding/json"
	"gophkeeper/client/internal/domain"
	"gophkeeper/client/internal/repository/records/db"
	"gophkeeper/client/internal/usecase/save"
	"time"
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

var _ save.CredentialsRepository = (*Repository)(nil)

func (r Repository) Save(ctx context.Context, credentials domain.Credentials) error {
	if r.queries == nil {
		conn, err := r.dbFactory()
		if err != nil {
			return err
		}
		r.queries = db.New(conn)
	}

	data, err := json.Marshal(credentials)
	if err != nil {
		return err
	}

	payload, err := r.encryptor(data)
	if err != nil {
		return err
	}

	err = r.queries.Add(ctx, db.AddParams{
		ID:            credentials.ID.String(),
		Kind:          domain.UserDataKindCredentials,
		UpdatedAtUnix: time.Now().UnixMilli(),
		Payload:       payload,
	})

	if err != nil {
		return err
	}

	return nil
}
