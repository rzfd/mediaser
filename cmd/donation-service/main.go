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
	"github.com/rzfd/mediashar/internal/service"
)

func main() {
	log.Println("Starting Donation Microservice...")

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
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("Failed to connect to donation database: %v", err)
	}

	// Run database migrations for donation-related tables only
	if err := migrateDonationTables(db); err != nil {
		log.Fatalf("Failed to migrate donation tables: %v", err)
	}

	// Initialize donation-specific repositories
	donationRepo := repositoryImpl.NewDonationRepository(db)
	userRepo := repositoryImpl.NewUserRepository(db)
	userCacheRepo := repositoryImpl.NewUserCacheRepository(db)
	
	// Initialize User Service client for external API calls
	userServiceURL := getEnv("USER_SERVICE_URL", "http://localhost:8080")
	userClient := serviceImpl.NewHTTPUserServiceClient(userServiceURL)
	
	// Initialize user aggregator service (combines cache + API calls)
	userAggregator := service.NewUserAggregatorService(userCacheRepo, userClient)
	
	log.Printf("DEBUG: Created donationRepo: %v", donationRepo != nil)
	log.Printf("DEBUG: Created userRepo: %v", userRepo != nil)
	log.Printf("DEBUG: Created userAggregator: %v", userAggregator != nil)
	log.Printf("DEBUG: User Service URL: %s", userServiceURL)
	
	// Initialize donation service with user aggregator
	donationService := serviceImpl.NewDonationServiceWithUserAggregator(donationRepo, userRepo, userAggregator)
	
	log.Printf("DEBUG: Created donationService: %v", donationService != nil)

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
		log.Printf("Donation Service listening on port %s", getEnv("GRPC_PORT", "9091"))
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Donation Service...")
	grpcSrv.GracefulStop()
	log.Println("Donation Service stopped")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func migrateDonationTables(db *gorm.DB) error {
	// First, drop existing foreign key constraints if they exist
	if err := dropExistingForeignKeys(db); err != nil {
		log.Printf("Warning: Could not drop existing foreign keys: %v", err)
	}
	
	// For PostgreSQL, we don't need to disable foreign key checks like MySQL
	// Instead, we'll configure GORM to not create foreign key constraints
	
	// Create a new session with custom configuration
	migrator := db.Session(&gorm.Session{})
	
	// Only migrate donation-related tables
	err := migrator.AutoMigrate(
		&models.User{},      // User table without FK relationships 
		&models.UserCache{}, // Cached user data for performance
		&models.Donation{},  // Donation table without FK relationships
	)
	
	if err != nil {
		return err
	}
	
	log.Println("Database migration completed successfully without foreign key constraints")
	return nil
}

func dropExistingForeignKeys(db *gorm.DB) error {
	// Drop foreign key constraints that might have been created by previous migrations
	constraints := []string{
		"fk_users_donations",
		"fk_users_received",
		"fk_donations_donator",
		"fk_donations_streamer",
		"fk_donations_users",
	}
	
	for _, constraint := range constraints {
		// Try to drop the constraint, ignore error if it doesn't exist
		if err := db.Exec(fmt.Sprintf("ALTER TABLE donations DROP CONSTRAINT IF EXISTS %s", constraint)).Error; err != nil {
			log.Printf("Info: Could not drop constraint %s (might not exist): %v", constraint, err)
		} else {
			log.Printf("Dropped foreign key constraint: %s", constraint)
		}
	}
	
	// Also try to clean up any existing data that might conflict
	log.Println("Cleaning up donation table to prevent constraint conflicts...")
	
	// Don't truncate as we want to preserve data, just ensure consistency
	log.Println("Foreign key cleanup completed")
	
	return nil
} 