package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rzfd/mediashar/configs"
	"github.com/rzfd/mediashar/internal/server"
	"github.com/rzfd/mediashar/pkg/logger"
	"github.com/rzfd/mediashar/pkg/metrics"
)

func main() {
	// Initialize logger
	loggerConfig := logger.Config{
		Level:       getEnv("LOG_LEVEL", "info"),
		Output:      getEnv("LOG_OUTPUT", "stdout"),
		LogFile:     getEnv("LOG_FILE", "logs/api-gateway.log"),
		ServiceName: "api-gateway",
	}
	logger.Init(loggerConfig)
	appLogger := logger.GetLogger()

	// Initialize metrics
	metrics.Init("api-gateway")
	
	appLogger.Info("Starting API Gateway...")

	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		appLogger.Fatal(err, "Failed to load configuration")
	}

	// Create API Gateway server
	gateway, err := server.NewAPIGateway(config)
	if err != nil {
		appLogger.Fatal(err, "Failed to create API Gateway")
	}

	// Start server
	go func() {
		appLogger.Info("Server starting on port 8080")
		if err := gateway.Start(); err != nil {
			appLogger.Fatal(err, "Failed to start server")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := gateway.Shutdown(ctx); err != nil {
		appLogger.Fatal(err, "Failed to shutdown server")
	}

	appLogger.Info("Server stopped gracefully")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 