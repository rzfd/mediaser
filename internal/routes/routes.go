package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rzfd/mediashar/internal/handler"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(e *echo.Echo, userHandler *handler.UserHandler, donationHandler *handler.DonationHandler, webhookHandler *handler.WebhookHandler, authHandler *handler.AuthHandler, qrisHandler *handler.QRISHandler, platformHandler *handler.PlatformHandler, midtransHandler *handler.MidtransHandler, currencyHandler *handler.CurrencyHandler, languageHandler *handler.LanguageHandler, mediaShareHandler *handler.MediaShareHandler, jwtSecret string) {
	// Health check routes (no prefix)
	healthHandler := handler.NewHealthHandler()
	e.GET("/health", healthHandler.HealthCheck)
	e.GET("/ready", healthHandler.ReadinessCheck)
	
	// API group
	api := e.Group("/api")
	
	// Health check in API group as well
	api.GET("/health", healthHandler.HealthCheck)
	api.GET("/ready", healthHandler.ReadinessCheck)

	// Setup routes by handler
	SetupAuthRoutes(api, authHandler, jwtSecret)
	SetupUserRoutes(api, userHandler, jwtSecret)
	SetupDonationRoutes(api, donationHandler, jwtSecret)
	SetupQRISRoutes(api, qrisHandler, jwtSecret)
	SetupMidtransRoutes(api, midtransHandler, jwtSecret)
	SetupWebhookRoutes(api, webhookHandler, qrisHandler)
	SetupPlatformRoutes(api, platformHandler, jwtSecret)
	SetupCurrencyRoutes(api, currencyHandler, jwtSecret)
	SetupLanguageRoutes(api, languageHandler, jwtSecret)
	SetupMediaShareRoutes(api, mediaShareHandler, jwtSecret)
} 