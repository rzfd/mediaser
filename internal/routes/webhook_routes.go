package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rzfd/mediashar/internal/handler"
)

// SetupWebhookRoutes configures webhook-related routes
func SetupWebhookRoutes(api *echo.Group, webhookHandler *handler.WebhookHandler, qrisHandler *handler.QRISHandler) {
	// Webhook routes (no authentication, but should be secured by webhook secrets)
	webhooks := api.Group("/webhooks")
	
	// Payment provider webhooks
	webhooks.POST("/paypal", webhookHandler.HandlePaypalWebhook)
	webhooks.POST("/stripe", webhookHandler.HandleStripeWebhook)
	webhooks.POST("/crypto", webhookHandler.HandleCryptoWebhook)
	
	// QRIS webhook
	webhooks.POST("/qris", qrisHandler.QRISCallback)
} 