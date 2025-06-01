package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/repository"
	"github.com/rzfd/mediashar/internal/service"
	"gorm.io/gorm"
)

type PlatformHandler struct {
	platformService *service.PlatformService
	platformRepo    repository.PlatformRepository
}

func NewPlatformHandler(platformService *service.PlatformService, platformRepo repository.PlatformRepository) *PlatformHandler {
	return &PlatformHandler{
		platformService: platformService,
		platformRepo:    platformRepo,
	}
}

// ValidateURL validates YouTube or TikTok URLs and extracts metadata
func (h *PlatformHandler) ValidateURL(c echo.Context) error {
	var req struct {
		URL string `json:"url" validate:"required,url"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// Temporarily disable validation for testing
	// if err := c.Validate(&req); err != nil {
	// 	return c.JSON(http.StatusBadRequest, map[string]interface{}{
	// 		"status":  "error",
	// 		"message": "URL is required and must be valid",
	// 	})
	// }

	// Basic URL check
	if req.URL == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "URL is required",
		})
	}

	result, err := h.platformService.ValidateURL(req.URL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Failed to validate URL",
		})
	}

	if !result.IsValid {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Invalid or unsupported URL. Supported platforms: YouTube, TikTok",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   result,
	})
}

// ConnectPlatform connects a social media platform to user account
func (h *PlatformHandler) ConnectPlatform(c echo.Context) error {
	userID := getUserIDFromContext(c)

	var req struct {
		PlatformType     string `json:"platform_type" validate:"required,oneof=youtube tiktok"`
		ChannelURL       string `json:"channel_url" validate:"required,url"`
		PlatformUsername string `json:"platform_username" validate:"required"`
		AutoSync         bool   `json:"auto_sync"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Validation failed",
		})
	}

	// Validate the URL first
	validation, err := h.platformService.ValidateURL(req.ChannelURL)
	if err != nil || !validation.IsValid {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Invalid platform URL",
		})
	}

	// Check if platform matches URL
	if validation.Platform != req.PlatformType {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Platform type doesn't match URL",
		})
	}

	// Check if platform already connected
	existingPlatform, err := h.platformRepo.GetPlatformByUserAndType(userID, req.PlatformType)
	if err == nil && existingPlatform != nil {
		return c.JSON(http.StatusConflict, map[string]interface{}{
			"status":  "error",
			"message": "Platform already connected for this user",
		})
	}

	// Create platform connection
	platform := &models.StreamingPlatform{
		UserID:           userID,
		PlatformType:     req.PlatformType,
		PlatformUsername: req.PlatformUsername,
		ChannelURL:       req.ChannelURL,
		IsActive:         true,
	}

	// Extract additional metadata from validation result
	if metadata := validation.Metadata; metadata != nil {
		if channelName, ok := metadata["creator"].(string); ok {
			platform.ChannelName = channelName
		}
		if profileImage, ok := metadata["thumbnail"].(string); ok {
			platform.ProfileImageURL = profileImage
		}
		if platformUserID, ok := metadata["channel_id"].(string); ok {
			platform.PlatformUserID = platformUserID
		} else if platformUserID, ok := metadata["video_id"].(string); ok {
			// For video URLs, extract channel info differently
			platform.PlatformUserID = platformUserID
		} else {
			// Fallback to username
			platform.PlatformUserID = req.PlatformUsername
		}
	}

	// Save to database
	err = h.platformRepo.CreatePlatform(platform)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Failed to connect platform",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"message": "Platform connected successfully",
		"data":    platform,
	})
}

// GetConnectedPlatforms returns list of connected platforms for user
func (h *PlatformHandler) GetConnectedPlatforms(c echo.Context) error {
	userID := getUserIDFromContext(c)

	platforms, err := h.platformRepo.GetPlatformsByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Failed to retrieve connected platforms",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   platforms,
	})
}

// GetPlatformByID returns a specific platform by ID
func (h *PlatformHandler) GetPlatformByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Invalid platform ID",
		})
	}

	platform, err := h.platformRepo.GetPlatformByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"status":  "error",
				"message": "Platform not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Failed to retrieve platform",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   platform,
	})
}

// CreateContent creates new streaming content
func (h *PlatformHandler) CreateContent(c echo.Context) error {
	var req struct {
		PlatformID   uint   `json:"platform_id" validate:"required"`
		ContentType  string `json:"content_type" validate:"required,oneof=live video short"`
		ContentID    string `json:"content_id" validate:"required"`
		ContentURL   string `json:"content_url" validate:"required,url"`
		Title        string `json:"title"`
		Description  string `json:"description"`
		ThumbnailURL string `json:"thumbnail_url"`
		Duration     *int   `json:"duration"`
		IsLive       bool   `json:"is_live"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Validation failed",
		})
	}

	// Verify platform exists
	_, err := h.platformRepo.GetPlatformByID(req.PlatformID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"status":  "error",
				"message": "Platform not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Failed to verify platform",
		})
	}

	// Create content
	content := &models.StreamingContent{
		PlatformID:   req.PlatformID,
		ContentType:  req.ContentType,
		ContentID:    req.ContentID,
		ContentURL:   req.ContentURL,
		Title:        req.Title,
		Description:  req.Description,
		ThumbnailURL: req.ThumbnailURL,
		Duration:     req.Duration,
		IsLive:       req.IsLive,
	}

	err = h.platformRepo.CreateContent(content)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Failed to create content",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"message": "Content created successfully",
		"data":    content,
	})
}

// GetContentByID returns specific content by ID
func (h *PlatformHandler) GetContentByID(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Invalid content ID",
		})
	}

	content, err := h.platformRepo.GetContentByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"status":  "error",
				"message": "Content not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Failed to retrieve content",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   content,
	})
}

// GetContentByURL returns content by URL
func (h *PlatformHandler) GetContentByURL(c echo.Context) error {
	var req struct {
		URL string `json:"url" validate:"required,url"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "URL is required and must be valid",
		})
	}

	content, err := h.platformRepo.GetContentByURL(req.URL)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"status":  "error",
				"message": "Content not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Failed to retrieve content",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   content,
	})
}

// GetLiveContent returns all live content
func (h *PlatformHandler) GetLiveContent(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	
	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	contents, err := h.platformRepo.GetLiveContent(offset, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Failed to retrieve live content",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   contents,
	})
}

// CreateContentDonation creates a donation linked to specific content
func (h *PlatformHandler) CreateContentDonation(c echo.Context) error {
	var req struct {
		DonationID   uint   `json:"donation_id" validate:"required"`
		ContentURL   string `json:"content_url" validate:"required,url"`
		PlatformType string `json:"platform_type" validate:"required,oneof=youtube tiktok"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Validation failed",
		})
	}

	// Try to find existing content by URL
	var contentID *uint
	content, err := h.platformRepo.GetContentByURL(req.ContentURL)
	if err == nil && content != nil {
		contentID = &content.ID
	}

	// Create content donation
	contentDonation := &models.ContentDonation{
		DonationID:   req.DonationID,
		ContentID:    contentID,
		PlatformType: req.PlatformType,
		ContentURL:   req.ContentURL,
	}

	err = h.platformRepo.CreateContentDonation(contentDonation)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Failed to create content donation",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status":  "success",
		"message": "Content donation created successfully",
		"data":    contentDonation,
	})
}

// GetContentDonationsByDonation returns content donations for a specific donation
func (h *PlatformHandler) GetContentDonationsByDonation(c echo.Context) error {
	idParam := c.Param("donationId")
	donationID, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Invalid donation ID",
		})
	}

	contentDonations, err := h.platformRepo.GetContentDonationsByDonationID(uint(donationID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Failed to retrieve content donations",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   contentDonations,
	})
}

// GetSupportedPlatforms returns list of supported platforms
func (h *PlatformHandler) GetSupportedPlatforms(c echo.Context) error {
	platforms := []map[string]interface{}{
		{
			"platform":    "youtube",
			"name":        "YouTube",
			"description": "Videos, Live Streams, Shorts, Channels",
			"url_formats": []string{
				"https://www.youtube.com/watch?v=VIDEO_ID",
				"https://youtu.be/VIDEO_ID",
				"https://www.youtube.com/live/VIDEO_ID",
				"https://www.youtube.com/shorts/VIDEO_ID",
				"https://www.youtube.com/@username",
			},
		},
		{
			"platform":    "tiktok",
			"name":        "TikTok",
			"description": "Videos, Live Streams, Profiles",
			"url_formats": []string{
				"https://www.tiktok.com/@username/video/VIDEO_ID",
				"https://vm.tiktok.com/SHORT_CODE",
				"https://www.tiktok.com/@username/live",
				"https://www.tiktok.com/@username",
			},
		},
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   platforms,
	})
}

// Helper function to get user ID from context
func getUserIDFromContext(c echo.Context) uint {
	// This should extract user ID from JWT token
	// For now, return a mock user ID
	// In production, implement proper JWT token parsing
	userInterface := c.Get("user")
	if userInterface != nil {
		if user, ok := userInterface.(map[string]interface{}); ok {
			if userID, ok := user["id"].(float64); ok {
				return uint(userID)
			}
		}
	}
	return 1 // Mock user ID for testing
} 