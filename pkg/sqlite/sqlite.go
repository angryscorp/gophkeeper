package sqlite

import "database/sql"

func Unlock(db *sql.DB, hexMasterKey string) error {
	if _, err := db.Exec("PRAGMA key = \"x'" + hexMasterKey + "'\";"); err != nil {
		return err
	}

	return nil
}

func Setup(db *sql.DB) {
	db.Exec("PRAGMA foreign_keys = ON;")
	db.Exec("PRAGMA journal_mode = WAL;")
	db.Exec("PRAGMA synchronous = NORMAL;")
}
