package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

func runWithGracefulShutdown(grpcServer *grpc.Server, addr string) error {
	// Listen on TCP
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	// Channel for receiving system signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Channel for receiving errors from the gRPC server
	serverErrors := make(chan error, 1)

	// Start a gRPC server in a goroutine
	go func() {
		log.Printf("Starting gRPC server on %s...", addr)
		if err := grpcServer.Serve(lis); err != nil {
			serverErrors <- fmt.Errorf("failed to serve gRPC server: %w", err)
		}
	}()

	// Wait for a system signal or an error from the gRPC server
	select {
	case err := <-serverErrors:
		return err
	case <-quit:
		log.Printf("Received shutdown signal, starting graceful shutdown...")
		return gracefulShutdown(grpcServer)
	}
}

func gracefulShutdown(grpcServer *grpc.Server) error {
	// Context for the graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Graceful shutdown
	done := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(done)
	}()

	// Wait for the graceful shutdown to finish or for the context to expire
	select {
	case <-done:
		log.Printf("gRPC server stopped gracefully")
		return nil
	case <-ctx.Done():
		// Force shutdown if the context expires
		log.Printf("Graceful shutdown timed out, forcing shutdown")
		grpcServer.Stop()
		return fmt.Errorf("graceful shutdown timed out")
	}
}
