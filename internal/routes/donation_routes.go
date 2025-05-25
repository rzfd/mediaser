package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rzfd/mediashar/internal/handler"
	"github.com/rzfd/mediashar/internal/middleware"
)

// SetupDonationRoutes configures donation-related routes
func SetupDonationRoutes(api *echo.Group, donationHandler *handler.DonationHandler, jwtSecret string) {
	// Protected donation routes (authentication required)
	protectedDonations := api.Group("/donations", middleware.JWTMiddleware(jwtSecret))
	protectedDonations.POST("", donationHandler.CreateDonation)
	protectedDonations.GET("", donationHandler.ListDonations)
	protectedDonations.GET("/:id", donationHandler.GetDonation)
	protectedDonations.GET("/latest", donationHandler.GetLatestDonations)

	// Streamer-only routes (authentication + streamer role required)
	streamerDonations := api.Group("/streamers", middleware.JWTMiddleware(jwtSecret), middleware.StreamerOnlyMiddleware())
	streamerDonations.GET("/:id/donations", donationHandler.GetStreamerDonations)
	streamerDonations.GET("/:id/total", donationHandler.GetTotalDonations)

	// Payment processing routes (authentication required)
	protectedPayments := api.Group("/payments", middleware.JWTMiddleware(jwtSecret))
	protectedPayments.POST("/process", donationHandler.ProcessPayment)
} 