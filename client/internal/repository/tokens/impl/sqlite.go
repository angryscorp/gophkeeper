package impl

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"gophkeeper/client/internal/repository/migration"
	"gophkeeper/client/internal/repository/tokens"
	"gophkeeper/client/internal/repository/tokens/db"
)

type Tokens struct {
	dbFileName string
	queries    *db.Queries
	db         *sql.DB
	closeDB    func()
}

func New(dbFileName string) (*Tokens, func(), error) {
	t := &Tokens{
		dbFileName: dbFileName,
		closeDB:    func() {},
	}
	closeFn := func() {
		if t.db != nil {
			_ = t.db.Close()
			t.db = nil
		}
		t.queries = nil
		t.closeDB = nil
	}
	t.closeDB = closeFn
	return t, closeFn, nil
}

var _ tokens.Tokens = (*Tokens)(nil)

func (t *Tokens) Ready() bool {
	return t.db != nil && t.queries != nil && t.closeDB != nil
}

func (t *Tokens) Unlock(dataKey []byte) error {
	// If DB is already open, close it
	if t.Ready() {
		t.closeDB()
	}

	// Compose DSN and open DB
	hexKey := hex.EncodeToString(dataKey)
	dsn := fmt.Sprintf("file:%s?_pragma_key=x'%s'", t.dbFileName, hexKey)
	database, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return err
	}
	// Checking that the key is correct and a database is readable
	var cipherVer string
	if err := database.QueryRow(`PRAGMA cipher_version;`).Scan(&cipherVer); err != nil {
		_ = database.Close()
		return fmt.Errorf("cipher check failed: %w", err)
	}

	t.db = database
	t.closeDB = func() { _ = database.Close() }
	t.queries = db.New(database)

	// Migration
	if err := migration.MigrateSQLite(t.db); err != nil {
		t.closeDB()
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}

func (t *Tokens) GetAccessToken(ctx context.Context) (string, error) {
	return t.queries.GetAccessToken(ctx)
}

func (t *Tokens) SaveAccessToken(ctx context.Context, token string) error {
	return t.queries.SaveAccessToken(ctx, token)
}
