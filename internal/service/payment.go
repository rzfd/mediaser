package service

import "github.com/rzfd/mediashar/internal/models"

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