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
		LogFile:     getEnv("LOG_FILE", "logs/currency-service.log"),
		ServiceName: "currency-service",
	}
	logger.Init(loggerConfig)
	appLogger := logger.GetLogger()

	// Initialize metrics
	metrics.Init("currency-service")

	appLogger.Info("Starting Currency Service...")

	// Create currency server
	currencyServer, err := server.NewCurrencyServer()
	if err != nil {
		appLogger.Fatal(err, "Failed to create currency server")
	}

	// Start gRPC server in goroutine
	go func() {
		if err := currencyServer.StartGRPC(); err != nil {
			appLogger.Fatal(err, "Failed to start gRPC server")
		}
	}()

	// Start HTTP server
	go func() {
		if err := currencyServer.StartHTTP(); err != nil {
			appLogger.Fatal(err, "Failed to start HTTP server")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	appLogger.Info("Currency service shutting down...")
	currencyServer.Stop()
	appLogger.Info("Currency Service stopped")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
} 