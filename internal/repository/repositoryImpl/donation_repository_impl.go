package repositoryImpl

import (
	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/repository"
	"gorm.io/gorm"
)

type donationRepository struct {
	db *gorm.DB
}

func NewDonationRepository(db *gorm.DB) repository.DonationRepository {
	return &donationRepository{db: db}
}

func (r *donationRepository) Create(donation *models.Donation) error {
	return r.db.Create(donation).Error
}

func (r *donationRepository) GetByID(id uint) (*models.Donation, error) {
	var donation models.Donation
	// Remove Preload for Donator and Streamer since they are excluded from GORM with gorm:"-"
	err := r.db.First(&donation, id).Error
	if err != nil {
		return nil, err
	}
	return &donation, nil
}

func (r *donationRepository) GetByTransactionID(transactionID string) (*models.Donation, error) {
	var donation models.Donation
	err := r.db.Where("transaction_id = ?", transactionID).First(&donation).Error
	if err != nil {
		return nil, err
	}
	return &donation, nil
}

func (r *donationRepository) Update(donation *models.Donation) error {
	return r.db.Save(donation).Error
}

func (r *donationRepository) Delete(id uint) error {
	return r.db.Delete(&models.Donation{}, id).Error
}

func (r *donationRepository) List(offset, limit int) ([]*models.Donation, error) {
	var donations []*models.Donation
	err := r.db.Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&donations).Error
	return donations, err
}

func (r *donationRepository) GetByDonatorID(donatorID uint, offset, limit int) ([]*models.Donation, error) {
	var donations []*models.Donation
	err := r.db.Where("donator_id = ?", donatorID).
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&donations).Error
	return donations, err
}

func (r *donationRepository) GetByStreamerID(streamerID uint, offset, limit int) ([]*models.Donation, error) {
	var donations []*models.Donation
	err := r.db.Where("streamer_id = ?", streamerID).
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&donations).Error
	return donations, err
}

func (r *donationRepository) UpdateStatus(id uint, status models.PaymentStatus) error {
	return r.db.Model(&models.Donation{}).Where("id = ?", id).Update("status", status).Error
}

func (r *donationRepository) GetLatestDonations(limit int) ([]*models.Donation, error) {
	var donations []*models.Donation
	err := r.db.Where("status = ?", models.PaymentCompleted).
		Order("created_at DESC").
		Limit(limit).
		Find(&donations).Error
	return donations, err
}

func (r *donationRepository) GetTotalAmountByStreamer(streamerID uint) (float64, error) {
	var total float64
	err := r.db.Model(&models.Donation{}).
		Where("streamer_id = ? AND status = ?", streamerID, models.PaymentCompleted).
		Select("COALESCE(SUM(amount), 0) as total").
		Scan(&total).Error
	return total, err
} 