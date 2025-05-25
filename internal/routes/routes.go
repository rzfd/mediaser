package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rzfd/mediashar/internal/handler"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(e *echo.Echo, userHandler *handler.UserHandler, donationHandler *handler.DonationHandler, webhookHandler *handler.WebhookHandler, authHandler *handler.AuthHandler, qrisHandler *handler.QRISHandler, jwtSecret string) {
	// API group
	api := e.Group("/api")

	// Setup routes by handler
	SetupAuthRoutes(api, authHandler, jwtSecret)
	SetupUserRoutes(api, userHandler, jwtSecret)
	SetupDonationRoutes(api, donationHandler, jwtSecret)
	SetupQRISRoutes(api, qrisHandler, jwtSecret)
	SetupWebhookRoutes(api, webhookHandler, qrisHandler)
} 