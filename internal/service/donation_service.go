package service

import (
	"errors"
	"time"

	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/repository"
)

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

type donationService struct {
	donationRepo repository.DonationRepository
	userRepo     repository.UserRepository
}

func NewDonationService(donationRepo repository.DonationRepository) DonationService {
	return &donationService{donationRepo: donationRepo}
}

func (s *donationService) Create(donation *models.Donation) error {
	// Validate inputs
	if donation.Amount <= 0 {
		return errors.New("donation amount must be greater than zero")
	}

	if donation.StreamerID == 0 {
		return errors.New("streamer ID is required")
	}

	// Set default status if not provided
	if donation.Status == "" {
		donation.Status = models.PaymentPending
	}

	return s.donationRepo.Create(donation)
}

func (s *donationService) CreateDonation(req *CreateDonationRequest) (*models.Donation, error) {
	// Validate inputs
	if req.Amount <= 0 {
		return nil, errors.New("donation amount must be greater than zero")
	}

	if req.StreamerID == 0 {
		return nil, errors.New("streamer ID is required")
	}

	// Create donation model
	donation := &models.Donation{
		Amount:      req.Amount,
		Currency:    req.Currency,
		Message:     req.Message,
		StreamerID:  req.StreamerID,
		DisplayName: req.DisplayName,
		IsAnonymous: req.IsAnonymous,
		Status:      models.PaymentPending,
	}

	// Set donator ID if provided (for non-anonymous donations)
	if req.DonatorID != nil {
		donation.DonatorID = *req.DonatorID
	}

	// Create donation in database
	if err := s.donationRepo.Create(donation); err != nil {
		return nil, err
	}

	return donation, nil
}

func (s *donationService) GetByID(id uint) (*models.Donation, error) {
	return s.donationRepo.GetByID(id)
}

func (s *donationService) GetByTransactionID(transactionID string) (*models.Donation, error) {
	return s.donationRepo.GetByTransactionID(transactionID)
}

func (s *donationService) List(page, pageSize int) ([]*models.Donation, error) {
	offset := (page - 1) * pageSize
	return s.donationRepo.List(offset, pageSize)
}

func (s *donationService) GetByDonatorID(donatorID uint, page, pageSize int) ([]*models.Donation, error) {
	offset := (page - 1) * pageSize
	return s.donationRepo.GetByDonatorID(donatorID, offset, pageSize)
}

func (s *donationService) GetByStreamerID(streamerID uint, page, pageSize int) ([]*models.Donation, error) {
	offset := (page - 1) * pageSize
	return s.donationRepo.GetByStreamerID(streamerID, offset, pageSize)
}

func (s *donationService) UpdateStatus(id uint, status models.PaymentStatus) error {
	return s.donationRepo.UpdateStatus(id, status)
}

func (s *donationService) ProcessPayment(donationID uint, transactionID string, provider models.PaymentProvider) error {
	donation, err := s.donationRepo.GetByID(donationID)
	if err != nil {
		return err
	}

	// Update donation with payment info
	donation.TransactionID = transactionID
	donation.PaymentProvider = provider
	donation.Status = models.PaymentCompleted
	
	now := time.Now()
	donation.PaymentTime = &now

	return s.donationRepo.Update(donation)
}

func (s *donationService) GetLatestDonations(limit int) ([]*models.Donation, error) {
	return s.donationRepo.GetLatestDonations(limit)
}

func (s *donationService) GetTotalAmountByStreamer(streamerID uint) (float64, error) {
	return s.donationRepo.GetTotalAmountByStreamer(streamerID)
} 