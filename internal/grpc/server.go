package grpc

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	
	"github.com/rzfd/mediashar/internal/service"
	"github.com/rzfd/mediashar/pkg/pb"
)

// GRPCServer wraps the existing services with gRPC endpoints
type GRPCServer struct {
	donationService    service.DonationService
	paymentService     service.PaymentService
	notificationService NotificationService // Future implementation
	server             *grpc.Server
}

// NotificationService interface for future implementation
type NotificationService interface {
	SendDonationNotification(ctx context.Context, userID uint, title, message string, data map[string]string) error
	SubscribeEvents(ctx context.Context, userID uint, eventTypes []string) (<-chan *pb.DonationEvent, error)
}

// NewGRPCServer creates a new gRPC server instance
func NewGRPCServer(
	donationService service.DonationService,
	paymentService service.PaymentService,
	notificationService NotificationService,
) *GRPCServer {
	return &GRPCServer{
		donationService:     donationService,
		paymentService:      paymentService,
		notificationService: notificationService,
		server:              grpc.NewServer(),
	}
}

// Start starts the gRPC server on the specified port
func (s *GRPCServer) Start(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// Register services
	donationServer := NewDonationGRPCServer(s.donationService)
	paymentServer := NewPaymentGRPCServer(s.paymentService)
	
	pb.RegisterDonationServiceServer(s.server, donationServer)
	pb.RegisterPaymentServiceServer(s.server, paymentServer)
	
	if s.notificationService != nil {
		notificationServer := NewNotificationGRPCServer(s.notificationService)
		pb.RegisterNotificationServiceServer(s.server, notificationServer)
	}

	// Enable reflection for development
	reflection.Register(s.server)

	log.Printf("gRPC server starting on port %s", port)
	return s.server.Serve(lis)
}

// Stop gracefully stops the gRPC server
func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
}

// GetServer returns the underlying gRPC server
func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
} 