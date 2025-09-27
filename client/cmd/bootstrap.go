package main

import (
	"gophkeeper/client/internal/config"
	grpcauth "gophkeeper/client/internal/grpc/auth"
	tokenrepo "gophkeeper/client/internal/repository/tokens/impl"
	"gophkeeper/client/internal/tui/menu"
	"gophkeeper/client/internal/usecase/auth"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func bootstrap(cfg config.Config) (*tea.Program, []func()) {
	var closeFuncs []func()

	// Repositories initialization
	repo, closeDB, err := tokenrepo.New(cfg.DBFileName)
	if err != nil {
		panic(err)
	}
	closeFuncs = append(closeFuncs, closeDB)

	// gRPC client connection
	conn, err := grpc.NewClient(
		cfg.ServerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // dial plaintext (debug)
	)
	if err != nil {
		log.Fatalf("connect: %v", err)
	}
	closeFuncs = append(closeFuncs, func() { _ = conn.Close() })

	// gRPC client
	authClient := grpcauth.New(conn)

	// TUI
	authUsecase := auth.New(authClient, repo)
	mainMenu := menu.New(
		authUsecase.Register,
		authUsecase.Login,
	)
	program := tea.NewProgram(mainMenu, tea.WithAltScreen())

	return program, closeFuncs
}
