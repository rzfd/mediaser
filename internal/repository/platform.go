package repository

import "github.com/rzfd/mediashar/internal/models"

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