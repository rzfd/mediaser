package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/rzfd/mediashar/configs"
	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/repository/repositoryImpl"
	"github.com/rzfd/mediashar/internal/service/serviceImpl"
	grpcServer "github.com/rzfd/mediashar/internal/grpc"
	"github.com/rzfd/mediashar/pkg/pb"
)

func main() {
	log.Println("💳 Starting Payment Microservice...")

	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection for payments
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		getEnv("PAYMENT_DB_HOST", config.DB.Host),
		getEnv("PAYMENT_DB_USERNAME", config.DB.Username),
		getEnv("PAYMENT_DB_PASSWORD", config.DB.Password),
		getEnv("PAYMENT_DB_NAME", "payment_db"),
		getEnv("PAYMENT_DB_PORT", config.DB.Port))
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to payment database: %v", err)
	}

	// Run database migrations for payment-related tables
	if err := migratePaymentTables(db); err != nil {
		log.Fatalf("Failed to migrate payment tables: %v", err)
	}

	// Initialize payment-specific repositories (we'll need to create these)
	donationRepo := repositoryImpl.NewDonationRepository(db) // For updating payment status
	userRepo := repositoryImpl.NewUserRepository(db)         // For user data
	
	// Initialize payment services
	donationService := serviceImpl.NewDonationService(donationRepo, userRepo)
	paymentService := serviceImpl.NewPaymentService(config, donationService, nil, nil, nil)

	// Create gRPC server
	lis, err := net.Listen("tcp", ":"+getEnv("GRPC_PORT", "9092"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcSrv := grpc.NewServer()
	
	// Register payment service
	paymentGRPCServer := grpcServer.NewPaymentGRPCServer(paymentService)
	pb.RegisterPaymentServiceServer(grpcSrv, paymentGRPCServer)

	// Enable reflection for development
	reflection.Register(grpcSrv)

	// Graceful shutdown
	go func() {
		log.Printf("Payment Service listening on port %s", getEnv("GRPC_PORT", "9092"))
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down Payment Service...")
	grpcSrv.GracefulStop()
	log.Println("Payment Service stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func migratePaymentTables(db *gorm.DB) error {
	// Payment-related tables including donations for status updates
	return db.AutoMigrate(
		&models.Donation{}, // Replicated for payment status updates
		&models.User{},     // Reference data
	)
} 