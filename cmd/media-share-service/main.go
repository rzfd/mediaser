package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/rzfd/mediashar/internal/config"
	"github.com/rzfd/mediashar/internal/grpc"
	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/repository/repositoryImpl"
	"github.com/rzfd/mediashar/internal/service/serviceImpl"
	"github.com/rzfd/mediashar/pkg/logger"
)

func main() {
	// Load configuration
	cfg := config.LoadMediaShareConfig()

	// Initialize logger
	logger.Init(cfg.Logger)
	appLogger := logger.GetLogger()

	appLogger.Info("Starting Media Share Service...")

	// Initialize database
	db, err := initDatabase(cfg.Database)
	if err != nil {
		appLogger.Fatal(err, "Failed to initialize database")
	}

	// Initialize services
	mediaShareService := initServices(db)

	// Create and start server
	server, err := grpc.NewMediaShareServer(cfg.Server.GRPCPort, mediaShareService)
	if err != nil {
		appLogger.Fatal(err, "Failed to create server")
	}

	// Start server in background and wait for shutdown
	go func() {
		if err := server.Start(); err != nil {
			appLogger.Fatal(err, "Failed to start server")
		}
	}()

	// Wait for shutdown signal
	server.WaitForShutdown()
}

// initDatabase initializes and migrates the database
func initDatabase(cfg config.DatabaseConfig) (*gorm.DB, error) {
	appLogger := logger.GetLogger()

	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	appLogger.Info("Database connected successfully")

	// Run migrations
	if err := db.AutoMigrate(
		&models.MediaShareSettings{},
		&models.MediaShare{},
	); err != nil {
		return nil, err
	}

	appLogger.Info("Database migrations completed")
	return db, nil
}

// initServices initializes the service layer
func initServices(db *gorm.DB) serviceImpl.MediaShareService {
	mediaShareRepo := repositoryImpl.NewMediaShareRepository(db)
	return serviceImpl.NewMediaShareService(mediaShareRepo)
}

 