package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/service"
	"github.com/rzfd/mediashar/pkg/pb"
)

// PaymentGRPCServer implements the gRPC PaymentService
type PaymentGRPCServer struct {
	pb.UnimplementedPaymentServiceServer
	paymentService service.PaymentService
}

// NewPaymentGRPCServer creates a new payment gRPC server
func NewPaymentGRPCServer(paymentService service.PaymentService) *PaymentGRPCServer {
	return &PaymentGRPCServer{
		paymentService: paymentService,
	}
}

// ProcessPayment processes a payment for a donation via gRPC
func (s *PaymentGRPCServer) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.ProcessPaymentResponse, error) {
	// Convert gRPC payment provider to model
	provider := convertPbToModelPaymentProvider(req.Provider)
	
	// For now, we'll create a mock donation since we need the donation object
	// In a real microservices architecture, you'd fetch this from the donation service
	donation := &models.Donation{
		ID:       uint(req.DonationId),
		Amount:   100.0, // Placeholder - would be fetched from donation service
		Currency: "IDR", // Placeholder
	}

	// Initiate payment
	transactionID, err := s.paymentService.InitiatePayment(donation, provider)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to process payment: %v", err)
	}

	return &pb.ProcessPaymentResponse{
		TransactionId: transactionID,
		PaymentUrl:    generatePaymentURL(provider, transactionID),
		QrCode:        generateQRCode(provider, transactionID),
		Status:        pb.PaymentStatus_PAYMENT_STATUS_PENDING,
	}, nil
}

// VerifyPayment verifies a payment status via gRPC
func (s *PaymentGRPCServer) VerifyPayment(ctx context.Context, req *pb.VerifyPaymentRequest) (*pb.VerifyPaymentResponse, error) {
	// Convert gRPC payment provider to model
	provider := convertPbToModelPaymentProvider(req.Provider)

	// Verify payment
	isVerified, err := s.paymentService.VerifyPayment(req.TransactionId, provider)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to verify payment: %v", err)
	}

	// Determine status based on verification
	var paymentStatus pb.PaymentStatus
	if isVerified {
		paymentStatus = pb.PaymentStatus_PAYMENT_STATUS_COMPLETED
	} else {
		paymentStatus = pb.PaymentStatus_PAYMENT_STATUS_PENDING
	}

	return &pb.VerifyPaymentResponse{
		IsVerified: isVerified,
		Status:     paymentStatus,
		Amount:     0, // Placeholder - would be fetched from payment provider
	}, nil
}

// HandleWebhook handles payment webhook notifications via gRPC
func (s *PaymentGRPCServer) HandleWebhook(ctx context.Context, req *pb.HandleWebhookRequest) (*pb.HandleWebhookResponse, error) {
	// Convert gRPC payment provider to model
	provider := convertPbToModelPaymentProvider(req.Provider)

	// Process webhook
	transactionID, err := s.paymentService.ProcessWebhook(req.Payload, provider)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to process webhook: %v", err)
	}

	return &pb.HandleWebhookResponse{
		Success:       true,
		TransactionId: transactionID,
		Message:       "Webhook processed successfully",
	}, nil
}

// Helper functions

func convertPbToModelPaymentProvider(provider pb.PaymentProvider) models.PaymentProvider {
	switch provider {
	case pb.PaymentProvider_PAYMENT_PROVIDER_MIDTRANS:
		return models.PaymentProviderMidtrans
	case pb.PaymentProvider_PAYMENT_PROVIDER_PAYPAL:
		return models.PaymentProviderPaypal
	case pb.PaymentProvider_PAYMENT_PROVIDER_STRIPE:
		return models.PaymentProviderStripe
	case pb.PaymentProvider_PAYMENT_PROVIDER_QRIS:
		return models.PaymentProviderQRIS
	case pb.PaymentProvider_PAYMENT_PROVIDER_CRYPTO:
		return models.PaymentProviderCrypto
	default:
		return models.PaymentProviderMidtrans // Default fallback
	}
}

func generatePaymentURL(provider models.PaymentProvider, transactionID string) string {
	// Generate payment URL based on provider
	// This is a placeholder implementation
	switch provider {
	case models.PaymentProviderMidtrans:
		return "https://app.sandbox.midtrans.com/snap/v1/transactions/" + transactionID
	case models.PaymentProviderPaypal:
		return "https://www.sandbox.paypal.com/checkoutnow?token=" + transactionID
	case models.PaymentProviderStripe:
		return "https://checkout.stripe.com/pay/" + transactionID
	default:
		return ""
	}
}

func generateQRCode(provider models.PaymentProvider, transactionID string) string {
	// Generate QR code for QRIS payments
	// This is a placeholder implementation
	if provider == models.PaymentProviderQRIS {
		return "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8/5+hHgAHggJ/PchI7wAAAABJRU5ErkJggg==" // Placeholder base64
	}
	return ""
} 