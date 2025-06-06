package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rzfd/mediashar/internal/handler"
	"github.com/rzfd/mediashar/internal/middleware"
)

// SetupLanguageRoutes configures language-related routes
func SetupLanguageRoutes(api *echo.Group, languageHandler *handler.LanguageHandler, jwtSecret string) {
	// Language API routes
	languageGroup := api.Group("/language")
	{
		// Public routes
		languageGroup.POST("/translate", languageHandler.TranslateText)
		languageGroup.GET("/list", languageHandler.GetSupportedLanguages)
		languageGroup.GET("/translation", languageHandler.GetTranslation)
		languageGroup.POST("/translation", languageHandler.AddTranslation)
		languageGroup.POST("/detect", languageHandler.DetectLanguage)
		languageGroup.POST("/translate/bulk", languageHandler.BulkTranslate)
		
		// Protected routes (authentication required)
		protectedLanguage := languageGroup.Group("", middleware.JWTMiddleware(jwtSecret))
		protectedLanguage.GET("/preference/:user_id", languageHandler.GetUserLanguagePreference)
		protectedLanguage.POST("/preference/:user_id", languageHandler.SetUserLanguagePreference)
	}
} 