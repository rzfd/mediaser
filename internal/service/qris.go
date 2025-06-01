package service

import (
	"time"

	"github.com/rzfd/mediashar/internal/models"
)

type QRISResponse struct {
	QRISString    string    `json:"qris_string"`
	QRCodeBase64  string    `json:"qr_code_base64"`
	ExpiryTime    time.Time `json:"expiry_time"`
	Amount        float64   `json:"amount"`
	TransactionID string    `json:"transaction_id"`
}

type QRISPaymentStatus struct {
	Status        string     `json:"status"`
	TransactionID string     `json:"transaction_id"`
	Amount        float64    `json:"amount"`
	PaidAt        *time.Time `json:"paid_at,omitempty"`
}

type QRISService interface {
	GenerateQRIS(donation *models.Donation) (*QRISResponse, error)
	ValidateQRISPayment(qrisID string) (*QRISPaymentStatus, error)
	ProcessQRISCallback(payload []byte) error
} 