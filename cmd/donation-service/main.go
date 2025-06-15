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
		LogFile:     getEnv("LOG_FILE", "logs/donation-service.log"),
		ServiceName: "donation-service",
	}
	logger.Init(loggerConfig)
	appLogger := logger.GetLogger()

	// Initialize metrics
	metrics.Init("donation-service")

	appLogger.Info("Starting Donation Service...")

	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		appLogger.Fatal(err, "Failed to load configuration")
	}

	// Create donation server
	donationServer, err := server.NewDonationServer(config)
	if err != nil {
		appLogger.Fatal(err, "Failed to create donation server")
	}

	// Start server
	go func() {
		appLogger.Info("Donation Service listening", "port", donationServer.GetPort())
		if err := donationServer.Start(); err != nil {
			appLogger.Fatal(err, "Failed to serve")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down Donation Service...")
	donationServer.Stop()
	appLogger.Info("Donation Service stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 