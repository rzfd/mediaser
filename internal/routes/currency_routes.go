package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/rzfd/mediashar/internal/handler"
	"github.com/rzfd/mediashar/internal/middleware"
)

// SetupCurrencyRoutes configures currency-related routes
func SetupCurrencyRoutes(api *echo.Group, currencyHandler *handler.CurrencyHandler, jwtSecret string) {
	// Currency API routes
	currencyGroup := api.Group("/currency")
	{
		// Public routes
		currencyGroup.POST("/convert", currencyHandler.ConvertCurrency)
		currencyGroup.GET("/rate", currencyHandler.GetExchangeRate)
		currencyGroup.GET("/list", currencyHandler.GetSupportedCurrencies)
		currencyGroup.POST("/update", currencyHandler.UpdateExchangeRates)
		
		// Protected routes (authentication required)
		protectedCurrency := currencyGroup.Group("", middleware.JWTMiddleware(jwtSecret))
		protectedCurrency.GET("/preference/:user_id", currencyHandler.GetUserCurrencyPreference)
		protectedCurrency.POST("/preference/:user_id", currencyHandler.SetUserCurrencyPreference)
	}
} 