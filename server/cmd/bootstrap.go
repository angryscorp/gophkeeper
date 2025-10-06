package main

import (
	grpcauth "gophkeeper/pkg/grpc/auth"
	grpcsync "gophkeeper/pkg/grpc/sync"
	"gophkeeper/server/internal/config"
	serverauth "gophkeeper/server/internal/grpc/auth"
	serversync "gophkeeper/server/internal/grpc/sync"
	challengesrepo "gophkeeper/server/internal/repository/challenges"
	usersrepo "gophkeeper/server/internal/repository/users"
	"gophkeeper/server/internal/tokens"
	"gophkeeper/server/internal/usecase/auth"
	"gophkeeper/server/internal/usecase/sync"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func bootstrap(cfg config.Config) (*grpc.Server, []func()) {
	var closeFuncs []func()

	// PEM keys
	publicKey, err := tokens.LoadPublicKey(cfg.PublicKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := tokens.LoadPrivateKey(cfg.PrivateKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	// Tokens
	signer := tokens.NewSigner(privateKey, audience, accessTokenTTL)
	verifier := tokens.NewVerifier(publicKey, audience)

	// Repositories initialization
	repoUsers, closeDB, err := usersrepo.New(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	closeFuncs = append(closeFuncs, closeDB)

	repoChallenges, closeDB, err := challengesrepo.New(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	closeFuncs = append(closeFuncs, closeDB)

	// gRPC server initialization
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			serversync.UnaryAuthForSync(verifier),
			serverauth.ErrorMappingServerInterceptor(),
		),
	)
	grpcauth.RegisterAuthServiceServer(
		server,
		serverauth.New(
			auth.New(repoUsers, repoChallenges, signer),
		),
	)
	grpcsync.RegisterSyncServiceServer(
		server,
		serversync.New(sync.New()),
	)

	if cfg.Debug {
		reflection.Register(server)
	}

	return server, closeFuncs
}
