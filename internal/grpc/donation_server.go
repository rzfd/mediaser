package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/service"
	"github.com/rzfd/mediashar/pkg/pb"
)

// DonationGRPCServer implements the gRPC DonationService
type DonationGRPCServer struct {
	pb.UnimplementedDonationServiceServer
	donationService service.DonationService
}

// NewDonationGRPCServer creates a new donation gRPC server
func NewDonationGRPCServer(donationService service.DonationService) *DonationGRPCServer {
	return &DonationGRPCServer{
		donationService: donationService,
	}
}

// CreateDonation creates a new donation via gRPC
func (s *DonationGRPCServer) CreateDonation(ctx context.Context, req *pb.CreateDonationRequest) (*pb.CreateDonationResponse, error) {
	// Convert gRPC request to service request
	createReq := &service.CreateDonationRequest{
		Amount:      req.Amount,
		Currency:    req.Currency,
		Message:     req.Message,
		StreamerID:  uint(req.StreamerId),
		DisplayName: req.DisplayName,
		IsAnonymous: req.IsAnonymous,
	}

	// Set donator ID if provided
	if req.DonatorId != nil {
		donatorID := uint(*req.DonatorId)
		createReq.DonatorID = &donatorID
	}

	// Create donation
	donation, err := s.donationService.CreateDonation(createReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create donation: %v", err)
	}

	// Generate transaction ID for response
	transactionID := generateTransactionID(donation.ID)

	return &pb.CreateDonationResponse{
		DonationId:    uint32(donation.ID),
		TransactionId: transactionID,
		PaymentUrl:    "", // Will be set by payment processor
		QrCodeBase64:  "", // Will be set by QRIS processor
		ExpiresAt:     timestamppb.New(time.Now().Add(15 * time.Minute)),
	}, nil
}

// GetDonation retrieves a donation by ID
func (s *DonationGRPCServer) GetDonation(ctx context.Context, req *pb.GetDonationRequest) (*pb.GetDonationResponse, error) {
	donation, err := s.donationService.GetByID(uint(req.DonationId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "donation not found: %v", err)
	}

	pbDonation := convertModelToPbDonation(donation)
	return &pb.GetDonationResponse{
		Donation: pbDonation,
	}, nil
}

// GetDonationsByStreamer retrieves donations for a specific streamer
func (s *DonationGRPCServer) GetDonationsByStreamer(ctx context.Context, req *pb.GetDonationsByStreamerRequest) (*pb.GetDonationsListResponse, error) {
	page := int(req.Page)
	pageSize := int(req.PageSize)
	
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	donations, err := s.donationService.GetByStreamerID(uint(req.StreamerId), page, pageSize)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get donations: %v", err)
	}

	var pbDonations []*pb.Donation
	for _, donation := range donations {
		pbDonations = append(pbDonations, convertModelToPbDonation(donation))
	}

	// Calculate total pages (simplified - in production, get actual count)
	totalPages := (len(pbDonations) + pageSize - 1) / pageSize

	return &pb.GetDonationsListResponse{
		Donations:   pbDonations,
		TotalCount:  int32(len(pbDonations)),
		CurrentPage: int32(page),
		TotalPages:  int32(totalPages),
	}, nil
}

// UpdateDonationStatus updates the status of a donation
func (s *DonationGRPCServer) UpdateDonationStatus(ctx context.Context, req *pb.UpdateDonationStatusRequest) (*pb.UpdateDonationStatusResponse, error) {
	status := convertPbToModelPaymentStatus(req.Status)
	
	err := s.donationService.UpdateStatus(uint(req.DonationId), status)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update donation status: %v", err)
	}

	return &pb.UpdateDonationStatusResponse{
		Success: true,
		Message: "Donation status updated successfully",
	}, nil
}

// StreamDonationEvents streams real-time donation events (placeholder)
func (s *DonationGRPCServer) StreamDonationEvents(req *pb.StreamDonationEventsRequest, stream pb.DonationService_StreamDonationEventsServer) error {
	// This is a placeholder implementation for streaming
	// In a real implementation, you would:
	// 1. Subscribe to donation events from a message queue or event store
	// 2. Stream events to the client as they occur
	
	// For now, just return a sample event
	event := &pb.DonationEvent{
		Type:      pb.EventType_EVENT_TYPE_DONATION_CREATED,
		Timestamp: timestamppb.Now(),
		Metadata: map[string]string{
			"streamer_id": string(rune(req.StreamerId)),
		},
	}

	if err := stream.Send(event); err != nil {
		return status.Errorf(codes.Internal, "failed to send event: %v", err)
	}

	return nil
}

// GetDonationStats retrieves donation statistics
func (s *DonationGRPCServer) GetDonationStats(ctx context.Context, req *pb.GetDonationStatsRequest) (*pb.GetDonationStatsResponse, error) {
	// Get total amount for streamer
	totalAmount, err := s.donationService.GetTotalAmountByStreamer(uint(req.StreamerId))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get donation stats: %v", err)
	}

	// Get recent donations to calculate count and average
	donations, err := s.donationService.GetByStreamerID(uint(req.StreamerId), 1, 1000)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get donations for stats: %v", err)
	}

	totalCount := len(donations)
	averageAmount := float64(0)
	if totalCount > 0 {
		averageAmount = totalAmount / float64(totalCount)
	}

	return &pb.GetDonationStatsResponse{
		TotalAmount:   totalAmount,
		TotalDonations: int32(totalCount),
		AverageAmount: averageAmount,
		DailyStats:    []*pb.DonationStat{}, // Placeholder for daily stats
	}, nil
}

// Helper functions

func convertModelToPbDonation(donation *models.Donation) *pb.Donation {
	pbDonation := &pb.Donation{
		Id:              uint32(donation.ID),
		Amount:          donation.Amount,
		Currency:        donation.Currency,
		Message:         donation.Message,
		StreamerId:      uint32(donation.StreamerID),
		DonatorId:       uint32(donation.DonatorID),
		DisplayName:     donation.DisplayName,
		IsAnonymous:     donation.IsAnonymous,
		Status:          convertModelToPbPaymentStatus(donation.Status),
		PaymentProvider: convertModelToPbPaymentProvider(donation.PaymentProvider),
		TransactionId:   donation.TransactionID,
		CreatedAt:       timestamppb.New(donation.CreatedAt),
		UpdatedAt:       timestamppb.New(donation.UpdatedAt),
	}

	if donation.PaymentTime != nil {
		pbDonation.PaymentTime = timestamppb.New(*donation.PaymentTime)
	}

	return pbDonation
}

func convertModelToPbPaymentStatus(status models.PaymentStatus) pb.PaymentStatus {
	switch status {
	case models.PaymentPending:
		return pb.PaymentStatus_PAYMENT_STATUS_PENDING
	case models.PaymentCompleted:
		return pb.PaymentStatus_PAYMENT_STATUS_COMPLETED
	case models.PaymentFailed:
		return pb.PaymentStatus_PAYMENT_STATUS_FAILED
	default:
		return pb.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED
	}
}

func convertPbToModelPaymentStatus(status pb.PaymentStatus) models.PaymentStatus {
	switch status {
	case pb.PaymentStatus_PAYMENT_STATUS_PENDING:
		return models.PaymentPending
	case pb.PaymentStatus_PAYMENT_STATUS_COMPLETED:
		return models.PaymentCompleted
	case pb.PaymentStatus_PAYMENT_STATUS_FAILED:
		return models.PaymentFailed
	default:
		return models.PaymentPending
	}
}

func convertModelToPbPaymentProvider(provider models.PaymentProvider) pb.PaymentProvider {
	switch provider {
	case models.PaymentProviderMidtrans:
		return pb.PaymentProvider_PAYMENT_PROVIDER_MIDTRANS
	case models.PaymentProviderPaypal:
		return pb.PaymentProvider_PAYMENT_PROVIDER_PAYPAL
	case models.PaymentProviderStripe:
		return pb.PaymentProvider_PAYMENT_PROVIDER_STRIPE
	case models.PaymentProviderQRIS:
		return pb.PaymentProvider_PAYMENT_PROVIDER_QRIS
	default:
		return pb.PaymentProvider_PAYMENT_PROVIDER_UNSPECIFIED
	}
}

func generateTransactionID(donationID uint) string {
	return "GRPC-DON-" + string(rune(donationID)) + "-" + string(rune(time.Now().Unix()))
} 