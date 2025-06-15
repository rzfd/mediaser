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
		LogFile:     getEnv("LOG_FILE", "logs/notification-service.log"),
		ServiceName: "notification-service",
	}
	logger.Init(loggerConfig)
	appLogger := logger.GetLogger()

	// Initialize metrics
	metrics.Init("notification-service")

	appLogger.Info("Starting Notification Service...")

	// Create notification server
	notificationServer, err := server.NewNotificationServer()
	if err != nil {
		appLogger.Fatal(err, "Failed to create notification server")
	}

	// Start server
	go func() {
		appLogger.Info("Notification Service listening", "port", notificationServer.GetPort())
		if err := notificationServer.Start(); err != nil {
			appLogger.Fatal(err, "Failed to serve")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("Shutting down Notification Service...")
	notificationServer.Stop()
	appLogger.Info("Notification Service stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 