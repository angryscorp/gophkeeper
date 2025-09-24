package main

import (
	"fmt"
	"gophkeeper/server/internal/config"
	"gophkeeper/server/internal/repository/migration"
	"gophkeeper/server/internal/tokens"
	"log"
	"time"
)

const (
	audience       = "gophkeeper"
	accessTokenTTL = 15 * time.Minute
)

func main() {
	// Load configuration from env
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatal(err)
	}

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
