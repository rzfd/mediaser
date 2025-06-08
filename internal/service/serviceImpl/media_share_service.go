package serviceImpl

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/repository/repositoryImpl"
)

type MediaShareService interface {
	// Settings
	GetSettingsByStreamerID(streamerID uint) (*models.MediaShareSettings, error)
	UpdateSettings(settings *models.MediaShareSettings) error
	
	// Media Share
	SubmitMediaShare(donationID, streamerID, donatorID uint, req *models.MediaShareRequest, donatorName string) (*models.MediaShareResponse, error)
	GetMediaQueue(streamerID uint, status string, page, pageSize int) ([]*models.MediaQueueItem, int64, error)
	ApproveMedia(streamerID, mediaID uint) error
	RejectMedia(streamerID, mediaID uint) error
	GetMediaStats(streamerID uint) (map[string]int64, error)
	ValidateMediaShare(streamerID uint, donationAmount float64, mediaType models.MediaShareType) error
}

type mediaShareService struct {
	repo repositoryImpl.MediaShareRepository
}

func NewMediaShareService(repo repositoryImpl.MediaShareRepository) MediaShareService {
	return &mediaShareService{repo: repo}
}

// Settings methods
func (s *mediaShareService) GetSettingsByStreamerID(streamerID uint) (*models.MediaShareSettings, error) {
	return s.repo.GetSettingsByStreamerID(streamerID)
}

func (s *mediaShareService) UpdateSettings(settings *models.MediaShareSettings) error {
	if settings.StreamerID == 0 {
		return errors.New("streamer ID is required")
	}
	
	if settings.MinDonationAmount < 0 {
		return errors.New("minimum donation amount cannot be negative")
	}
	
	return s.repo.CreateOrUpdateSettings(settings)
}

// Media Share methods
func (s *mediaShareService) SubmitMediaShare(donationID, streamerID, donatorID uint, req *models.MediaShareRequest, donatorName string) (*models.MediaShareResponse, error) {
	// Validate URL format
	if err := s.validateURL(req.URL, req.Type); err != nil {
		return nil, err
	}
	
	// Check if media share already exists for this donation
	existing, err := s.repo.GetByDonationID(donationID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("media share already submitted for this donation")
	}
	
	// Get streamer settings for validation
	settings, err := s.repo.GetSettingsByStreamerID(streamerID)
	if err != nil {
		return nil, err
	}
	
	// Validate media share eligibility
	if err := s.validateMediaShareWithSettings(settings, req.DonationAmount, req.Type); err != nil {
		return nil, err
	}
	
	// Generate thumbnail and extract metadata
	thumbnail, duration := s.extractMediaMetadata(req.URL, req.Type)
	
	// Create media share
	mediaShare := &models.MediaShare{
		DonationID:     donationID,
		StreamerID:     streamerID,
		DonatorID:      donatorID,
		Type:           req.Type,
		URL:            req.URL,
		Title:          req.Title,
		Message:        req.Message,
		Status:         models.MediaShareStatusPending,
		DonationAmount: req.DonationAmount,
		Currency:       "IDR", // Default to IDR for now
		DonatorName:    donatorName,
		Thumbnail:      thumbnail,
		Duration:       duration,
	}
	
	// Auto-approve if enabled
	if settings.AutoApprove {
		mediaShare.Status = models.MediaShareStatusApproved
	}
	
	if err := s.repo.Create(mediaShare); err != nil {
		return nil, err
	}
	
	return &models.MediaShareResponse{
		ID:        mediaShare.ID,
		Type:      mediaShare.Type,
		URL:       mediaShare.URL,
		Title:     mediaShare.Title,
		Message:   mediaShare.Message,
		Status:    mediaShare.Status,
		CreatedAt: mediaShare.CreatedAt,
	}, nil
}

func (s *mediaShareService) GetMediaQueue(streamerID uint, status string, page, pageSize int) ([]*models.MediaQueueItem, int64, error) {
	offset := (page - 1) * pageSize
	
	items, err := s.repo.GetQueueByStreamerID(streamerID, status, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	
	total, err := s.repo.GetTotalQueueCount(streamerID, status)
	if err != nil {
		return nil, 0, err
	}
	
	return items, total, nil
}

func (s *mediaShareService) ApproveMedia(streamerID, mediaID uint) error {
	// Verify the media belongs to the streamer
	media, err := s.repo.GetByID(mediaID)
	if err != nil {
		return err
	}
	
	if media.StreamerID != streamerID {
		return errors.New("unauthorized: media does not belong to this streamer")
	}
	
	return s.repo.UpdateStatus(mediaID, models.MediaShareStatusApproved)
}

func (s *mediaShareService) RejectMedia(streamerID, mediaID uint) error {
	// Verify the media belongs to the streamer
	media, err := s.repo.GetByID(mediaID)
	if err != nil {
		return err
	}
	
	if media.StreamerID != streamerID {
		return errors.New("unauthorized: media does not belong to this streamer")
	}
	
	return s.repo.UpdateStatus(mediaID, models.MediaShareStatusRejected)
}

func (s *mediaShareService) GetMediaStats(streamerID uint) (map[string]int64, error) {
	return s.repo.GetStatsByStreamerID(streamerID)
}

func (s *mediaShareService) ValidateMediaShare(streamerID uint, donationAmount float64, mediaType models.MediaShareType) error {
	settings, err := s.repo.GetSettingsByStreamerID(streamerID)
	if err != nil {
		return err
	}
	
	return s.validateMediaShareWithSettings(settings, donationAmount, mediaType)
}

// Helper methods
func (s *mediaShareService) validateURL(mediaURL string, mediaType models.MediaShareType) error {
	parsedURL, err := url.Parse(mediaURL)
	if err != nil {
		return errors.New("invalid URL format")
	}
	
	switch mediaType {
	case models.MediaShareTypeYoutube:
		if !s.isValidYouTubeURL(parsedURL) {
			return errors.New("invalid YouTube URL format")
		}
	case models.MediaShareTypeTiktok:
		if !s.isValidTikTokURL(parsedURL) {
			return errors.New("invalid TikTok URL format")
		}
	default:
		return errors.New("unsupported media type")
	}
	
	return nil
}

func (s *mediaShareService) isValidYouTubeURL(u *url.URL) bool {
	// Valid YouTube domains
	validDomains := []string{"youtube.com", "www.youtube.com", "youtu.be", "m.youtube.com"}
	
	domainValid := false
	for _, domain := range validDomains {
		if u.Host == domain {
			domainValid = true
			break
		}
	}
	
	if !domainValid {
		return false
	}
	
	// Check for video ID
	if u.Host == "youtu.be" {
		// Short format: youtu.be/VIDEO_ID
		return len(strings.TrimPrefix(u.Path, "/")) > 0
	} else {
		// Long format: youtube.com/watch?v=VIDEO_ID
		return u.Query().Get("v") != ""
	}
}

func (s *mediaShareService) isValidTikTokURL(u *url.URL) bool {
	// Valid TikTok domains
	validDomains := []string{"tiktok.com", "www.tiktok.com", "m.tiktok.com"}
	
	domainValid := false
	for _, domain := range validDomains {
		if u.Host == domain {
			domainValid = true
			break
		}
	}
	
	if !domainValid {
		return false
	}
	
	// Check path format: /@username/video/VIDEO_ID
	pathRegex := regexp.MustCompile(`^/@[^/]+/video/\d+`)
	return pathRegex.MatchString(u.Path)
}

func (s *mediaShareService) validateMediaShareWithSettings(settings *models.MediaShareSettings, donationAmount float64, mediaType models.MediaShareType) error {
	if !settings.MediaShareEnabled {
		return errors.New("media share is disabled for this streamer")
	}
	
	if donationAmount < settings.MinDonationAmount {
		return fmt.Errorf("donation amount (%.2f) is below minimum required (%.2f)", 
			donationAmount, settings.MinDonationAmount)
	}
	
	switch mediaType {
	case models.MediaShareTypeYoutube:
		if !settings.AllowYoutube {
			return errors.New("YouTube media sharing is disabled for this streamer")
		}
	case models.MediaShareTypeTiktok:
		if !settings.AllowTiktok {
			return errors.New("TikTok media sharing is disabled for this streamer")
		}
	default:
		return errors.New("unsupported media type")
	}
	
	return nil
}

func (s *mediaShareService) extractMediaMetadata(mediaURL string, mediaType models.MediaShareType) (thumbnail string, duration int) {
	// For now, return basic placeholder data
	// In a real implementation, you would call the respective APIs
	
	switch mediaType {
	case models.MediaShareTypeYoutube:
		// Extract YouTube video ID and generate thumbnail
		if videoID := s.extractYouTubeVideoID(mediaURL); videoID != "" {
			thumbnail = fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", videoID)
		} else {
			thumbnail = "https://via.placeholder.com/480x360/ff0000/ffffff?text=YouTube"
		}
		duration = 180 // Default 3 minutes
		
	case models.MediaShareTypeTiktok:
		thumbnail = "https://via.placeholder.com/300x400/ff0050/ffffff?text=TikTok"
		duration = 60 // Default 1 minute
	}
	
	return thumbnail, duration
}

func (s *mediaShareService) extractYouTubeVideoID(mediaURL string) string {
	parsedURL, err := url.Parse(mediaURL)
	if err != nil {
		return ""
	}
	
	if parsedURL.Host == "youtu.be" {
		return strings.TrimPrefix(parsedURL.Path, "/")
	}
	
	return parsedURL.Query().Get("v")
} 