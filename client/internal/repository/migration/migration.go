package migration

import (
	"database/sql"
	"embed"
	"fmt"
	"gophkeeper/pkg/sqlite"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*
var embedMigrations embed.FS

func MigrateSQLite(dsn, hexMasterKey string) error {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)

	}

	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	if _, err := db.Exec("PRAGMA key = \"x'" + hexMasterKey + "'\";"); err != nil {
		return err
	}

	sqlite.Setup(db)
	
	goose.SetBaseFS(embedMigrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}
