package impl

import (
	"context"
	"database/sql"
	"errors"
	"gophkeeper/client/internal/repository/tokens"
	"gophkeeper/client/internal/repository/tokens/db"
	"gophkeeper/pkg/sqlite"
)

type Tokens struct {
	queries *db.Queries
	db      *sql.DB
}

func New(dsn, hexMasterKey string) (*Tokens, func(), error) {
	database, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, func() {}, err
	}

	if err := sqlite.Unlock(database, hexMasterKey); err != nil {
		_ = database.Close()
		return nil, func() {}, err
	}

	sqlite.Setup(database)

	return &Tokens{
		queries: db.New(database),
		db:      database,
	}, func() { _ = database.Close() }, nil
}

var _ tokens.Tokens = (*Tokens)(nil)

func (t Tokens) GetAccessToken(ctx context.Context) (string, error) {
	dbTokens, err := t.queries.GetAccessToken(ctx)
	if err != nil {
		return "", err
	}

	if len(dbTokens) == 0 {
		return "", errors.New("no tokens found")
	}

	if len(dbTokens) > 1 {
		return "", errors.New("unexpected state")
	}

	return dbTokens[0], nil
}

func (t Tokens) SaveAccessToken(ctx context.Context, token string) error {
	return t.queries.SaveAccessToken(ctx, token)
}
