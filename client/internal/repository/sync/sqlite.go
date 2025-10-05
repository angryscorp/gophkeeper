package sync

import (
	"database/sql"
	"gophkeeper/client/internal/repository/sync/db"
)

type Repository struct {
	queries   *db.Queries
	dbFactory func() (*sql.DB, error)
}

func New(
	dbFactory func() (*sql.DB, error),
) *Repository {
	return &Repository{
		dbFactory: dbFactory,
	}
}
