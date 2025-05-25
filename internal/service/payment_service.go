package service

import (
	"errors"

	"github.com/rzfd/mediashar/configs"
	"github.com/rzfd/mediashar/internal/models"
)

// PaymentProcessor defines the interface for processing payments
type PaymentProcessor interface {
	ProcessPayment(amount float64, currency string, description string) (string, error)
	VerifyPayment(transactionID string) (bool, error)
	RefundPayment(transactionID string) error
}

// PaymentService handles payment processing for donations
type PaymentService interface {
	InitiatePayment(donation *models.Donation, provider models.PaymentProvider) (string, error)
	VerifyPayment(transactionID string, provider models.PaymentProvider) (bool, error)
	ProcessWebhook(payload []byte, provider models.PaymentProvider) (string, error)
}

type paymentService struct {
	config             *configs.Config
	donationService    DonationService
	paypalProcessor    PaymentProcessor
	stripeProcessor    PaymentProcessor
	cryptoProcessor    PaymentProcessor
}

func NewPaymentService(
	config *configs.Config,
	donationService DonationService,
	paypalProcessor PaymentProcessor,
	stripeProcessor PaymentProcessor,
	cryptoProcessor PaymentProcessor,
) PaymentService {
	return &paymentService{
		config:             config,
		donationService:    donationService,
		paypalProcessor:    paypalProcessor,
		stripeProcessor:    stripeProcessor,
		cryptoProcessor:    cryptoProcessor,
	}
}

func (s *paymentService) InitiatePayment(donation *models.Donation, provider models.PaymentProvider) (string, error) {
	description := "Donation to " + donation.Streamer.Username

	var processor PaymentProcessor
	switch provider {
	case models.PaymentProviderPaypal:
		processor = s.paypalProcessor
	case models.PaymentProviderStripe:
		processor = s.stripeProcessor
	case models.PaymentProviderCrypto:
		processor = s.cryptoProcessor
	default:
		return "", errors.New("unsupported payment provider")
	}

	transactionID, err := processor.ProcessPayment(donation.Amount, donation.Currency, description)
	if err != nil {
		return "", err
	}

	// This would be an initial transaction ID or payment intent ID
	// The actual payment processing would happen asynchronously
	return transactionID, nil
}

func (s *paymentService) VerifyPayment(transactionID string, provider models.PaymentProvider) (bool, error) {
	var processor PaymentProcessor
	switch provider {
	case models.PaymentProviderPaypal:
		processor = s.paypalProcessor
	case models.PaymentProviderStripe:
		processor = s.stripeProcessor
	case models.PaymentProviderCrypto:
		processor = s.cryptoProcessor
	default:
		return false, errors.New("unsupported payment provider")
	}

	return processor.VerifyPayment(transactionID)
}

func (s *paymentService) ProcessWebhook(payload []byte, provider models.PaymentProvider) (string, error) {
	// This would handle webhook callbacks from payment providers
	// It would parse the payload, extract transaction ID, verify the payment,
	// and update the donation status accordingly

	// Simplified implementation for now
	switch provider {
	case models.PaymentProviderPaypal:
		// Parse PayPal webhook payload
		// Extract transaction ID and payment status
		// Update donation status
		return "paypal-transaction-id", nil
	case models.PaymentProviderStripe:
		// Parse Stripe webhook payload
		// Extract transaction ID and payment status
		// Update donation status
		return "stripe-transaction-id", nil
	case models.PaymentProviderCrypto:
		// Parse crypto payment webhook payload
		// Extract transaction ID and payment status
		// Update donation status
		return "crypto-transaction-id", nil
	default:
		return "", errors.New("unsupported payment provider")
	}
} 