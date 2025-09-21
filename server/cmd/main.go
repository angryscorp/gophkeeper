package main

import (
	"gophkeeper/server/internal/config"
	"gophkeeper/server/internal/repository/migration"
	"log"
)

func main() {
	// Load configuration from env
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	// Database migrations
	if err := migration.MigratePostgres(cfg.DatabaseDSN); err != nil {
		log.Fatal(err)
	}

	// gRPC server initialization
	grpcServer, closeFuncs := bootstrap(cfg)

	// Clean up resources after a graceful shutdown
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
