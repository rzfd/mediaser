package server

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpcServer "github.com/rzfd/mediashar/internal/grpc"
	"github.com/rzfd/mediashar/internal/utils"
	"github.com/rzfd/mediashar/pkg/pb"
)

type NotificationServer struct {
	server  *grpc.Server
	service *grpcServer.MockNotificationService
	port    string
}

func NewNotificationServer() (*NotificationServer, error) {
	// Initialize notification service
	notificationService := grpcServer.NewMockNotificationService()

	// Create gRPC server
	grpcSrv := grpc.NewServer()
	
	// Register notification service
	notificationGRPCServer := grpcServer.NewNotificationGRPCServer(notificationService)
	pb.RegisterNotificationServiceServer(grpcSrv, notificationGRPCServer)

	// Enable reflection for development
	reflection.Register(grpcSrv)

	return &NotificationServer{
		server:  grpcSrv,
		service: notificationService,
		port:    utils.GetEnv("GRPC_PORT", "9093"),
	}, nil
}

func (s *NotificationServer) Start() error {
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", s.port, err)
	}

	return s.server.Serve(lis)
}

func (s *NotificationServer) Stop() {
	s.server.GracefulStop()
}

func (s *NotificationServer) GetPort() string {
	return s.port
} 