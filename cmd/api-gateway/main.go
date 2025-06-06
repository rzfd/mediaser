package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/rzfd/mediashar/configs"
	"github.com/rzfd/mediashar/internal/handler"
	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/repository/repositoryImpl"
	"github.com/rzfd/mediashar/internal/routes"
	"github.com/rzfd/mediashar/internal/service"
	"github.com/rzfd/mediashar/internal/service/serviceImpl"
	"github.com/rzfd/mediashar/pkg/pb"
)

type APIGateway struct {
	donationClient     pb.DonationServiceClient
	paymentClient      pb.PaymentServiceClient
	notificationClient pb.NotificationServiceClient
	echo               *echo.Echo
	config             *configs.Config
}

func main() {
	log.Println("Starting API Gateway...")

	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection for gateway (user management, auth, etc.)
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

	// Run migrations for gateway-specific tables
	if err := migrateGatewayTables(db); err != nil {
		log.Fatalf("Failed to migrate gateway tables: %v", err)
	}

	// Connect to microservices
	donationConn, err := grpc.Dial(getEnv("DONATION_SERVICE_URL", "localhost:9091"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to donation service: %v", err)
	}
	defer donationConn.Close()

	paymentConn, err := grpc.Dial(getEnv("PAYMENT_SERVICE_URL", "localhost:9092"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to payment service: %v", err)
	}
	defer paymentConn.Close()

	notificationConn, err := grpc.Dial(getEnv("NOTIFICATION_SERVICE_URL", "localhost:9093"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to notification service: %v", err)
	}
	defer notificationConn.Close()

	// Create API Gateway
	gateway := &APIGateway{
		donationClient:     pb.NewDonationServiceClient(donationConn),
		paymentClient:      pb.NewPaymentServiceClient(paymentConn),
		notificationClient: pb.NewNotificationServiceClient(notificationConn),
		config:             config,
	}

	// Initialize local services (Auth, User management, Currency, Language)
	userRepo := repositoryImpl.NewUserRepository(db)
	platformRepo := repositoryImpl.NewPlatformRepository(db)
	currencyRepo := repositoryImpl.NewCurrencyRepository(db)
	languageRepo := repositoryImpl.NewLanguageRepository(db)

	userService := serviceImpl.NewUserService(userRepo)
	authService := serviceImpl.NewAuthService(config.Auth.JWTSecret, config.Auth.TokenExpiry/3600)
	platformService := serviceImpl.NewPlatformService()
	currencyService := service.NewCurrencyService(currencyRepo)
	languageService := service.NewLanguageService(languageRepo)

	// Create mock donation service for handlers that need it
	mockDonationService := &MockDonationService{gateway: gateway}
	qrisService := serviceImpl.NewQRISService("MERCHANT123", "MediaShar Donation", mockDonationService)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService, mockDonationService)
	authHandler := handler.NewAuthHandler(userService, authService)
	platformHandler := handler.NewPlatformHandler(platformService, platformRepo)
	qrisHandler := handler.NewQRISHandler(qrisService, mockDonationService)
	currencyHandler := handler.NewCurrencyHandler(currencyService)
	languageHandler := handler.NewLanguageHandler(languageService)

	// Create gateway-specific handlers
	donationHandler := gateway.NewDonationHandler()
	webhookHandler := gateway.NewWebhookHandler()
	midtransHandler := gateway.NewMidtransHandler()

	// Initialize Echo
	e := echo.New()
	gateway.echo = e

	// Global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

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

	// Setup routes
	routes.SetupRoutes(e, userHandler, donationHandler, webhookHandler, authHandler, qrisHandler, platformHandler, midtransHandler, currencyHandler, languageHandler, config.Auth.JWTSecret)

	// Health check for gateway
	e.GET("/health", gateway.HealthCheck)
	e.GET("/services/health", gateway.ServicesHealthCheck)

	// Start server
	log.Printf("API Gateway starting on port %s", config.Server.Port)
	log.Printf("Connected to:")
	log.Printf("   - Donation Service: %s", getEnv("DONATION_SERVICE_URL", "localhost:9091"))
	log.Printf("   - Payment Service: %s", getEnv("PAYMENT_SERVICE_URL", "localhost:9092"))
	log.Printf("   - Notification Service: %s", getEnv("NOTIFICATION_SERVICE_URL", "localhost:9093"))

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Server.Port)))
}

// MockDonationService implements the donation service interface for backwards compatibility
type MockDonationService struct {
	gateway *APIGateway
}

func (m *MockDonationService) Create(donation *models.Donation) error {
	// Convert to gRPC call
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcReq := &pb.CreateDonationRequest{
		Amount:        donation.Amount,
		Currency:      string(donation.Currency),
		Message:       donation.Message,
		StreamerId:    uint32(donation.StreamerID),
		DisplayName:   donation.DisplayName,
		IsAnonymous:   donation.IsAnonymous,
		PaymentMethod: "qris",
	}

	if donation.DonatorID != 0 {
		grpcReq.DonatorId = uint32(donation.DonatorID)
	}

	resp, err := m.gateway.donationClient.CreateDonation(ctx, grpcReq)
	if err != nil {
		return err
	}

	// Update the donation with the response
	donation.ID = uint(resp.DonationId)
	fmt.Printf("Created donation with ID: %d\n", donation.ID)
	return nil
}

func (m *MockDonationService) CreateDonation(req *service.CreateDonationRequest) (*models.Donation, error) {
	// Convert to gRPC call
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcReq := &pb.CreateDonationRequest{
		Amount:        req.Amount,
		Currency:      req.Currency,
		Message:       req.Message,
		StreamerId:    uint32(req.StreamerID),
		DisplayName:   req.DisplayName,
		IsAnonymous:   req.IsAnonymous,
		PaymentMethod: "qris",
	}

	if req.DonatorID != nil {
		grpcReq.DonatorId = uint32(*req.DonatorID)
	}

	resp, err := m.gateway.donationClient.CreateDonation(ctx, grpcReq)
	if err != nil {
		return nil, err
	}

	// Create donation model from response
	donation := &models.Donation{
		Amount:      req.Amount,
		Currency:    models.SupportedCurrency(req.Currency),
		Message:     req.Message,
		StreamerID:  req.StreamerID,
		DisplayName: req.DisplayName,
		IsAnonymous: req.IsAnonymous,
		Status:      models.PaymentPending,
	}
	donation.ID = uint(resp.DonationId)

	if req.DonatorID != nil {
		donation.DonatorID = *req.DonatorID
	}

	return donation, nil
}

func (m *MockDonationService) GetByID(id uint) (*models.Donation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Printf("MockDonationService.GetByID called with ID: %d\n", id)
	
	if m.gateway.donationClient == nil {
		fmt.Println("Donation client is nil")
		return nil, fmt.Errorf("donation service not available")
	}

	resp, err := m.gateway.donationClient.GetDonation(ctx, &pb.GetDonationRequest{
		DonationId: uint32(id),
	})
	if err != nil {
		fmt.Printf("gRPC GetDonation failed: %v\n", err)
		
		// IMPROVED: Try to get the real donation data first by calling donation creation endpoint
		// or return an error instead of hardcoded mock data
		return nil, fmt.Errorf("donation not found or service unavailable: %w", err)
	}

	fmt.Printf("gRPC GetDonation succeeded for ID: %d\n", id)
	fmt.Printf("Retrieved donation: Amount=%.2f, Currency=%s\n", resp.Donation.Amount, resp.Donation.Currency)

	// Convert protobuf to model
	donation := &models.Donation{
		Amount:      resp.Donation.Amount,
		Currency:    models.SupportedCurrency(resp.Donation.Currency),
		Message:     resp.Donation.Message,
		StreamerID:  uint(resp.Donation.StreamerId),
		DonatorID:   uint(resp.Donation.DonatorId),
		DisplayName: resp.Donation.DisplayName,
		IsAnonymous: resp.Donation.IsAnonymous,
	}
	donation.ID = uint(resp.Donation.Id)

	return donation, nil
}

func (m *MockDonationService) GetByTransactionID(transactionID string) (*models.Donation, error) {
	// For now, return a mock donation
	return &models.Donation{}, nil
}

func (m *MockDonationService) List(page, pageSize int) ([]*models.Donation, error) {
	// Implement gRPC call
	return []*models.Donation{}, nil
}

func (m *MockDonationService) GetByDonatorID(donatorID uint, page, pageSize int) ([]*models.Donation, error) {
	// Implement gRPC call
	return []*models.Donation{}, nil
}

func (m *MockDonationService) GetByStreamerID(streamerID uint, page, pageSize int) ([]*models.Donation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := m.gateway.donationClient.GetDonationsByStreamer(ctx, &pb.GetDonationsByStreamerRequest{
		StreamerId: uint32(streamerID),
		Page:       int32(page),
		PageSize:   int32(pageSize),
	})
	if err != nil {
		return nil, err
	}

	var donations []*models.Donation
	for _, pbDonation := range resp.Donations {
		donation := &models.Donation{
			Amount:      pbDonation.Amount,
			Currency:    models.SupportedCurrency(pbDonation.Currency),
			Message:     pbDonation.Message,
			StreamerID:  uint(pbDonation.StreamerId),
			DonatorID:   uint(pbDonation.DonatorId),
			DisplayName: pbDonation.DisplayName,
			IsAnonymous: pbDonation.IsAnonymous,
		}
		donation.ID = uint(pbDonation.Id)
		donations = append(donations, donation)
	}

	return donations, nil
}

func (m *MockDonationService) UpdateStatus(id uint, status models.PaymentStatus) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var pbStatus pb.PaymentStatus
	switch status {
	case models.PaymentPending:
		pbStatus = pb.PaymentStatus_PAYMENT_STATUS_PENDING
	case models.PaymentCompleted:
		pbStatus = pb.PaymentStatus_PAYMENT_STATUS_COMPLETED
	case models.PaymentFailed:
		pbStatus = pb.PaymentStatus_PAYMENT_STATUS_FAILED
	default:
		pbStatus = pb.PaymentStatus_PAYMENT_STATUS_PENDING
	}

	_, err := m.gateway.donationClient.UpdateDonationStatus(ctx, &pb.UpdateDonationStatusRequest{
		DonationId: uint32(id),
		Status:     pbStatus,
	})

	return err
}

func (m *MockDonationService) ProcessPayment(donationID uint, transactionID string, provider models.PaymentProvider) error {
	// Implement via payment service gRPC
	return nil
}

func (m *MockDonationService) GetLatestDonations(limit int) ([]*models.Donation, error) {
	// Implement gRPC call
	return []*models.Donation{}, nil
}

func (m *MockDonationService) GetTotalAmountByStreamer(streamerID uint) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := m.gateway.donationClient.GetDonationStats(ctx, &pb.GetDonationStatsRequest{
		StreamerId: uint32(streamerID),
	})
	if err != nil {
		return 0, err
	}

	return resp.TotalAmount, nil
}

// Gateway handlers
func (gw *APIGateway) NewDonationHandler() *handler.DonationHandler {
	mockService := &MockDonationService{gateway: gw}
	return handler.NewDonationHandler(mockService)
}

func (gw *APIGateway) NewWebhookHandler() *handler.WebhookHandler {
	mockPaymentService := &MockPaymentService{gateway: gw}
	return handler.NewWebhookHandler(mockPaymentService)
}

func (gw *APIGateway) NewMidtransHandler() *handler.MidtransHandler {
	// Create a mock payment service that delegates to payment microservice
	mockMidtransService := &MockMidtransService{gateway: gw}
	mockDonationService := &MockDonationService{gateway: gw}
	return handler.NewMidtransHandler(mockMidtransService, mockDonationService)
}

// MockPaymentService for compatibility
type MockPaymentService struct {
	gateway *APIGateway
}

func (m *MockPaymentService) InitiatePayment(donation *models.Donation, provider models.PaymentProvider) (string, error) {
	// Implement via gRPC
	return "mock-transaction-id", nil
}

func (m *MockPaymentService) VerifyPayment(transactionID string, provider models.PaymentProvider) (bool, error) {
	// Implement via gRPC
	return true, nil
}

func (m *MockPaymentService) ProcessWebhook(payload []byte, provider models.PaymentProvider) (string, error) {
	// Implement via gRPC
	return "mock-transaction-id", nil
}

// MockMidtransService for compatibility
type MockMidtransService struct {
	gateway *APIGateway
}

func (m *MockMidtransService) CreateSnapTransaction(req *service.MidtransPaymentRequest) (*service.MidtransPaymentResponse, error) {
	// Implement via payment service gRPC
	return &service.MidtransPaymentResponse{
		Token:       "mock-token",
		RedirectURL: "https://app.sandbox.midtrans.com/snap/v1/transactions/mock-token",
		OrderID:     req.OrderID,
	}, nil
}

func (m *MockMidtransService) HandleNotification(notification *service.MidtransNotification) error {
	// Implement via payment service gRPC
	return nil
}

func (m *MockMidtransService) VerifySignature(notification *service.MidtransNotification) bool {
	// Implement via payment service gRPC
	return true
}

func (m *MockMidtransService) GetTransactionStatus(orderID string) (*service.MidtransNotification, error) {
	// Implement via payment service gRPC
	return &service.MidtransNotification{
		TransactionStatus: "pending",
		OrderID:           orderID,
	}, nil
}

func (m *MockMidtransService) ProcessDonationPayment(donation *models.Donation) (*service.MidtransPaymentResponse, error) {
	fmt.Printf("MockMidtransService.ProcessDonationPayment called for donation ID: %d\n", donation.ID)
	
	// Load config for Midtrans credentials
	config, err := configs.LoadConfig()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	
	// Create real Midtrans service
	realMidtransService := serviceImpl.NewMidtransService(config, &MockDonationService{gateway: m.gateway})
	
	// Call the real Midtrans API
	response, err := realMidtransService.ProcessDonationPayment(donation)
	if err != nil {
		fmt.Printf("Real Midtrans API call failed: %v\n", err)
		// Still try to return the response from real API
		return response, nil
	}
	
	fmt.Printf("Real Midtrans API call succeeded - Token: %s\n", response.Token[:20]+"...")
	return response, nil
}

func (gw *APIGateway) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"service": "api-gateway",
		"version": "1.0.0",
	})
}

func (gw *APIGateway) ServicesHealthCheck(c echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	status := map[string]string{
		"gateway": "healthy",
	}

	// Check donation service
	_, err := gw.donationClient.GetDonationStats(ctx, &pb.GetDonationStatsRequest{StreamerId: 1})
	if err != nil {
		status["donation_service"] = "unhealthy: " + err.Error()
	} else {
		status["donation_service"] = "healthy"
	}

	// Check payment service
	_, err = gw.paymentClient.VerifyPayment(ctx, &pb.VerifyPaymentRequest{
		TransactionId: "health-check",
		Provider:      pb.PaymentProvider_PAYMENT_PROVIDER_MIDTRANS,
	})
	if err != nil {
		status["payment_service"] = "unhealthy: " + err.Error()
	} else {
		status["payment_service"] = "healthy"
	}

	// Check notification service
	_, err = gw.notificationClient.SendDonationNotification(ctx, &pb.SendNotificationRequest{
		UserId:  1,
		Title:   "Health Check",
		Message: "System health check",
	})
	if err != nil {
		status["notification_service"] = "unhealthy: " + err.Error()
	} else {
		status["notification_service"] = "healthy"
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":   "services_status",
		"services": status,
	})
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func migrateGatewayTables(db *gorm.DB) error {
	// Gateway only needs user management and platform tables
	return db.AutoMigrate(
		&models.User{},
		&models.StreamingPlatform{},
		&models.StreamingContent{},
		&models.ContentDonation{},
	)
} 