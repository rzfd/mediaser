package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rzfd/mediashar/configs"
	"github.com/rzfd/mediashar/internal/handler"
	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/repository"
	"github.com/rzfd/mediashar/internal/routes"
	"github.com/rzfd/mediashar/internal/service"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		config.DB.Host,
		config.DB.Username,
		config.DB.Password,
		config.DB.Name,
		config.DB.Port)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run database migrations
	if err := models.MigrateDB(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize repository layer
	donationRepo := repository.NewDonationRepository(db)
	userRepo := repository.NewUserRepository(db)
	
	// Initialize service layer
	donationService := service.NewDonationService(donationRepo)
	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(config.Auth.JWTSecret, config.Auth.TokenExpiry/3600) // Convert seconds to hours
	qrisService := service.NewQRISService("MERCHANT123", "MediaShar Donation", donationService)
	
	// Initialize payment service (with nil processors for now)
	paymentService := service.NewPaymentService(config, donationService, nil, nil, nil)
	
	// Initialize handler layer
	donationHandler := handler.NewDonationHandler(donationService)
	userHandler := handler.NewUserHandler(userService, donationService)
	authHandler := handler.NewAuthHandler(userService, authService)
	webhookHandler := handler.NewWebhookHandler(paymentService)
	qrisHandler := handler.NewQRISHandler(qrisService, donationService)

	// Initialize Echo
	e := echo.New()
	
	// Global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Setup routes with authentication
	routes.SetupRoutes(e, userHandler, donationHandler, webhookHandler, authHandler, qrisHandler, config.Auth.JWTSecret)

	// Start server
	log.Printf("Server starting on port %s", config.Server.Port)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Server.Port)))
}