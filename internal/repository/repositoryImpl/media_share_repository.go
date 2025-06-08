package repositoryImpl

import (
	"time"

	"gorm.io/gorm"
	"github.com/rzfd/mediashar/internal/models"
)

type MediaShareRepository interface {
	// Settings
	GetSettingsByStreamerID(streamerID uint) (*models.MediaShareSettings, error)
	CreateOrUpdateSettings(settings *models.MediaShareSettings) error
	
	// Media Share
	Create(mediaShare *models.MediaShare) error
	GetByID(id uint) (*models.MediaShare, error)
	GetByDonationID(donationID uint) (*models.MediaShare, error)
	GetQueueByStreamerID(streamerID uint, status string, limit, offset int) ([]*models.MediaQueueItem, error)
	GetTotalQueueCount(streamerID uint, status string) (int64, error)
	UpdateStatus(id uint, status models.MediaShareStatus) error
	GetStatsByStreamerID(streamerID uint) (map[string]int64, error)
}

type mediaShareRepository struct {
	db *gorm.DB
}

func NewMediaShareRepository(db *gorm.DB) MediaShareRepository {
	return &mediaShareRepository{db: db}
}

// Settings methods
func (r *mediaShareRepository) GetSettingsByStreamerID(streamerID uint) (*models.MediaShareSettings, error) {
	var settings models.MediaShareSettings
	err := r.db.Where("streamer_id = ?", streamerID).First(&settings).Error
	if err == gorm.ErrRecordNotFound {
		// Return default settings if not found
		return &models.MediaShareSettings{
			StreamerID:         streamerID,
			MediaShareEnabled:  true,
			MinDonationAmount:  5000,
			Currency:           "IDR",
			AllowYoutube:       true,
			AllowTiktok:        true,
			AutoApprove:        false,
			MaxDurationYoutube: 300,
			MaxDurationTiktok:  180,
			WelcomeMessage:     "Terima kasih atas donasi Anda! Silakan bagikan media favorit Anda.",
		}, nil
	}
	return &settings, err
}

func (r *mediaShareRepository) CreateOrUpdateSettings(settings *models.MediaShareSettings) error {
	return r.db.Save(settings).Error
}

// Media Share methods
func (r *mediaShareRepository) Create(mediaShare *models.MediaShare) error {
	return r.db.Create(mediaShare).Error
}

func (r *mediaShareRepository) GetByID(id uint) (*models.MediaShare, error) {
	var mediaShare models.MediaShare
	err := r.db.Preload("Donation").Preload("Streamer").Preload("Donator").
		Where("id = ?", id).First(&mediaShare).Error
	return &mediaShare, err
}

func (r *mediaShareRepository) GetByDonationID(donationID uint) (*models.MediaShare, error) {
	var mediaShare models.MediaShare
	err := r.db.Where("donation_id = ?", donationID).First(&mediaShare).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &mediaShare, err
}

func (r *mediaShareRepository) GetQueueByStreamerID(streamerID uint, status string, limit, offset int) ([]*models.MediaQueueItem, error) {
	var results []*models.MediaQueueItem
	
	query := r.db.Table("media_shares").
		Select(`
			media_shares.id,
			media_shares.type,
			media_shares.url,
			media_shares.title,
			media_shares.message,
			media_shares.status,
			media_shares.donator_name,
			media_shares.donation_amount,
			media_shares.currency,
			media_shares.thumbnail,
			media_shares.created_at as submitted_at,
			media_shares.processed_at
		`).
		Where("media_shares.streamer_id = ?", streamerID).
		Order("media_shares.created_at DESC")
	
	if status != "all" && status != "" {
		query = query.Where("media_shares.status = ?", status)
	}
	
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	
	err := query.Scan(&results).Error
	return results, err
}

func (r *mediaShareRepository) GetTotalQueueCount(streamerID uint, status string) (int64, error) {
	var count int64
	query := r.db.Model(&models.MediaShare{}).Where("streamer_id = ?", streamerID)
	
	if status != "all" && status != "" {
		query = query.Where("status = ?", status)
	}
	
	err := query.Count(&count).Error
	return count, err
}

func (r *mediaShareRepository) UpdateStatus(id uint, status models.MediaShareStatus) error {
	now := time.Now()
	return r.db.Model(&models.MediaShare{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       status,
			"processed_at": &now,
		}).Error
}

func (r *mediaShareRepository) GetStatsByStreamerID(streamerID uint) (map[string]int64, error) {
	stats := make(map[string]int64)
	
	// Get counts by status
	type StatusCount struct {
		Status string
		Count  int64
	}
	
	var statusCounts []StatusCount
	err := r.db.Model(&models.MediaShare{}).
		Select("status, COUNT(*) as count").
		Where("streamer_id = ?", streamerID).
		Group("status").
		Scan(&statusCounts).Error
	
	if err != nil {
		return nil, err
	}
	
	// Initialize all status counts to 0
	stats["pending"] = 0
	stats["approved"] = 0
	stats["rejected"] = 0
	stats["total"] = 0
	
	// Fill actual counts
	for _, sc := range statusCounts {
		stats[sc.Status] = sc.Count
		stats["total"] += sc.Count
	}
	
	// Get total donation amount from media shares
	var totalAmount float64
	err = r.db.Model(&models.MediaShare{}).
		Select("COALESCE(SUM(donation_amount), 0)").
		Where("streamer_id = ?", streamerID).
		Scan(&totalAmount).Error
	
	if err != nil {
		return stats, err
	}
	
	stats["total_amount"] = int64(totalAmount)
	
	return stats, nil
} 