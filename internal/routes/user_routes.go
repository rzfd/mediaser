package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rzfd/mediashar/internal/handler"
	"github.com/rzfd/mediashar/internal/middleware"
)

// SetupUserRoutes configures user-related routes
func SetupUserRoutes(api *echo.Group, userHandler *handler.UserHandler, jwtSecret string) {
	// Public user routes (no authentication required)
	api.GET("/users/:id", userHandler.GetUser)
	api.GET("/streamers", userHandler.ListStreamers)

	// Protected user routes (authentication required)
	protectedUsers := api.Group("/users", middleware.JWTMiddleware(jwtSecret))
	protectedUsers.POST("", userHandler.CreateUser)
	protectedUsers.PUT("/:id", userHandler.UpdateUser)
	protectedUsers.GET("/:id/donations", userHandler.GetUserDonations)
} 