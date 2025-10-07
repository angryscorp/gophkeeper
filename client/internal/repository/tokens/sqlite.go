package tokens

import (
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"gophkeeper/client/internal/repository/migration"
	"gophkeeper/client/internal/repository/tokens/db"
	"gophkeeper/client/internal/usecase/auth"
	"gophkeeper/client/internal/usecase/sync"
)

// Tokens manages secure storage of access tokens in an encrypted SQLite database.
// It implements both sync.Tokens and auth.Tokens use cases.
// The DB is protected with SQLCipher, unlocked using a derived data key.
type Tokens struct {
	dbFileName string
	queries    *db.Queries
	db         *sql.DB
	closeDB    func()
}

// New creates a new Tokens repository bound to the given database file name.
// It returns the repository and a close function that must be called to release resources.
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

// Conn returns the underlying *sql.DB if the database is unlocked and ready.
// Otherwise it returns an error.
func (t *Tokens) Conn() (*sql.DB, error) {
	if t.Ready() {
		return t.db, nil
	}

	return nil, errors.New("connection not ready")
}

var _ sync.Tokens = (*Tokens)(nil)

// Ready reports whether the encrypted database has been successfully unlocked
// and the queries are initialized.
func (t *Tokens) Ready() bool {
	return t.db != nil && t.queries != nil && t.closeDB != nil
}

// GetAccessToken loads the current access token from the database.
// Returns an error if the DB is not unlocked or no token is found.
func (t *Tokens) GetAccessToken(ctx context.Context) (string, error) {
	if !t.Ready() {
		return "", errors.New("db is not ready")
	}
	return t.queries.GetAccessToken(ctx)
}

var _ auth.Tokens = (*Tokens)(nil)

// Unlock opens the encrypted database using the provided data key.
// It derives a hex-encoded key, applies it via PRAGMA, verifies cipher_version,
// and runs migrations if necessary.
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

// SaveAccessToken stores the given access token in the encrypted database.
// Returns an error if the DB is not unlocked.
func (t *Tokens) SaveAccessToken(ctx context.Context, token string) error {
	if !t.Ready() {
		return errors.New("db is not ready")
	}
	return t.queries.SaveAccessToken(ctx, token)
}
