package server

import (
	"fmt"
	"net"
	"net/http"

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
	"github.com/rzfd/mediashar/pkg/metrics"
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

func (ps *PaymentServer) Start() error {
	// Start metrics HTTP server in background
	go ps.startMetricsServer()
	
	lis, err := net.Listen("tcp", ":"+ps.port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	return ps.server.Serve(lis)
}

func (ps *PaymentServer) Stop() {
	ps.server.GracefulStop()
}

func (ps *PaymentServer) GetPort() string {
	return ps.port
}

func (ps *PaymentServer) startMetricsServer() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", metrics.MetricsHandler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"payment-service"}`))
	})
	
	metricsPort := utils.GetEnv("METRICS_PORT", "8092")
	http.ListenAndServe(":"+metricsPort, mux)
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