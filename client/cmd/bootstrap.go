package main

import (
	"gophkeeper/client/internal/config"
	"gophkeeper/client/internal/grpc/auth"
	tokenrepo "gophkeeper/client/internal/repository/tokens/impl"
	"gophkeeper/client/internal/tui/menu"
	usecaseAuth "gophkeeper/client/internal/usecase/auth"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func bootstrap(cfg config.Config) (*tea.Program, []func()) {
	var closeFuncs []func()

	// Repositories initialization
	repo, closeDB, err := tokenrepo.New(cfg.DatabaseDSN(), "hexMasterKey")
	if err != nil {
		panic(err)
	}
	closeFuncs = append(closeFuncs, closeDB)

	// grpc client conn
	conn, err := grpc.NewClient(
		cfg.ServerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // dial plaintext (dev)
	)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	closeFuncs = append(closeFuncs, func() { _ = conn.Close() })

	// grpc client
	authClient := auth.New(conn)

	// usecase
	authUsecase := usecaseAuth.NewAuth(authClient, repo)

	// TUI
	mainMenu := menu.New(func() *usecaseAuth.Auth { return authUsecase })
	program := tea.NewProgram(mainMenu, tea.WithAltScreen())

	return program, closeFuncs
}
