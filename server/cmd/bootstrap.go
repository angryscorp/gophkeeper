package main

import (
	grpcauth "gophkeeper/pkg/grpc/auth"
	"gophkeeper/server/internal/config"
	serverauth "gophkeeper/server/internal/grpc/auth"
	usersrepo "gophkeeper/server/internal/repository/users/impl"
	"gophkeeper/server/internal/usecase/auth"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func bootstrap(cfg config.Config) (*grpc.Server, []func()) {
	var closeFuncs []func()

	// Repositories initialization
	repo, closeDB, err := usersrepo.New(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}
	closeFuncs = append(closeFuncs, closeDB)

	// gRPC server initialization
	grpcServer := grpc.NewServer()
	grpcauth.RegisterAuthServiceServer(grpcServer, serverauth.New(auth.New(repo)))
	if cfg.Debug {
		reflection.Register(grpcServer)
	}

	return grpcServer, closeFuncs
}
