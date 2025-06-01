package repository

import (
	"github.com/rzfd/mediashar/internal/models"
	"gorm.io/gorm"
)

// PlatformRepository interface defines methods for platform operations
type PlatformRepository interface {
	// StreamingPlatform operations
	CreatePlatform(platform *models.StreamingPlatform) error
	GetPlatformByID(id uint) (*models.StreamingPlatform, error)
	GetPlatformsByUserID(userID uint) ([]*models.StreamingPlatform, error)
	GetPlatformByUserAndType(userID uint, platformType string) (*models.StreamingPlatform, error)
	GetActivePlatforms(offset, limit int) ([]*models.StreamingPlatform, error)
	
	// StreamingContent operations
	CreateContent(content *models.StreamingContent) error
	GetContentByID(id uint) (*models.StreamingContent, error)
	GetContentByPlatformID(platformID uint, offset, limit int) ([]*models.StreamingContent, error)
	GetContentByURL(contentURL string) (*models.StreamingContent, error)
	GetLiveContent(offset, limit int) ([]*models.StreamingContent, error)
	GetContentByType(contentType string, offset, limit int) ([]*models.StreamingContent, error)
	
	// ContentDonation operations
	CreateContentDonation(contentDonation *models.ContentDonation) error
	GetContentDonationByID(id uint) (*models.ContentDonation, error)
	GetContentDonationsByDonationID(donationID uint) ([]*models.ContentDonation, error)
	GetContentDonationsByContentID(contentID uint) ([]*models.ContentDonation, error)
	GetContentDonationsByPlatform(platformType string, offset, limit int) ([]*models.ContentDonation, error)
}

type platformRepository struct {
	db *gorm.DB
}

// NewPlatformRepository creates a new platform repository
func NewPlatformRepository(db *gorm.DB) PlatformRepository {
	return &platformRepository{db: db}
}

// StreamingPlatform operations

func (r *platformRepository) CreatePlatform(platform *models.StreamingPlatform) error {
	return r.db.Create(platform).Error
}

func (r *platformRepository) GetPlatformByID(id uint) (*models.StreamingPlatform, error) {
	var platform models.StreamingPlatform
	err := r.db.Preload("User").Preload("StreamingContent").First(&platform, id).Error
	if err != nil {
		return nil, err
	}
	return &platform, nil
}

func (r *platformRepository) GetPlatformsByUserID(userID uint) ([]*models.StreamingPlatform, error) {
	var platforms []*models.StreamingPlatform
	err := r.db.Where("user_id = ?", userID).
		Preload("StreamingContent").
		Find(&platforms).Error
	return platforms, err
}

func (r *platformRepository) GetPlatformByUserAndType(userID uint, platformType string) (*models.StreamingPlatform, error) {
	var platform models.StreamingPlatform
	err := r.db.Where("user_id = ? AND platform_type = ?", userID, platformType).
		Preload("User").
		Preload("StreamingContent").
		First(&platform).Error
	if err != nil {
		return nil, err
	}
	return &platform, nil
}

func (r *platformRepository) GetActivePlatforms(offset, limit int) ([]*models.StreamingPlatform, error) {
	var platforms []*models.StreamingPlatform
	err := r.db.Where("is_active = ?", true).
		Preload("User").
		Offset(offset).
		Limit(limit).
		Find(&platforms).Error
	return platforms, err
}

// StreamingContent operations

func (r *platformRepository) CreateContent(content *models.StreamingContent) error {
	return r.db.Create(content).Error
}

func (r *platformRepository) GetContentByID(id uint) (*models.StreamingContent, error) {
	var content models.StreamingContent
	err := r.db.Preload("Platform").
		Preload("Platform.User").
		Preload("ContentDonations").
		First(&content, id).Error
	if err != nil {
		return nil, err
	}
	return &content, nil
}

func (r *platformRepository) GetContentByPlatformID(platformID uint, offset, limit int) ([]*models.StreamingContent, error) {
	var contents []*models.StreamingContent
	err := r.db.Where("platform_id = ?", platformID).
		Preload("Platform").
		Preload("ContentDonations").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&contents).Error
	return contents, err
}

func (r *platformRepository) GetContentByURL(contentURL string) (*models.StreamingContent, error) {
	var content models.StreamingContent
	err := r.db.Where("content_url = ?", contentURL).
		Preload("Platform").
		Preload("Platform.User").
		Preload("ContentDonations").
		First(&content).Error
	if err != nil {
		return nil, err
	}
	return &content, nil
}

func (r *platformRepository) GetLiveContent(offset, limit int) ([]*models.StreamingContent, error) {
	var contents []*models.StreamingContent
	err := r.db.Where("is_live = ?", true).
		Preload("Platform").
		Preload("Platform.User").
		Offset(offset).
		Limit(limit).
		Order("started_at DESC").
		Find(&contents).Error
	return contents, err
}

func (r *platformRepository) GetContentByType(contentType string, offset, limit int) ([]*models.StreamingContent, error) {
	var contents []*models.StreamingContent
	err := r.db.Where("content_type = ?", contentType).
		Preload("Platform").
		Preload("Platform.User").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&contents).Error
	return contents, err
}

// ContentDonation operations

func (r *platformRepository) CreateContentDonation(contentDonation *models.ContentDonation) error {
	return r.db.Create(contentDonation).Error
}

func (r *platformRepository) GetContentDonationByID(id uint) (*models.ContentDonation, error) {
	var contentDonation models.ContentDonation
	err := r.db.Preload("Donation").
		Preload("Content").
		Preload("Content.Platform").
		First(&contentDonation, id).Error
	if err != nil {
		return nil, err
	}
	return &contentDonation, nil
}

func (r *platformRepository) GetContentDonationsByDonationID(donationID uint) ([]*models.ContentDonation, error) {
	var contentDonations []*models.ContentDonation
	err := r.db.Where("donation_id = ?", donationID).
		Preload("Content").
		Preload("Content.Platform").
		Find(&contentDonations).Error
	return contentDonations, err
}

func (r *platformRepository) GetContentDonationsByContentID(contentID uint) ([]*models.ContentDonation, error) {
	var contentDonations []*models.ContentDonation
	err := r.db.Where("content_id = ?", contentID).
		Preload("Donation").
		Find(&contentDonations).Error
	return contentDonations, err
}

func (r *platformRepository) GetContentDonationsByPlatform(platformType string, offset, limit int) ([]*models.ContentDonation, error) {
	var contentDonations []*models.ContentDonation
	err := r.db.Where("platform_type = ?", platformType).
		Preload("Donation").
		Preload("Content").
		Preload("Content.Platform").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&contentDonations).Error
	return contentDonations, err
} 