package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rzfd/mediashar/internal/config"
	"github.com/rzfd/mediashar/internal/grpc"
	"github.com/rzfd/mediashar/pkg/logger"
	"github.com/rzfd/mediashar/pkg/metrics"
)

func main() {
	// Load configuration
	cfg := config.LoadMediaShareConfig()

	// Initialize logger
	logger.Init(cfg.Logger)
	appLogger := logger.GetLogger()

	// Initialize metrics
	metrics.Init("media-share-service")

	appLogger.Info("Starting Media Share Service...")

	// Create and start server
	server, err := grpc.NewMediaShareServer(cfg.Server.GRPCPort, nil)
	if err != nil {
		appLogger.Fatal(err, "Failed to create server")
	}

	// Start server in background
	go func() {
		if err := server.Start(); err != nil {
			appLogger.Fatal(err, "Failed to start server")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down Media Share Service...")
	server.Stop()
	appLogger.Info("Media Share Service stopped")
}

 