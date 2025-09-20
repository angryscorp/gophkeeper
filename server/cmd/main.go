package main

import (
	"gophkeeper/server/internal/config"
	"gophkeeper/server/internal/repository/migration"
	"log"
)

func main() {
	cfg := config.Config{
		DatabaseDSN: "postgres://db_user:db_password@postgres:5432/gophkeeper?sslmode=disable",
		Debug:       true,
		ServerAddr:  ":8443",
	}

	// Database migrations
	if err := migration.MigratePostgres(cfg.DatabaseDSN); err != nil {
		panic(err)
	}

	// gRPC server initialization
	grpcServer, closeFuncs := bootstrap(cfg)
	defer func() {
		for _, f := range closeFuncs {
			f()
		}
	}()

	// Run the server with graceful shutdown
	if err := runWithGracefulShutdown(grpcServer, cfg.ServerAddr); err != nil {
		log.Printf("Server error: %v", err)
	}
}
