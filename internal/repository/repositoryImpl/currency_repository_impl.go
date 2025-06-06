package repositoryImpl

import (
	"context"
	"fmt"
	"time"

	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/repository"
	"gorm.io/gorm"
)

// CurrencyRepositoryImpl implements CurrencyRepository using GORM
type CurrencyRepositoryImpl struct {
	db *gorm.DB
}

// NewCurrencyRepository creates a new currency repository instance
func NewCurrencyRepository(db *gorm.DB) repository.CurrencyRepository {
	return &CurrencyRepositoryImpl{db: db}
}

// GetExchangeRate retrieves exchange rate from database
func (r *CurrencyRepositoryImpl) GetExchangeRate(ctx context.Context, from, to models.SupportedCurrency) (float64, error) {
	var rate models.CurrencyRate
	
	// Look for direct rate
	err := r.db.WithContext(ctx).Where(
		"from_currency = ? AND to_currency = ? AND is_active = ?",
		from, to, true,
	).Order("updated_at DESC").First(&rate).Error
	
	if err == nil {
		// Check if rate is not too old (older than 1 hour)
		if time.Now().Unix()-rate.LastUpdated < 3600 {
			return rate.Rate, nil
		}
	}
	
	// Try inverse rate if direct rate not found or too old
	err = r.db.WithContext(ctx).Where(
		"from_currency = ? AND to_currency = ? AND is_active = ?",
		to, from, true,
	).Order("updated_at DESC").First(&rate).Error
	
	if err == nil && rate.Rate > 0 {
		// Check if inverse rate is not too old
		if time.Now().Unix()-rate.LastUpdated < 3600 {
			return 1.0 / rate.Rate, nil
		}
	}
	
	return 0, fmt.Errorf("exchange rate not found for %s to %s", from, to)
}

// SaveExchangeRate saves exchange rate to database
func (r *CurrencyRepositoryImpl) SaveExchangeRate(ctx context.Context, rate *models.CurrencyRate) error {
	// Check if rate already exists
	var existingRate models.CurrencyRate
	err := r.db.WithContext(ctx).Where(
		"from_currency = ? AND to_currency = ?",
		rate.FromCurrency, rate.ToCurrency,
	).First(&existingRate).Error
	
	if err == gorm.ErrRecordNotFound {
		// Create new rate
		return r.db.WithContext(ctx).Create(rate).Error
	} else if err != nil {
		return err
	}
	
	// Update existing rate
	existingRate.Rate = rate.Rate
	existingRate.Source = rate.Source
	existingRate.IsActive = rate.IsActive
	existingRate.LastUpdated = rate.LastUpdated
	
	return r.db.WithContext(ctx).Save(&existingRate).Error
}

// GetAllExchangeRates retrieves all active exchange rates
func (r *CurrencyRepositoryImpl) GetAllExchangeRates(ctx context.Context) ([]*models.CurrencyRate, error) {
	var rates []*models.CurrencyRate
	
	err := r.db.WithContext(ctx).Where("is_active = ?", true).
		Order("updated_at DESC").Find(&rates).Error
	
	return rates, err
}

// GetCurrencyInfo retrieves currency information
func (r *CurrencyRepositoryImpl) GetCurrencyInfo(ctx context.Context, currency models.SupportedCurrency) (*models.CurrencyInfo, error) {
	var info models.CurrencyInfo
	
	err := r.db.WithContext(ctx).Where("code = ? AND is_active = ?", currency, true).
		First(&info).Error
	
	if err == gorm.ErrRecordNotFound {
		// Create default currency info if not found
		info = models.CurrencyInfo{
			Code:        currency,
			Name:        currency.GetName(),
			Symbol:      currency.GetSymbol(),
			DecimalUnit: 2,
			IsActive:    true,
			Country:     currency.GetCountry(),
			Region:      currency.GetRegion(),
		}
		
		if err := r.db.WithContext(ctx).Create(&info).Error; err != nil {
			return nil, err
		}
		
		return &info, nil
	}
	
	return &info, err
}

// SaveCurrencyInfo saves currency information
func (r *CurrencyRepositoryImpl) SaveCurrencyInfo(ctx context.Context, info *models.CurrencyInfo) error {
	var existing models.CurrencyInfo
	err := r.db.WithContext(ctx).Where("code = ?", info.Code).First(&existing).Error
	
	if err == gorm.ErrRecordNotFound {
		return r.db.WithContext(ctx).Create(info).Error
	} else if err != nil {
		return err
	}
	
	// Update existing
	existing.Name = info.Name
	existing.Symbol = info.Symbol
	existing.DecimalUnit = info.DecimalUnit
	existing.IsActive = info.IsActive
	existing.Country = info.Country
	existing.Region = info.Region
	
	return r.db.WithContext(ctx).Save(&existing).Error
}

// GetUserCurrencyPreference retrieves user's currency preference
func (r *CurrencyRepositoryImpl) GetUserCurrencyPreference(ctx context.Context, userID uint) (*models.UserCurrencyPreference, error) {
	var pref models.UserCurrencyPreference
	
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&pref).Error
	
	if err == gorm.ErrRecordNotFound {
		// Create default preference
		pref = models.UserCurrencyPreference{
			UserID:             userID,
			PrimaryCurrency:    models.CurrencyIDR,
			SecondaryCurrency:  models.CurrencyUSD,
			AutoConvert:        false,
			ShowBothCurrencies: false,
		}
		
		if err := r.db.WithContext(ctx).Create(&pref).Error; err != nil {
			return nil, err
		}
		
		return &pref, nil
	}
	
	return &pref, err
}

// SaveUserCurrencyPreference saves user's currency preference
func (r *CurrencyRepositoryImpl) SaveUserCurrencyPreference(ctx context.Context, pref *models.UserCurrencyPreference) error {
	var existing models.UserCurrencyPreference
	err := r.db.WithContext(ctx).Where("user_id = ?", pref.UserID).First(&existing).Error
	
	if err == gorm.ErrRecordNotFound {
		return r.db.WithContext(ctx).Create(pref).Error
	} else if err != nil {
		return err
	}
	
	// Update existing
	existing.PrimaryCurrency = pref.PrimaryCurrency
	existing.SecondaryCurrency = pref.SecondaryCurrency
	existing.AutoConvert = pref.AutoConvert
	existing.ShowBothCurrencies = pref.ShowBothCurrencies
	
	return r.db.WithContext(ctx).Save(&existing).Error
} 