package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rzfd/mediashar/internal/handler"
	"github.com/rzfd/mediashar/internal/middleware"
)

// SetupQRISRoutes configures QRIS-related routes
func SetupQRISRoutes(api *echo.Group, qrisHandler *handler.QRISHandler, jwtSecret string) {
	// Public QRIS routes (optional authentication for anonymous donations)
	qrisPublic := api.Group("", middleware.OptionalJWTMiddleware(jwtSecret))
	qrisPublic.POST("/qris/donate", qrisHandler.CreateQRISDonation)

	// Protected QRIS routes (authentication required)
	protectedQRIS := api.Group("/qris", middleware.JWTMiddleware(jwtSecret))
	protectedQRIS.POST("/donations/:id/generate", qrisHandler.GenerateQRIS)
	protectedQRIS.GET("/status/:transaction_id", qrisHandler.CheckQRISStatus)
} 