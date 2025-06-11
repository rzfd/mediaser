package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rzfd/mediashar/internal/handler"
	"github.com/rzfd/mediashar/internal/middleware"
)

// SetupAuthRoutes configures authentication-related routes
func SetupAuthRoutes(api *echo.Group, authHandler *handler.AuthHandler, jwtSecret string) {
	// Public auth routes (no authentication required)
	auth := api.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/google", authHandler.GoogleLogin)
	auth.POST("/refresh", authHandler.RefreshToken)

	// Protected auth routes (authentication required)
	protectedAuth := api.Group("/auth", middleware.JWTMiddleware(jwtSecret))
	protectedAuth.GET("/profile", authHandler.GetProfile)
	protectedAuth.PUT("/profile", authHandler.UpdateProfile)
	protectedAuth.POST("/change-password", authHandler.ChangePassword)
	protectedAuth.POST("/logout", authHandler.Logout)
} 