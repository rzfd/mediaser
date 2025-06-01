package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rzfd/mediashar/internal/handler"
	"github.com/rzfd/mediashar/internal/middleware"
)

// SetupMidtransRoutes sets up all Midtrans related routes
func SetupMidtransRoutes(api *echo.Group, midtransHandler *handler.MidtransHandler, jwtSecret string) {
	midtrans := api.Group("/midtrans")

	// Public routes (no authentication required)
	midtrans.POST("/webhook", midtransHandler.HandleWebhook) // Webhook from Midtrans
	midtrans.GET("/status/:orderId", midtransHandler.GetTransactionStatus)

	// Protected routes (authentication required)
	protected := midtrans.Group("", middleware.JWTMiddleware(jwtSecret))
	protected.POST("/payment/:donationId", midtransHandler.CreatePayment)
} 