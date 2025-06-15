package server

import (
	"fmt"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	grpcServer "github.com/rzfd/mediashar/internal/grpc"
	"github.com/rzfd/mediashar/internal/utils"
	"github.com/rzfd/mediashar/pkg/metrics"
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

func (ns *NotificationServer) Start() error {
	// Start metrics HTTP server in background
	go ns.startMetricsServer()
	
	lis, err := net.Listen("tcp", ":"+ns.port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	return ns.server.Serve(lis)
}

func (ns *NotificationServer) startMetricsServer() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", metrics.MetricsHandler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"notification-service"}`))
	})
	
	metricsPort := utils.GetEnv("METRICS_PORT", "8093")
	http.ListenAndServe(":"+metricsPort, mux)
}

func (s *NotificationServer) Stop() {
	s.server.GracefulStop()
}

func (s *NotificationServer) GetPort() string {
	return s.port
} 