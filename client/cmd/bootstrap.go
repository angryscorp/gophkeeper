package main

import (
	"gophkeeper/client/internal/config"
	"gophkeeper/client/internal/crypto"
	grpcauth "gophkeeper/client/internal/grpc/auth"
	grpcsync "gophkeeper/client/internal/grpc/sync"
	recordrepo "gophkeeper/client/internal/repository/records"
	tokenrepo "gophkeeper/client/internal/repository/tokens"
	"gophkeeper/client/internal/tui/menu"
	"gophkeeper/client/internal/usecase/auth"
	"gophkeeper/client/internal/usecase/save"
	"gophkeeper/client/internal/usecase/sync"
	"gophkeeper/pkg/buildinfo"
	pkgcrypto "gophkeeper/pkg/crypto"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func bootstrap(cfg config.Config) (*tea.Program, []func()) {
	var closeFuncs []func()

	// Crypto proxy
	cryptoProxy := crypto.New(pkgcrypto.Encrypt)

	// Repositories initialization
	tokensRepo, closeDB, err := tokenrepo.New(cfg.DBFileName)
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
	syncClient := grpcsync.New(conn)

	// Usecases
	authUsecase := auth.New(authClient, tokensRepo, cryptoProxy.SetDataKey)
	syncUsecase := sync.New(syncClient, tokensRepo)
	saveUsecase := save.New(recordrepo.New(tokensRepo.Conn, cryptoProxy.Encrypt))

	// TUI
	mainMenu := menu.New(
		authUsecase.Register,
		authUsecase.Login,
		saveUsecase,
		func() error {
			return syncUsecase.Ping()
		},
		buildinfo.New(Version, BuildTime).String,
	)
	program := tea.NewProgram(mainMenu, tea.WithAltScreen())

	return program, closeFuncs
}
