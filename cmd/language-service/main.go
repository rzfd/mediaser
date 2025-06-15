package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/rzfd/mediashar/internal/server"
	"github.com/rzfd/mediashar/pkg/logger"
	"github.com/rzfd/mediashar/pkg/metrics"
)

func main() {
	// Initialize logger
	loggerConfig := logger.Config{
		Level:       getEnv("LOG_LEVEL", "info"),
		Output:      getEnv("LOG_OUTPUT", "stdout"),
		LogFile:     getEnv("LOG_FILE", "logs/language-service.log"),
		ServiceName: "language-service",
	}
	logger.Init(loggerConfig)
	appLogger := logger.GetLogger()

	// Initialize metrics
	metrics.Init("language-service")

	appLogger.Info("Starting Language Service...")

	// Create language server
	languageServer, err := server.NewLanguageServer()
	if err != nil {
		appLogger.Fatal(err, "Failed to create language server")
	}

	// Start gRPC server in goroutine
	go func() {
		if err := languageServer.StartGRPC(); err != nil {
			appLogger.Fatal(err, "Failed to start gRPC server")
	}
	}()

	// Start HTTP server
	go func() {
		if err := languageServer.StartHTTP(); err != nil {
			appLogger.Fatal(err, "Failed to start HTTP server")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	appLogger.Info("Language service shutting down...")
	languageServer.Stop()
	appLogger.Info("Language Service stopped")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}