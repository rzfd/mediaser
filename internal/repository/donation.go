package repository

import "github.com/rzfd/mediashar/internal/models"

type DonationRepository interface {
	Create(donation *models.Donation) error
	GetByID(id uint) (*models.Donation, error)
	GetByTransactionID(transactionID string) (*models.Donation, error)
	Update(donation *models.Donation) error
	Delete(id uint) error
	List(offset, limit int) ([]*models.Donation, error)
	GetByDonatorID(donatorID uint, offset, limit int) ([]*models.Donation, error)
	GetByStreamerID(streamerID uint, offset, limit int) ([]*models.Donation, error)
	UpdateStatus(id uint, status models.PaymentStatus) error
	GetLatestDonations(limit int) ([]*models.Donation, error)
	GetTotalAmountByStreamer(streamerID uint) (float64, error)
} 