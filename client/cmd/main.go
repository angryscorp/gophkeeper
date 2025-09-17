package main

import (
	"fmt"
	"gophkeeper/client/internal/config"
	"gophkeeper/client/internal/grpc/auth"
	"gophkeeper/client/internal/repository/migration"
	tokenrepo "gophkeeper/client/internal/repository/tokens/impl"
	"gophkeeper/client/internal/tui/menu"
	usecaseAuth "gophkeeper/client/internal/usecase/auth"
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
	defer func() { _ = conn.Close() }()

	// grpc client
	authClient := auth.New(conn)

	// usecase
	authUsecase := usecaseAuth.NewAuth(authClient)

	// test call
	err = authUsecase.Register("test_username")
	if err != nil {
		log.Fatalf("Register error: %v", err)
	}
	fmt.Printf("Response OK\n")

	// TUI
	mainMenu := menu.New(func() *usecaseAuth.Auth { return authUsecase })
	program := tea.NewProgram(mainMenu, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Printf("Something went wrong: %v", err)
		os.Exit(1)
	}
}
