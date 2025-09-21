package main

import (
	"gophkeeper/pkg/grpc/auth"
	"gophkeeper/server/internal/config"
	authimpl "gophkeeper/server/internal/grpc/auth"
	usersrepo "gophkeeper/server/internal/repository/users/impl"
	usecaseauth "gophkeeper/server/internal/usecase/auth"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func bootstrap(cfg config.Config) (*grpc.Server, []func()) {
	closeFuncs := []func(){}

	// Repositories initialization
	repo, closeDB, err := usersrepo.New(cfg.DatabaseDSN)
	if err != nil {
		panic(err)
	}

	closeFuncs = append(closeFuncs, closeDB)

	// Usecases initialization
	usecase := usecaseauth.New(repo)

	// GRPC server initialization
	grpcServer := grpc.NewServer()
	auth.RegisterAuthServiceServer(grpcServer, authimpl.New(usecase))

	if cfg.Debug {
		reflection.Register(grpcServer)
	}

	return grpcServer, closeFuncs
}
