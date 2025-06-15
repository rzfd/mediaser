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

type DonationServer struct {
	server  *grpc.Server
	service service.DonationService
	port    string
}

func NewDonationServer(config *configs.Config) (*DonationServer, error) {
	// Initialize database
	db, err := initDonationDatabase(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize services
	donationService := initDonationServices(db)

	// Create gRPC server
	grpcSrv := grpc.NewServer()
	
	// Register donation service
	donationGRPCServer := grpcServer.NewDonationGRPCServer(donationService)
	pb.RegisterDonationServiceServer(grpcSrv, donationGRPCServer)

	// Enable reflection for development
	reflection.Register(grpcSrv)

	return &DonationServer{
		server:  grpcSrv,
		service: donationService,
		port:    utils.GetEnv("GRPC_PORT", "9091"),
	}, nil
}

func (s *DonationServer) Start() error {
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", s.port, err)
	}

	return s.server.Serve(lis)
}

func (s *DonationServer) Stop() {
	s.server.GracefulStop()
}

func (s *DonationServer) GetPort() string {
	return s.port
}

func initDonationDatabase(config *configs.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		utils.GetEnv("DONATION_DB_HOST", config.DB.Host),
		utils.GetEnv("DONATION_DB_USERNAME", config.DB.Username),
		utils.GetEnv("DONATION_DB_PASSWORD", config.DB.Password),
		utils.GetEnv("DONATION_DB_NAME", "donation_db"),
		utils.GetEnv("DONATION_DB_PORT", config.DB.Port))
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := migrateDonationTables(db); err != nil {
		return nil, err
	}

	return db, nil
}

func initDonationServices(db *gorm.DB) service.DonationService {
	// Initialize repositories
	donationRepo := repositoryImpl.NewDonationRepository(db)
	userRepo := repositoryImpl.NewUserRepository(db)
	userCacheRepo := repositoryImpl.NewUserCacheRepository(db)
	
	// Initialize user service client
	userServiceURL := utils.GetEnv("USER_SERVICE_URL", "http://localhost:8080")
	userClient := serviceImpl.NewHTTPUserServiceClient(userServiceURL)
	
	// Initialize user aggregator service
	userAggregator := service.NewUserAggregatorService(userCacheRepo, userClient)
	
	// Initialize donation service
	return serviceImpl.NewDonationServiceWithUserAggregator(donationRepo, userRepo, userAggregator)
}

func migrateDonationTables(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.UserCache{},
		&models.Donation{},
	)
} 