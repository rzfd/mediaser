package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rzfd/mediashar/internal/handler"
	"github.com/rzfd/mediashar/internal/middleware"
)

// SetupPlatformRoutes sets up all platform integration routes
func SetupPlatformRoutes(api *echo.Group, platformHandler *handler.PlatformHandler, jwtSecret string) {
	// Public routes (no authentication required)
	api.POST("/content/validate", platformHandler.ValidateURL)
	api.GET("/platforms/supported", platformHandler.GetSupportedPlatforms)

	// Protected routes (authentication required)
	authGroup := api.Group("", middleware.JWTMiddleware(jwtSecret))

	// Platform Management
	authGroup.POST("/platforms/connect", platformHandler.ConnectPlatform)
	authGroup.GET("/platforms", platformHandler.GetConnectedPlatforms)
	authGroup.GET("/platforms/:id", platformHandler.GetPlatformByID)

	// Content Management
	authGroup.POST("/content", platformHandler.CreateContent)
	authGroup.GET("/content/:id", platformHandler.GetContentByID)
	authGroup.POST("/content/by-url", platformHandler.GetContentByURL)
	authGroup.GET("/content/live", platformHandler.GetLiveContent)

	// Content Donations
	authGroup.POST("/content-donations", platformHandler.CreateContentDonation)
	authGroup.GET("/content-donations/donation/:donationId", platformHandler.GetContentDonationsByDonation)

	// Alternative routes for backward compatibility
	authGroup.POST("/donations/to-content", platformHandler.CreateContentDonation)
}

// SetupPlatformRoutesWithCustomMiddleware allows custom middleware setup
func SetupPlatformRoutesWithCustomMiddleware(e *echo.Echo, platformHandler *handler.PlatformHandler, authMiddleware echo.MiddlewareFunc) {
	// Platform Integration Group
	platformGroup := e.Group("/api")

	// Public routes
	platformGroup.POST("/content/validate", platformHandler.ValidateURL)
	platformGroup.GET("/platforms/supported", platformHandler.GetSupportedPlatforms)

	// Protected routes with custom middleware
	authGroup := platformGroup.Group("", authMiddleware)

	// Platform Management
	authGroup.POST("/platforms/connect", platformHandler.ConnectPlatform)
	authGroup.GET("/platforms", platformHandler.GetConnectedPlatforms)
	authGroup.GET("/platforms/:id", platformHandler.GetPlatformByID)

	// Content Management
	authGroup.POST("/content", platformHandler.CreateContent)
	authGroup.GET("/content/:id", platformHandler.GetContentByID)
	authGroup.POST("/content/by-url", platformHandler.GetContentByURL)
	authGroup.GET("/content/live", platformHandler.GetLiveContent)

	// Content Donations
	authGroup.POST("/content-donations", platformHandler.CreateContentDonation)
	authGroup.GET("/content-donations/donation/:donationId", platformHandler.GetContentDonationsByDonation)

	// Alternative routes
	authGroup.POST("/donations/to-content", platformHandler.CreateContentDonation)
} 