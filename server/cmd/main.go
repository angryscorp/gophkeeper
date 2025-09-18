package main

import (
	"fmt"
	"gophkeeper/pkg/grpc/auth"

	authimpl "gophkeeper/server/internal/grpc/auth"
	usecaseauth "gophkeeper/server/internal/usecase/auth"

	"gophkeeper/server/internal/config"
	"gophkeeper/server/internal/repository/migration"
	usersrepo "gophkeeper/server/internal/repository/users/impl"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	fmt.Println("hello world - server")

	cfg := config.Config{
		DatabaseDSN: "postgres://db_user:db_password@postgres:5432/gophkeeper?sslmode=disable",
		Debug:       true,
	}

	if err := migration.MigratePostgres(cfg.DatabaseDSN); err != nil {
		panic(err)
	}

	repo, closeDB, err := usersrepo.New(cfg.DatabaseDSN)
	if err != nil {
		panic(err)
	}

	defer closeDB()

	// usecase
	usecase := usecaseauth.New(repo)

	// start gRPC server
	lis, err := net.Listen("tcp", ":8443")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	auth.RegisterAuthServiceServer(grpcServer, authimpl.New(usecase))

	// for debug
	reflection.Register(grpcServer)

	log.Printf("Auth gRPC server started")
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
