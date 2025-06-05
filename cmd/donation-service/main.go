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
	log.Println("🚀 Starting Donation Microservice...")

	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection for donations
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		getEnv("DONATION_DB_HOST", config.DB.Host),
		getEnv("DONATION_DB_USERNAME", config.DB.Username),
		getEnv("DONATION_DB_PASSWORD", config.DB.Password),
		getEnv("DONATION_DB_NAME", "donation_db"),
		getEnv("DONATION_DB_PORT", config.DB.Port))
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to donation database: %v", err)
	}

	// Run database migrations for donation-related tables only
	if err := migrateDonationTables(db); err != nil {
		log.Fatalf("Failed to migrate donation tables: %v", err)
	}

	// Initialize donation-specific repositories
	donationRepo := repositoryImpl.NewDonationRepository(db)
	
	// Initialize donation service
	donationService := serviceImpl.NewDonationService(donationRepo)

	// Create gRPC server
	lis, err := net.Listen("tcp", ":"+getEnv("GRPC_PORT", "9091"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcSrv := grpc.NewServer()
	
	// Register donation service
	donationGRPCServer := grpcServer.NewDonationGRPCServer(donationService)
	pb.RegisterDonationServiceServer(grpcSrv, donationGRPCServer)

	// Enable reflection for development
	reflection.Register(grpcSrv)

	// Graceful shutdown
	go func() {
		log.Printf("✅ Donation Service listening on port %s", getEnv("GRPC_PORT", "9091"))
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down Donation Service...")
	grpcSrv.GracefulStop()
	log.Println("✅ Donation Service stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func migrateDonationTables(db *gorm.DB) error {
	// Only migrate donation-related tables
	return db.AutoMigrate(
		&models.Donation{},
		&models.User{}, // Still needed for foreign key relationships
	)
} 