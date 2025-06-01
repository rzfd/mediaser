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
	e.Use(middleware.CORS())

	// Setup routes with authentication
	routes.SetupRoutes(e, userHandler, donationHandler, webhookHandler, authHandler, qrisHandler, platformHandler, midtransHandler, config.Auth.JWTSecret)

	// Start server
	log.Printf("Server starting on port %s", config.Server.Port)
	log.Printf("Platform integration enabled - Available endpoints:")
	log.Printf("  POST /api/content/validate")
	log.Printf("  GET  /api/platforms/supported")
	log.Printf("  POST /api/platforms/connect (auth required)")
	log.Printf("  GET  /api/platforms (auth required)")
	log.Printf("  POST /api/content (auth required)")
	log.Printf("  GET  /api/content/live")
	log.Printf("Midtrans payment integration enabled:")
	log.Printf("  POST /api/midtrans/payment/:donationId (auth required)")
	log.Printf("  POST /api/midtrans/webhook")
	log.Printf("  GET  /api/midtrans/status/:orderId")
	
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Server.Port)))
}