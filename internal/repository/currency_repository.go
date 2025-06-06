package repository

import (
	"context"
	"github.com/rzfd/mediashar/internal/models"
)

// CurrencyRepository defines the interface for currency data operations
type CurrencyRepository interface {
	// GetExchangeRate retrieves exchange rate from database
	GetExchangeRate(ctx context.Context, from, to models.SupportedCurrency) (float64, error)
	
	// SaveExchangeRate saves exchange rate to database
	SaveExchangeRate(ctx context.Context, rate *models.CurrencyRate) error
	
	// GetAllExchangeRates retrieves all active exchange rates
	GetAllExchangeRates(ctx context.Context) ([]*models.CurrencyRate, error)
	
	// GetCurrencyInfo retrieves currency information
	GetCurrencyInfo(ctx context.Context, currency models.SupportedCurrency) (*models.CurrencyInfo, error)
	
	// SaveCurrencyInfo saves currency information
	SaveCurrencyInfo(ctx context.Context, info *models.CurrencyInfo) error
	
	// GetUserCurrencyPreference retrieves user's currency preference
	GetUserCurrencyPreference(ctx context.Context, userID uint) (*models.UserCurrencyPreference, error)
	
	// SaveUserCurrencyPreference saves user's currency preference
	SaveUserCurrencyPreference(ctx context.Context, pref *models.UserCurrencyPreference) error
} 