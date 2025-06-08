package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rzfd/mediashar/internal/handler"
	"github.com/rzfd/mediashar/internal/middleware"
)

// SetupMediaShareRoutes configures all media share routes
func SetupMediaShareRoutes(api *echo.Group, mediaShareHandler *handler.MediaShareHandler, jwtSecret string) {
	// Create JWT middleware
	jwtAuth := middleware.JWTMiddleware(jwtSecret)
	
	// Streamer settings routes (requires authentication)
	streamers := api.Group("/streamers", jwtAuth)
	streamers.GET("/:streamer_id/media-settings", mediaShareHandler.GetStreamerSettings)
	streamers.PUT("/:streamer_id/media-settings", mediaShareHandler.UpdateStreamerSettings)
	streamers.GET("/:streamer_id/media-queue", mediaShareHandler.GetMediaQueue)
	streamers.GET("/:streamer_id/media-stats", mediaShareHandler.GetMediaStats)

	// Donation-related media share routes (requires authentication)
	donations := api.Group("/donations", jwtAuth)
	donations.POST("/:donation_id/media-share", mediaShareHandler.SubmitMediaShare)

	// Media management routes (requires authentication)
	media := api.Group("/media", jwtAuth)
	media.PUT("/:media_id/approve", mediaShareHandler.ApproveMedia)
	media.PUT("/:media_id/reject", mediaShareHandler.RejectMedia)
} 