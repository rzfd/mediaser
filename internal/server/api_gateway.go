package server

import (
	"context"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/rzfd/mediashar/configs"
	"github.com/rzfd/mediashar/internal/adapter"
	"github.com/rzfd/mediashar/internal/handler"
	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/repository/repositoryImpl"
	"github.com/rzfd/mediashar/internal/routes"
	"github.com/rzfd/mediashar/internal/service"
	"github.com/rzfd/mediashar/internal/service/serviceImpl"
	"github.com/rzfd/mediashar/internal/utils"
	customMiddleware "github.com/rzfd/mediashar/internal/middleware"
	"github.com/rzfd/mediashar/pkg/logger"
	"github.com/rzfd/mediashar/pkg/metrics"
	"github.com/rzfd/mediashar/pkg/pb"
)

type APIGateway struct {
	donationClient     pb.DonationServiceClient
	paymentClient      pb.PaymentServiceClient
	notificationClient pb.NotificationServiceClient
	echo               *echo.Echo
	config             *configs.Config
}

type Handlers struct {
	UserHandler       *handler.UserHandler
	AuthHandler       *handler.AuthHandler
	PlatformHandler   *handler.PlatformHandler
	QRISHandler       *handler.QRISHandler
	CurrencyHandler   *handler.CurrencyHandler
	LanguageHandler   *handler.LanguageHandler
	MediaShareHandler *handler.MediaShareHandler
	DonationHandler   *handler.DonationHandler
	WebhookHandler    *handler.WebhookHandler
	MidtransHandler   *handler.MidtransHandler
}

func NewAPIGateway(config *configs.Config) (*APIGateway, error) {
	// Initialize database
	db, err := initDatabase(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Connect to microservices
	gateway, err := connectToMicroservices(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to microservices: %w", err)
	}

	// Initialize services and handlers
	handlers := initializeHandlers(db, gateway, config)

	// Setup Echo server
	e := setupEchoServer(handlers, config)
	gateway.echo = e
	gateway.config = config

	// Start user metrics updater
	go startUserMetricsUpdater(db)

	return gateway, nil
}

func (gw *APIGateway) Start() error {
	return gw.echo.Start(fmt.Sprintf(":%s", gw.config.Server.Port))
}

func (gw *APIGateway) Shutdown(ctx context.Context) error {
	return gw.echo.Shutdown(ctx)
}

func initDatabase(config *configs.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		config.DB.Host,
		config.DB.Username,
		config.DB.Password,
		config.DB.Name,
		config.DB.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := migrateGatewayTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

func connectToMicroservices(config *configs.Config) (*APIGateway, error) {
	appLogger := logger.GetLogger()

	// Connect to donation service
	donationURL := utils.GetEnv("DONATION_SERVICE_URL", "localhost:9091")
	donationConn, err := grpc.Dial(donationURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to donation service: %w", err)
	}
	appLogger.Info("Connected to donation service", "url", donationURL)

	// Connect to payment service
	paymentURL := utils.GetEnv("PAYMENT_SERVICE_URL", "localhost:9092")
	paymentConn, err := grpc.Dial(paymentURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment service: %w", err)
	}
	appLogger.Info("Connected to payment service", "url", paymentURL)

	// Connect to notification service
	notificationURL := utils.GetEnv("NOTIFICATION_SERVICE_URL", "localhost:9093")
	notificationConn, err := grpc.Dial(notificationURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to notification service: %w", err)
	}
	appLogger.Info("Connected to notification service", "url", notificationURL)

	return &APIGateway{
		donationClient:     pb.NewDonationServiceClient(donationConn),
		paymentClient:      pb.NewPaymentServiceClient(paymentConn),
		notificationClient: pb.NewNotificationServiceClient(notificationConn),
	}, nil
}

func initializeHandlers(db *gorm.DB, gateway *APIGateway, config *configs.Config) *Handlers {
	// Initialize repositories
	userRepo := repositoryImpl.NewUserRepository(db)
	platformRepo := repositoryImpl.NewPlatformRepository(db)
	currencyRepo := repositoryImpl.NewCurrencyRepository(db)
	languageRepo := repositoryImpl.NewLanguageRepository(db)
	mediaShareRepo := repositoryImpl.NewMediaShareRepository(db)

	// Initialize services
	userService := serviceImpl.NewUserService(userRepo)
	authService := serviceImpl.NewAuthService(config.Auth.JWTSecret, config.Auth.TokenExpiry/3600)
	platformService := serviceImpl.NewPlatformService()
	currencyService := service.NewCurrencyService(currencyRepo)
	languageService := service.NewLanguageService(languageRepo)
	mediaShareService := serviceImpl.NewMediaShareService(mediaShareRepo)

	// Create service adapters
	donationService := adapter.NewDonationServiceAdapter(gateway.donationClient)
	paymentService := adapter.NewPaymentServiceAdapter(gateway.paymentClient)
	
	// Use real Midtrans service instead of adapter
	midtransService := serviceImpl.NewMidtransService(config, donationService)
	
	qrisService := serviceImpl.NewQRISService("MERCHANT123", "MediaShar Donation", donationService)

	// Initialize handlers
	return &Handlers{
		UserHandler:       handler.NewUserHandler(userService, donationService),
		AuthHandler:       handler.NewAuthHandler(userService, authService),
		PlatformHandler:   handler.NewPlatformHandler(platformService, platformRepo),
		QRISHandler:       handler.NewQRISHandler(qrisService, donationService),
		CurrencyHandler:   handler.NewCurrencyHandler(currencyService),
		LanguageHandler:   handler.NewLanguageHandler(languageService),
		MediaShareHandler: handler.NewMediaShareHandler(mediaShareService),
		DonationHandler:   handler.NewDonationHandler(donationService),
		WebhookHandler:    handler.NewWebhookHandler(paymentService),
		MidtransHandler:   handler.NewMidtransHandler(midtransService, donationService),
	}
}

func setupEchoServer(handlers *Handlers, config *configs.Config) *echo.Echo {
	e := echo.New()

	// Global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(customMiddleware.MetricsMiddleware("api-gateway"))

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			"https://localhost:3000",
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

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]interface{}{
			"service": "mediashar-api",
			"status":  "healthy",
			"version": "1.0.0",
		})
	})
	
	// Metrics endpoint
	e.GET("/metrics", func(c echo.Context) error {
		handler := metrics.MetricsHandler()
		handler.ServeHTTP(c.Response().Writer, c.Request())
		return nil
	})

	// Setup routes
	routes.SetupRoutes(e, 
		handlers.UserHandler, 
		handlers.DonationHandler, 
		handlers.WebhookHandler, 
		handlers.AuthHandler, 
		handlers.QRISHandler, 
		handlers.PlatformHandler, 
		handlers.MidtransHandler, 
		handlers.CurrencyHandler, 
		handlers.LanguageHandler, 
		handlers.MediaShareHandler, 
		config.Auth.JWTSecret)

	return e
}

func startUserMetricsUpdater(db *gorm.DB) {
	appLogger := logger.GetLogger()
	appLogger.Info("Starting user metrics updater...")
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	updateUserMetrics := func() {
		var totalUsers int64
		if err := db.Table("users").Where("deleted_at IS NULL").Count(&totalUsers).Error; err == nil {
			metrics.GetMetrics().TotalUsersRegistered.Set(float64(totalUsers))
			appLogger.Info("User metrics updated", "total_users", totalUsers)
		} else {
			appLogger.Error(err, "Failed to update user metrics")
		}
	}
	
	// Update immediately
	updateUserMetrics()
	
	// Update every 30 seconds
	for {
		select {
		case <-ticker.C:
			updateUserMetrics()
		}
	}
}

func migrateGatewayTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.StreamingPlatform{},
		&models.StreamingContent{},
		&models.ContentDonation{},
		&models.CurrencyRate{},
		&models.CurrencyInfo{},
		&models.UserCurrencyPreference{},
		&models.LanguageConfig{},
		&models.LanguageInfo{},
		&models.UserLanguagePreference{},
		&models.MediaShare{},
		&models.MediaShareSettings{},
	)
}

 