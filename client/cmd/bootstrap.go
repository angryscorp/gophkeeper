package main

import (
	"crypto/tls"
	"gophkeeper/client/internal/config"
	"gophkeeper/client/internal/crypto"
	grpcauth "gophkeeper/client/internal/grpc/auth"
	grpcsync "gophkeeper/client/internal/grpc/sync"
	recordrepo "gophkeeper/client/internal/repository/records"
	syncrepo "gophkeeper/client/internal/repository/sync"
	tokenrepo "gophkeeper/client/internal/repository/tokens"
	"gophkeeper/client/internal/tui/menu"
	"gophkeeper/client/internal/usecase/auth"
	"gophkeeper/client/internal/usecase/help"
	"gophkeeper/client/internal/usecase/list"
	"gophkeeper/client/internal/usecase/save"
	"gophkeeper/client/internal/usecase/sync"
	pkgcrypto "gophkeeper/pkg/crypto"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func bootstrap(cfg config.Config) (*tea.Program, []func()) {
	var closeFuncs []func()

	// Crypto proxy keeps data key inside
	cryptoProxy := crypto.New(pkgcrypto.Encrypt, pkgcrypto.Decrypt)

	// Repositories initialization
	tokensRepo, closeDB, err := tokenrepo.New(cfg.DBFileName)
	if err != nil {
		log.Fatal(err)
	}
	closeFuncs = append(closeFuncs, closeDB)

	// gRPC client connection
	var creds grpc.DialOption
	if cfg.Debug {
		creds = grpc.WithTransportCredentials(insecure.NewCredentials())
	} else {
		tlsCfg := &tls.Config{
			MinVersion: tls.VersionTLS13,
			ServerName: cfg.ServerName,
		}
		creds = grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg))
	}
	conn, err := grpc.NewClient(cfg.ServerAddr, creds)
	if err != nil {
		log.Fatal(err)
	}
	closeFuncs = append(closeFuncs, func() { _ = conn.Close() })

	// gRPC client
	authClient := grpcauth.New(conn)
	syncClient := grpcsync.New(conn)

	// Usecases
	authUsecase := auth.New(authClient, tokensRepo, cryptoProxy.SetDataKey)
	syncUsecase := sync.New(syncClient, tokensRepo, syncrepo.New(tokensRepo.Conn))
	saveUsecase := save.New(recordrepo.New(tokensRepo.Conn), cryptoProxy.Encrypt)
	listUsecase := list.New(recordrepo.New(tokensRepo.Conn), cryptoProxy.Decrypt)
	helpUsecase := help.New(Version, BuildTime)

	// TUI
	env := menu.Environment{
		RegFactory:   authUsecase.Register,
		LoginFactory: authUsecase.Login,
		DataSaver:    saveUsecase,
		SyncFactory:  syncUsecase.Sync,
		DataFactory:  listUsecase.GetAllRecords,
		HelpFactory:  helpUsecase.Help,
	}
	program := tea.NewProgram(menu.New(env), tea.WithAltScreen())
	return program, closeFuncs
}
