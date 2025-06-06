package serviceImpl

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/rzfd/mediashar/configs"
	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/service"
)

type midtransService struct {
	config          *configs.Config
	snapClient      snap.Client
	donationService service.DonationService
}

func NewMidtransService(config *configs.Config, donationService service.DonationService) service.MidtransService {
	// Initialize Midtrans client
	var env midtrans.EnvironmentType
	if config.Midtrans.Environment == "production" {
		env = midtrans.Production
	} else {
		env = midtrans.Sandbox
	}

	snapClient := snap.Client{}
	snapClient.New(config.Midtrans.ServerKey, env)

	return &midtransService{
		config:          config,
		snapClient:      snapClient,
		donationService: donationService,
	}
}

func (s *midtransService) CreateSnapTransaction(req *service.MidtransPaymentRequest) (*service.MidtransPaymentResponse, error) {
	// Convert amount to integer (Midtrans expects amount in smallest currency unit)
	amount := int64(req.Amount)

	// Create Snap request
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  req.OrderID,
			GrossAmt: amount,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: req.CustomerName,
			Email: req.CustomerEmail,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    "donation",
				Price: amount,
				Qty:   1,
				Name:  req.Description,
			},
		},
	}

	// Create transaction
	snapTokenResp, err := s.snapClient.CreateTransaction(snapReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create Midtrans transaction: %w", err)
	}

	return &service.MidtransPaymentResponse{
		Token:       snapTokenResp.Token,
		RedirectURL: snapTokenResp.RedirectURL,
		OrderID:     req.OrderID,
	}, nil
}

func (s *midtransService) ProcessDonationPayment(donation *models.Donation) (*service.MidtransPaymentResponse, error) {
	// Generate unique order ID
	orderID := fmt.Sprintf("DONATION-%d-%d", donation.ID, time.Now().Unix())

	// Get donator details
	customerName := donation.DisplayName
	customerEmail := "donor@mediashar.com" // Default email

	if donation.DonatorID != 0 {
		// If we have donator info, you can fetch from database
		// For now, use display name
		customerName = donation.DisplayName
	}

	req := &service.MidtransPaymentRequest{
		OrderID:      orderID,
		Amount:       donation.Amount,
		Currency:     string(donation.Currency),
		CustomerName: customerName,
		CustomerEmail: customerEmail,
		Description:  fmt.Sprintf("Donation to %s", donation.Streamer.Username),
		CallbackURL:  "https://yourdomain.com/donation/success", // Replace with your domain
	}

	response, err := s.CreateSnapTransaction(req)
	if err != nil {
		return nil, err
	}

	// Update donation with Midtrans order ID
	donation.TransactionID = orderID
	donation.PaymentProvider = models.PaymentProviderMidtrans
	donation.Status = models.PaymentPending

	// You might want to update the donation in the database here
	// s.donationService.Update(donation)

	return response, nil
}

func (s *midtransService) HandleNotification(notification *service.MidtransNotification) error {
	// Verify signature first
	if !s.VerifySignature(notification) {
		return fmt.Errorf("invalid signature")
	}

	// Find donation by order ID
	donation, err := s.donationService.GetByTransactionID(notification.OrderID)
	if err != nil {
		return fmt.Errorf("donation not found: %w", err)
	}

	// Update donation status based on transaction status
	var newStatus models.PaymentStatus
	switch notification.TransactionStatus {
	case "capture", "settlement":
		newStatus = models.PaymentCompleted
	case "pending":
		newStatus = models.PaymentPending
	case "deny", "expire", "cancel":
		newStatus = models.PaymentFailed
	default:
		newStatus = models.PaymentPending
	}

	// Update donation status
	err = s.donationService.UpdateStatus(donation.ID, newStatus)
	if err != nil {
		return fmt.Errorf("failed to update donation status: %w", err)
	}

	return nil
}

func (s *midtransService) VerifySignature(notification *service.MidtransNotification) bool {
	// Create signature string
	signatureString := notification.OrderID + notification.StatusCode + notification.GrossAmount + s.config.Midtrans.ServerKey
	
	// Create SHA512 hash
	hash := sha512.New()
	hash.Write([]byte(signatureString))
	signature := hex.EncodeToString(hash.Sum(nil))

	return signature == notification.SignatureKey
}

func (s *midtransService) GetTransactionStatus(orderID string) (*service.MidtransNotification, error) {
	// This would typically call Midtrans API to get transaction status
	// For now, return a mock response
	return &service.MidtransNotification{
		OrderID:           orderID,
		TransactionStatus: "pending",
		StatusCode:        "201",
	}, nil
} 