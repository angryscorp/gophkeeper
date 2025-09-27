package main

import (
	"fmt"
	grpcauth "gophkeeper/pkg/grpc/auth"
	"gophkeeper/server/internal/config"
	serverauth "gophkeeper/server/internal/grpc/auth"
	challengesrepo "gophkeeper/server/internal/repository/challenges/impl"
	usersrepo "gophkeeper/server/internal/repository/users/impl"
	"gophkeeper/server/internal/tokens"
	"gophkeeper/server/internal/usecase/auth"
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

	// Debug
	fmt.Printf("\nSigner: %+v\n\nVerifier: %+v\n\n", signer, verifier)

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
	grpcServer := grpc.NewServer()
	grpcauth.RegisterAuthServiceServer(
		grpcServer,
		serverauth.New(
			auth.New(repoUsers, repoChallenges, signer),
		),
	)
	if cfg.Debug {
		reflection.Register(grpcServer)
	}

	return grpcServer, closeFuncs
}
