package main

import (
	"os"
	"os/signal"
	"syscall"

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
		LogFile:     getEnv("LOG_FILE", "logs/payment-service.log"),
		ServiceName: "payment-service",
	}
	logger.Init(loggerConfig)
	appLogger := logger.GetLogger()

	// Initialize metrics
	metrics.Init("payment-service")

	appLogger.Info("Starting Payment Service...")

	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		appLogger.Fatal(err, "Failed to load configuration")
	}

	// Create payment server
	paymentServer, err := server.NewPaymentServer(config)
	if err != nil {
		appLogger.Fatal(err, "Failed to create payment server")
	}

	// Start server
	go func() {
		appLogger.Info("Payment Service listening", "port", paymentServer.GetPort())
		if err := paymentServer.Start(); err != nil {
			appLogger.Fatal(err, "Failed to serve")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down Payment Service...")
	paymentServer.Stop()
	appLogger.Info("Payment Service stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 