package grpc

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/rzfd/mediashar/internal/adapter"
	"github.com/rzfd/mediashar/internal/service/serviceImpl"
	"github.com/rzfd/mediashar/pkg/logger"
	"github.com/rzfd/mediashar/pkg/pb"
)

// MediaShareServer represents the gRPC server
type MediaShareServer struct {
	grpcServer *grpc.Server
	listener   net.Listener
	port       string
	logger     *logger.Logger
}

// NewMediaShareServer creates a new gRPC server instance
func NewMediaShareServer(port string, mediaShareService serviceImpl.MediaShareService) (*MediaShareServer, error) {
	appLogger := logger.GetLogger()
	
	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register services
	converter := adapter.NewMediaShareConverter()
	mediaShareServer := &MediaShareGRPCHandler{
		service:   mediaShareService,
		converter: converter,
	}
	pb.RegisterMediaShareServiceServer(grpcServer, mediaShareServer)

	// Register health service
	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// Enable reflection for debugging
	reflection.Register(grpcServer)

	// Create listener
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, err
	}

	return &MediaShareServer{
		grpcServer: grpcServer,
		listener:   listener,
		port:       port,
		logger:     appLogger,
	}, nil
}

// Start starts the gRPC server
func (s *MediaShareServer) Start() error {
	s.logger.Info("Media Share Service started", "port", s.port)
	return s.grpcServer.Serve(s.listener)
}

// Stop gracefully stops the server
func (s *MediaShareServer) Stop() {
	s.logger.Info("Shutting down Media Share Service...")
	s.grpcServer.GracefulStop()
	s.logger.Info("Media Share Service stopped gracefully")
}

// WaitForShutdown waits for interrupt signal and gracefully shuts down
func (s *MediaShareServer) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	s.Stop()
} 