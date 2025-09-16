package main

import (
	"context"
	"fmt"
	"gophkeeper/client/internal/config"
	"gophkeeper/client/internal/grpc/auth"
	"gophkeeper/client/internal/repository/migration"
	tokenrepo "gophkeeper/client/internal/repository/tokens/impl"
	"gophkeeper/client/internal/tui/menu"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	dbName      = "vault.db"
	busyTimeout = 5000
)

func main() {
	// config
	cfg := config.Config{
		DatabaseDSN: fmt.Sprintf("file:%s?_pragma=busy_timeout=%d", dbName, busyTimeout),
		Debug:       true,
	}

	// migration
	if err := migration.MigrateSQLite(cfg.DatabaseDSN, "hexMasterKey"); err != nil {
		panic(err)
	}

	// token repo
	repo, closeDB, err := tokenrepo.New(cfg.DatabaseDSN, "hexMasterKey")
	if err != nil {
		panic(err)
	}

	defer closeDB()

	fmt.Println(repo)

	// grpc conn
	conn, err := grpc.NewClient(
		"localhost:8443",
		grpc.WithTransportCredentials(insecure.NewCredentials()), // dial plaintext (dev)
	)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}

	defer func() {
		_ = conn.Close()
	}()

	// grpc client
	authClient := auth.New(conn)

	// test call
	resp, err := authClient.Register(context.Background(), "test_username")
	if err != nil {
		log.Fatalf("Register error: %v", err)
	}
	fmt.Printf("Response OK: %+v\n", resp)

	// TUI
	program := tea.NewProgram(menu.New(), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Printf("Something went wrong: %v", err)
		os.Exit(1)
	}
}
