package server

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/rzfd/mediashar/configs"
	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/repository/repositoryImpl"
	"github.com/rzfd/mediashar/internal/service"
	"github.com/rzfd/mediashar/internal/service/serviceImpl"
	"github.com/rzfd/mediashar/internal/utils"
	grpcServer "github.com/rzfd/mediashar/internal/grpc"
	"github.com/rzfd/mediashar/pkg/pb"
)

type PaymentServer struct {
	server  *grpc.Server
	service service.PaymentService
	port    string
}

func NewPaymentServer(config *configs.Config) (*PaymentServer, error) {
	// Initialize database
	db, err := initPaymentDatabase(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize services
	paymentService := initPaymentServices(db, config)

	// Create gRPC server
	grpcSrv := grpc.NewServer()
	
	// Register payment service
	paymentGRPCServer := grpcServer.NewPaymentGRPCServer(paymentService)
	pb.RegisterPaymentServiceServer(grpcSrv, paymentGRPCServer)

	// Enable reflection for development
	reflection.Register(grpcSrv)

	return &PaymentServer{
		server:  grpcSrv,
		service: paymentService,
		port:    utils.GetEnv("GRPC_PORT", "9092"),
	}, nil
}

func (s *PaymentServer) Start() error {
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", s.port, err)
	}

	return s.server.Serve(lis)
}

func (s *PaymentServer) Stop() {
	s.server.GracefulStop()
}

func (s *PaymentServer) GetPort() string {
	return s.port
}

func initPaymentDatabase(config *configs.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		utils.GetEnv("PAYMENT_DB_HOST", config.DB.Host),
		utils.GetEnv("PAYMENT_DB_USERNAME", config.DB.Username),
		utils.GetEnv("PAYMENT_DB_PASSWORD", config.DB.Password),
		utils.GetEnv("PAYMENT_DB_NAME", "payment_db"),
		utils.GetEnv("PAYMENT_DB_PORT", config.DB.Port))
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := migratePaymentTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

func initPaymentServices(db *gorm.DB, config *configs.Config) service.PaymentService {
	// Initialize repositories
	donationRepo := repositoryImpl.NewDonationRepository(db)
	userRepo := repositoryImpl.NewUserRepository(db)
	
	// Initialize services
	donationService := serviceImpl.NewDonationService(donationRepo, userRepo)
	
	return serviceImpl.NewPaymentService(config, donationService, nil, nil, nil)
}

func migratePaymentTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Donation{},
		&models.User{},
	)
} 