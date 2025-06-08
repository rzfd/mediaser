package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/service/serviceImpl"
	"github.com/rzfd/mediashar/pkg/logger"
)

type MediaShareHandler struct {
	service serviceImpl.MediaShareService
}

func NewMediaShareHandler(service serviceImpl.MediaShareService) *MediaShareHandler {
	return &MediaShareHandler{service: service}
}

// GetStreamerSettings godoc
// @Summary Get streamer media share settings
// @Tags MediaShare
// @Accept json
// @Produce json
// @Param streamer_id path int true "Streamer ID"
// @Success 200 {object} models.MediaShareSettings
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/streamers/{streamer_id}/media-settings [get]
func (h *MediaShareHandler) GetStreamerSettings(c echo.Context) error {
	appLogger := logger.GetLogger()
	
	streamerID, err := strconv.ParseUint(c.Param("streamer_id"), 10, 32)
	if err != nil {
		appLogger.Error(err, "Invalid streamer ID", "streamer_id", c.Param("streamer_id"))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid streamer ID"})
	}

	settings, err := h.service.GetSettingsByStreamerID(uint(streamerID))
	if err != nil {
		appLogger.Error(err, "Failed to get streamer settings", "streamer_id", streamerID)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get settings"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    settings,
	})
}

// UpdateStreamerSettings godoc
// @Summary Update streamer media share settings
// @Tags MediaShare
// @Accept json
// @Produce json
// @Param streamer_id path int true "Streamer ID"
// @Param settings body models.MediaShareSettings true "Settings"
// @Success 200 {object} models.MediaShareSettings
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/streamers/{streamer_id}/media-settings [put]
func (h *MediaShareHandler) UpdateStreamerSettings(c echo.Context) error {
	appLogger := logger.GetLogger()
	
	streamerID, err := strconv.ParseUint(c.Param("streamer_id"), 10, 32)
	if err != nil {
		appLogger.Error(err, "Invalid streamer ID", "streamer_id", c.Param("streamer_id"))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid streamer ID"})
	}

	var settings models.MediaShareSettings
	if err := c.Bind(&settings); err != nil {
		appLogger.Error(err, "Invalid request body")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Ensure the streamer ID matches
	settings.StreamerID = uint(streamerID)

	if err := h.service.UpdateSettings(&settings); err != nil {
		appLogger.Error(err, "Failed to update settings", "streamer_id", streamerID)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Settings updated successfully",
		"data":    settings,
	})
}

// SubmitMediaShare godoc
// @Summary Submit media share with donation
// @Tags MediaShare
// @Accept json
// @Produce json
// @Param donation_id path int true "Donation ID"
// @Param media_request body models.MediaShareRequest true "Media Share Request"
// @Success 201 {object} models.MediaShareResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/donations/{donation_id}/media-share [post]
func (h *MediaShareHandler) SubmitMediaShare(c echo.Context) error {
	appLogger := logger.GetLogger()
	
	donationID, err := strconv.ParseUint(c.Param("donation_id"), 10, 32)
	if err != nil {
		appLogger.Error(err, "Invalid donation ID", "donation_id", c.Param("donation_id"))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid donation ID"})
	}

	var req models.MediaShareRequest
	if err := c.Bind(&req); err != nil {
		appLogger.Error(err, "Invalid request body")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Get user info from context (set by auth middleware)
	userID := getMediaShareUserID(c)
	if userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
	}

	// For now, we'll use dummy streamer ID and donator name
	// In real implementation, you'd get this from the donation record
	streamerID := uint(1) // TODO: Get from donation record
	donatorName := "Anonymous" // TODO: Get from user record

	response, err := h.service.SubmitMediaShare(uint(donationID), streamerID, userID, &req, donatorName)
	if err != nil {
		appLogger.Error(err, "Failed to submit media share", "donation_id", donationID, "user_id", userID)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Media share submitted successfully",
		"data":    response,
	})
}

// GetMediaQueue godoc
// @Summary Get media queue for streamer
// @Tags MediaShare
// @Accept json
// @Produce json
// @Param streamer_id path int true "Streamer ID"
// @Param status query string false "Status filter (pending, approved, rejected, all)" default(all)
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/streamers/{streamer_id}/media-queue [get]
func (h *MediaShareHandler) GetMediaQueue(c echo.Context) error {
	appLogger := logger.GetLogger()
	
	streamerID, err := strconv.ParseUint(c.Param("streamer_id"), 10, 32)
	if err != nil {
		appLogger.Error(err, "Invalid streamer ID", "streamer_id", c.Param("streamer_id"))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid streamer ID"})
	}

	status := c.QueryParam("status")
	if status == "" {
		status = "all"
	}

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.QueryParam("page_size"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	items, total, err := h.service.GetMediaQueue(uint(streamerID), status, page, pageSize)
	if err != nil {
		appLogger.Error(err, "Failed to get media queue", "streamer_id", streamerID)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get media queue"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"items":      items,
			"total":      total,
			"page":       page,
			"page_size":  pageSize,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetMediaStats godoc
// @Summary Get media share statistics for streamer
// @Tags MediaShare
// @Accept json
// @Produce json
// @Param streamer_id path int true "Streamer ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/streamers/{streamer_id}/media-stats [get]
func (h *MediaShareHandler) GetMediaStats(c echo.Context) error {
	appLogger := logger.GetLogger()
	
	streamerID, err := strconv.ParseUint(c.Param("streamer_id"), 10, 32)
	if err != nil {
		appLogger.Error(err, "Invalid streamer ID", "streamer_id", c.Param("streamer_id"))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid streamer ID"})
	}

	stats, err := h.service.GetMediaStats(uint(streamerID))
	if err != nil {
		appLogger.Error(err, "Failed to get media stats", "streamer_id", streamerID)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get statistics"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"data":    stats,
	})
}

// ApproveMedia godoc
// @Summary Approve media share
// @Tags MediaShare
// @Accept json
// @Produce json
// @Param media_id path int true "Media ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/media/{media_id}/approve [put]
func (h *MediaShareHandler) ApproveMedia(c echo.Context) error {
	appLogger := logger.GetLogger()
	
	mediaID, err := strconv.ParseUint(c.Param("media_id"), 10, 32)
	if err != nil {
		appLogger.Error(err, "Invalid media ID", "media_id", c.Param("media_id"))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid media ID"})
	}

	userID := getMediaShareUserID(c)
	if userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
	}

	if err := h.service.ApproveMedia(userID, uint(mediaID)); err != nil {
		appLogger.Error(err, "Failed to approve media", "media_id", mediaID, "user_id", userID)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"success": "true",
		"message": "Media approved successfully",
	})
}

// RejectMedia godoc
// @Summary Reject media share
// @Tags MediaShare
// @Accept json
// @Produce json
// @Param media_id path int true "Media ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/media/{media_id}/reject [put]
func (h *MediaShareHandler) RejectMedia(c echo.Context) error {
	appLogger := logger.GetLogger()
	
	mediaID, err := strconv.ParseUint(c.Param("media_id"), 10, 32)
	if err != nil {
		appLogger.Error(err, "Invalid media ID", "media_id", c.Param("media_id"))
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid media ID"})
	}

	userID := getMediaShareUserID(c)
	if userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
	}

	if err := h.service.RejectMedia(userID, uint(mediaID)); err != nil {
		appLogger.Error(err, "Failed to reject media", "media_id", mediaID, "user_id", userID)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"success": "true",
		"message": "Media rejected successfully",
	})
}

// Helper function to get user ID from context for media share operations
func getMediaShareUserID(c echo.Context) uint {
	userIDInterface := c.Get("user_id")
	if userIDInterface == nil {
		return 0
	}
	
	switch userID := userIDInterface.(type) {
	case uint:
		return userID
	case float64:
		return uint(userID)
	case string:
		if id, err := strconv.ParseUint(userID, 10, 32); err == nil {
			return uint(id)
		}
	}
	
	return 0
}

 