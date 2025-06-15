package adapter

import (
	"fmt"
	"time"

	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/service"
	"github.com/rzfd/mediashar/pkg/pb"
)

// PaymentServiceAdapter adapts gRPC payment service calls
type PaymentServiceAdapter struct {
	paymentClient pb.PaymentServiceClient
}

func NewPaymentServiceAdapter(paymentClient pb.PaymentServiceClient) *PaymentServiceAdapter {
	return &PaymentServiceAdapter{
		paymentClient: paymentClient,
	}
}

func (p *PaymentServiceAdapter) InitiatePayment(donation *models.Donation, provider models.PaymentProvider) (string, error) {
	return "", nil
}

func (p *PaymentServiceAdapter) VerifyPayment(transactionID string, provider models.PaymentProvider) (bool, error) {
	return false, nil
}

func (p *PaymentServiceAdapter) ProcessWebhook(payload []byte, provider models.PaymentProvider) (string, error) {
	return "", nil
}

// MidtransServiceAdapter adapts Midtrans service calls
type MidtransServiceAdapter struct {
	paymentClient pb.PaymentServiceClient
}

func NewMidtransServiceAdapter(paymentClient pb.PaymentServiceClient) *MidtransServiceAdapter {
	return &MidtransServiceAdapter{
		paymentClient: paymentClient,
	}
}

func (m *MidtransServiceAdapter) CreateSnapTransaction(req *service.MidtransPaymentRequest) (*service.MidtransPaymentResponse, error) {
	// For now, return a mock response with sandbox token
	// In production, this should call actual Midtrans API
	return &service.MidtransPaymentResponse{
		Token:       generateMockSnapToken(req.OrderID),
		RedirectURL: fmt.Sprintf("https://app.sandbox.midtrans.com/snap/v2/vtweb/%s", generateMockSnapToken(req.OrderID)),
		OrderID:     req.OrderID,
	}, nil
}

func (m *MidtransServiceAdapter) HandleNotification(notification *service.MidtransNotification) error {
	return nil
}

func (m *MidtransServiceAdapter) VerifySignature(notification *service.MidtransNotification) bool {
	return true
}

func (m *MidtransServiceAdapter) GetTransactionStatus(orderID string) (*service.MidtransNotification, error) {
	return &service.MidtransNotification{}, nil
}

func (m *MidtransServiceAdapter) ProcessDonationPayment(donation *models.Donation) (*service.MidtransPaymentResponse, error) {
	// Generate unique order ID
	orderID := fmt.Sprintf("DONATION-%d-%d", donation.ID, time.Now().Unix())

	// Create mock Snap token for testing
	snapToken := generateMockSnapToken(orderID)
	
	return &service.MidtransPaymentResponse{
		Token:       snapToken,
		RedirectURL: fmt.Sprintf("https://app.sandbox.midtrans.com/snap/v2/vtweb/%s", snapToken),
		OrderID:     orderID,
	}, nil
}

// generateMockSnapToken generates a mock Snap token for testing
func generateMockSnapToken(orderID string) string {
	// This is a mock token format that looks like real Midtrans token
	// In production, this would come from actual Midtrans API
	return fmt.Sprintf("snap-token-%s-%d", orderID, time.Now().Unix())
} 