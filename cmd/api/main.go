package main

import (
	"fmt"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rzfd/mediashar/configs"
	"github.com/rzfd/mediashar/internal/handler"
	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/repository/repositoryImpl"
	"github.com/rzfd/mediashar/internal/routes"
	"github.com/rzfd/mediashar/internal/service/serviceImpl"
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
	donationRepo := repositoryImpl.NewDonationRepository(db)
	userRepo := repositoryImpl.NewUserRepository(db)
	platformRepo := repositoryImpl.NewPlatformRepository(db)
	
	// Initialize service layer
	donationService := serviceImpl.NewDonationService(donationRepo)
	userService := serviceImpl.NewUserService(userRepo)
	authService := serviceImpl.NewAuthService(config.Auth.JWTSecret, config.Auth.TokenExpiry/3600) // Convert seconds to hours
	qrisService := serviceImpl.NewQRISService("MERCHANT123", "MediaShar Donation", donationService)
	platformService := serviceImpl.NewPlatformService()
	midtransService := serviceImpl.NewMidtransService(config, donationService)
	
	// Initialize payment service (with nil processors for now)
	paymentService := serviceImpl.NewPaymentService(config, donationService, nil, nil, nil)
	
	// Initialize handler layer
	donationHandler := handler.NewDonationHandler(donationService)
	userHandler := handler.NewUserHandler(userService, donationService)
	authHandler := handler.NewAuthHandler(userService, authService)
	webhookHandler := handler.NewWebhookHandler(paymentService)
	qrisHandler := handler.NewQRISHandler(qrisService, donationService)
	platformHandler := handler.NewPlatformHandler(platformService, platformRepo)
	midtransHandler := handler.NewMidtransHandler(midtransService, donationService)

	// Initialize Echo
	e := echo.New()
	
	// Global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	
	// CORS middleware with specific configuration for frontend
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:8000",  // Frontend testing interface
			"http://localhost:3000",  // Common React dev port
			"http://localhost:3001",  // Alternative React port
			"http://127.0.0.1:8000",  // Alternative localhost format
			"http://127.0.0.1:3000",  // Alternative localhost format
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With",
		},
		AllowCredentials: true,
		ExposeHeaders: []string{
			"Content-Length",
		},
	}))

	// Setup routes with authentication
	routes.SetupRoutes(e, userHandler, donationHandler, webhookHandler, authHandler, qrisHandler, platformHandler, midtransHandler, config.Auth.JWTSecret)

	// Start server
	log.Printf("Server starting on port %s", config.Server.Port)
	
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Server.Port)))
}