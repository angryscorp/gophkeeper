package main

import (
	"fmt"
	"gophkeeper/client/internal/config"
	"gophkeeper/client/internal/repository/migration"
	tokenrepo "gophkeeper/client/internal/repository/tokens/impl"
	"gophkeeper/client/internal/tui/menu"
	"os"

	_ "github.com/mattn/go-sqlite3"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	dbName      = "vault.db"
	busyTimeout = 5000
)

func main() {
	cfg := config.Config{
		DatabaseDSN: fmt.Sprintf("file:%s?_pragma=busy_timeout=%d", dbName, busyTimeout),
		Debug:       true,
	}

	if err := migration.MigrateSQLite(cfg.DatabaseDSN, "hexMasterKey"); err != nil {
		panic(err)
	}

	repo, closeDB, err := tokenrepo.New(cfg.DatabaseDSN, "hexMasterKey")
	if err != nil {
		panic(err)
	}

	defer closeDB()

	fmt.Println(repo)

	program := tea.NewProgram(menu.New(), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Printf("Something went wrong: %v", err)
		os.Exit(1)
	}
}
