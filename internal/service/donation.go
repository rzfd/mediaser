package service

import "github.com/rzfd/mediashar/internal/models"

type CreateDonationRequest struct {
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Message     string  `json:"message"`
	StreamerID  uint    `json:"streamer_id"`
	DonatorID   *uint   `json:"donator_id,omitempty"`
	DisplayName string  `json:"display_name"`
	IsAnonymous bool    `json:"is_anonymous"`
}

type DonationService interface {
	Create(donation *models.Donation) error
	CreateDonation(req *CreateDonationRequest) (*models.Donation, error)
	GetByID(id uint) (*models.Donation, error)
	GetByTransactionID(transactionID string) (*models.Donation, error)
	List(page, pageSize int) ([]*models.Donation, error)
	GetByDonatorID(donatorID uint, page, pageSize int) ([]*models.Donation, error)
	GetByStreamerID(streamerID uint, page, pageSize int) ([]*models.Donation, error)
	UpdateStatus(id uint, status models.PaymentStatus) error
	ProcessPayment(donationID uint, transactionID string, provider models.PaymentProvider) error
	GetLatestDonations(limit int) ([]*models.Donation, error)
	GetTotalAmountByStreamer(streamerID uint) (float64, error)
} 