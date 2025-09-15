package main

import (
	"fmt"
	"gophkeeper/pkg/grpc/auth"

	authimpl "gophkeeper/server/internal/grpc/auth"

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

	_, closeDB, err := usersrepo.New(cfg.DatabaseDSN)
	if err != nil {
		panic(err)
	}

	defer closeDB()

	// start gRPC server
	lis, err := net.Listen("tcp", ":8443")
	if err != nil {
		panic(err)
	}

	grpcServer := grpc.NewServer()
	auth.RegisterAuthServiceServer(grpcServer, authimpl.New())

	reflection.Register(grpcServer)

	log.Printf("Auth gRPC server started")
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
