package serviceImpl

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/repository"
	"github.com/rzfd/mediashar/internal/service"
)

type donationService struct {
	donationRepo    repository.DonationRepository
	userRepo        repository.UserRepository
	userAggregator  service.UserAggregatorService // User aggregator for cache + API
}

func NewDonationService(donationRepo repository.DonationRepository, userRepo repository.UserRepository) service.DonationService {
	return &donationService{
		donationRepo: donationRepo,
		userRepo:     userRepo,
	}
}

// NewDonationServiceWithUserAggregator creates donation service with user aggregator (recommended)
func NewDonationServiceWithUserAggregator(donationRepo repository.DonationRepository, userRepo repository.UserRepository, userAggregator service.UserAggregatorService) service.DonationService {
	return &donationService{
		donationRepo:   donationRepo,
		userRepo:       userRepo,
		userAggregator: userAggregator,
	}
}

func getEnv(key, defaultValue string) string {
	// This function should be imported or implemented
	return defaultValue
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

	// Ensure users exist before creating donation
	if err := s.ensureUsersExist(donation); err != nil {
		return err
	}

	return s.donationRepo.Create(donation)
}

func (s *donationService) CreateDonation(req *service.CreateDonationRequest) (*models.Donation, error) {
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

	// Ensure users exist before creating donation
	if err := s.ensureUsersExist(donation); err != nil {
		return nil, err
	}

	// Create donation in database
	if err := s.donationRepo.Create(donation); err != nil {
		return nil, err
	}

	// Populate user data for response
	if err := s.populateUserData(donation); err != nil {
		// Log the error but don't fail the donation creation
		fmt.Printf("Warning: Failed to populate user data: %v\n", err)
	}

	return donation, nil
}

// populateUserData fetches user data via User Aggregator (cache + API)
func (s *donationService) populateUserData(donation *models.Donation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Use user aggregator if available
	if s.userAggregator != nil {
		// Populate streamer data
		if donation.StreamerID > 0 {
			streamer, err := s.userAggregator.GetUser(ctx, donation.StreamerID)
			if err != nil {
				fmt.Printf("Warning: Could not fetch streamer data: %v\n", err)
			} else {
				donation.Streamer = *streamer
			}
		}

		// Populate donator data (only for non-anonymous donations)
		if donation.DonatorID > 0 && !donation.IsAnonymous {
			donator, err := s.userAggregator.GetUser(ctx, donation.DonatorID)
			if err != nil {
				fmt.Printf("Warning: Could not fetch donator data: %v\n", err)
			} else {
				donation.Donator = *donator
			}
		}
		
		return nil
	}

	// Fallback to local repository if user aggregator not available
	if donation.StreamerID > 0 {
		if streamer, err := s.userRepo.GetByID(donation.StreamerID); err == nil {
			donation.Streamer = *streamer
		}
	}

	if donation.DonatorID > 0 && !donation.IsAnonymous {
		if donator, err := s.userRepo.GetByID(donation.DonatorID); err == nil {
			donation.Donator = *donator
		}
	}

	return nil
}

// ensureUsersExist validates user existence via User Aggregator (cache + API)
func (s *donationService) ensureUsersExist(donation *models.Donation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Validate streamer exists
	if donation.StreamerID == 0 {
		return errors.New("streamer ID is required")
	}

	// Use user aggregator if available
	if s.userAggregator != nil {
		// Validate streamer exists and is actually a streamer
		if err := s.userAggregator.ValidateUser(ctx, donation.StreamerID, true); err != nil {
			return fmt.Errorf("streamer validation failed: %v", err)
		}

		// For non-anonymous donations, validate donator exists
		if !donation.IsAnonymous {
			if donation.DonatorID == 0 {
				return errors.New("donator ID is required for non-anonymous donations")
			}
			
			if err := s.userAggregator.ValidateUser(ctx, donation.DonatorID, false); err != nil {
				return fmt.Errorf("donator validation failed: %v", err)
			}
		}
	} else {
		// No external validation available, just check required fields
		fmt.Printf("⚠️ No user validation service available, skipping user existence check\n")
		
		// For non-anonymous donations, donator ID is still required
		if !donation.IsAnonymous && donation.DonatorID == 0 {
			return errors.New("donator ID is required for non-anonymous donations")
		}
	}

	// Log the validation result
	fmt.Printf("✅ User validation passed - StreamerID: %d, DonatorID: %d, Anonymous: %t\n", 
		donation.StreamerID, donation.DonatorID, donation.IsAnonymous)

	return nil
}

func (s *donationService) GetByID(id uint) (*models.Donation, error) {
	donation, err := s.donationRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Populate user data for response
	if err := s.populateUserData(donation); err != nil {
		// Log the error but don't fail the get operation
		fmt.Printf("Warning: Failed to populate user data: %v\n", err)
	}

	return donation, nil
}

func (s *donationService) GetByTransactionID(transactionID string) (*models.Donation, error) {
	donation, err := s.donationRepo.GetByTransactionID(transactionID)
	if err != nil {
		return nil, err
	}

	// Populate user data for response
	if err := s.populateUserData(donation); err != nil {
		// Log the error but don't fail the get operation
		fmt.Printf("Warning: Failed to populate user data: %v\n", err)
	}

	return donation, nil
}

func (s *donationService) List(page, pageSize int) ([]*models.Donation, error) {
	offset := (page - 1) * pageSize
	donations, err := s.donationRepo.List(offset, pageSize)
	if err != nil {
		return nil, err
	}

	// Populate user data for each donation
	for _, donation := range donations {
		if err := s.populateUserData(donation); err != nil {
			fmt.Printf("Warning: Failed to populate user data for donation ID %d: %v\n", donation.ID, err)
		}
	}

	return donations, nil
}

func (s *donationService) GetByDonatorID(donatorID uint, page, pageSize int) ([]*models.Donation, error) {
	offset := (page - 1) * pageSize
	donations, err := s.donationRepo.GetByDonatorID(donatorID, offset, pageSize)
	if err != nil {
		return nil, err
	}

	// Populate user data for each donation
	for _, donation := range donations {
		if err := s.populateUserData(donation); err != nil {
			fmt.Printf("Warning: Failed to populate user data for donation ID %d: %v\n", donation.ID, err)
		}
	}

	return donations, nil
}

func (s *donationService) GetByStreamerID(streamerID uint, page, pageSize int) ([]*models.Donation, error) {
	offset := (page - 1) * pageSize
	donations, err := s.donationRepo.GetByStreamerID(streamerID, offset, pageSize)
	if err != nil {
		return nil, err
	}

	// Populate user data for each donation
	for _, donation := range donations {
		if err := s.populateUserData(donation); err != nil {
			fmt.Printf("Warning: Failed to populate user data for donation ID %d: %v\n", donation.ID, err)
		}
	}

	return donations, nil
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
	donations, err := s.donationRepo.GetLatestDonations(limit)
	if err != nil {
		return nil, err
	}

	// Populate user data for each donation
	for _, donation := range donations {
		if err := s.populateUserData(donation); err != nil {
			fmt.Printf("Warning: Failed to populate user data for donation ID %d: %v\n", donation.ID, err)
		}
	}

	return donations, nil
}

func (s *donationService) GetTotalAmountByStreamer(streamerID uint) (float64, error) {
	return s.donationRepo.GetTotalAmountByStreamer(streamerID)
} 