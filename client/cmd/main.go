package main

import (
	"gophkeeper/client/internal/config"
	"gophkeeper/client/internal/repository/migration"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const cfgFileName = "config.json"

func main() {
	// Get executable path
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	// Load config from a file
	configPath := filepath.Join(filepath.Dir(execPath), cfgFileName)
	cfg, err := config.LoadFromFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Database migrations
	if err := migration.MigrateSQLite(cfg.DatabaseDSN(), "hexMasterKey"); err != nil {
		log.Fatal(err)
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
