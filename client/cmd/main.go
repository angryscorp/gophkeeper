package main

import (
	"fmt"
	"gophkeeper/client/internal/config"
	"gophkeeper/client/internal/repository/migration"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbName      = "vault.db"
	busyTimeout = 5000
)

func main() {
	// config
	cfg := config.Config{
		DatabaseDSN: fmt.Sprintf("file:%s?_pragma=busy_timeout=%d", dbName, busyTimeout),
		ServerAddr:  "localhost:8443",
		Debug:       true,
	}

	// Database migrations
	if err := migration.MigrateSQLite(cfg.DatabaseDSN, "hexMasterKey"); err != nil {
		panic(err)
	}

	// Bootstrap
	program, closes := bootstrap(cfg)
	defer func() {
		for _, f := range closes {
			f()
		}
	}()

	// Run the program
	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}
