package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/repository"
)

// CurrencyService defines the interface for currency operations
type CurrencyService interface {
	// Exchange rate operations
	GetExchangeRate(ctx context.Context, from, to models.SupportedCurrency) (float64, error)
	UpdateExchangeRates(ctx context.Context) error
	
	// Currency conversion
	ConvertAmount(ctx context.Context, amount float64, from, to models.SupportedCurrency) (float64, error)
	
	// Currency information
	GetSupportedCurrencies(ctx context.Context) ([]models.SupportedCurrency, error)
	
	// Formatting
	FormatCurrency(amount float64, currency models.SupportedCurrency) string
}

// ConvertedAmount represents a converted currency amount
type ConvertedAmount struct {
	OriginalAmount   float64                    `json:"original_amount"`
	OriginalCurrency models.SupportedCurrency  `json:"original_currency"`
	ConvertedAmount  float64                    `json:"converted_amount"`
	TargetCurrency   models.SupportedCurrency  `json:"target_currency"`
	ExchangeRate     float64                    `json:"exchange_rate"`
	ConvertedAt      time.Time                  `json:"converted_at"`
}

// ExchangeRateResponse represents response from external API
type ExchangeRateResponse struct {
	Success   bool                           `json:"success"`
	Timestamp int64                          `json:"timestamp"`
	Base      string                         `json:"base"`
	Date      string                         `json:"date"`
	Rates     map[string]float64             `json:"rates"`
}

// CurrencyFormatter provides currency formatting utilities
type CurrencyFormatter struct {
	DecimalPlaces map[models.SupportedCurrency]int
	Separators    map[models.SupportedCurrency]CurrencySeparator
}

type CurrencySeparator struct {
	Decimal   string // "." or ","
	Thousand  string // "," or "."
}

// GetDefaultCurrencyFormatter returns a default currency formatter
func GetDefaultCurrencyFormatter() *CurrencyFormatter {
	return &CurrencyFormatter{
		DecimalPlaces: map[models.SupportedCurrency]int{
			models.CurrencyIDR: 0, // Indonesian Rupiah - no decimals
			models.CurrencyUSD: 2, // US Dollar - 2 decimals
			models.CurrencyCNY: 2, // Chinese Yuan - 2 decimals
			models.CurrencyEUR: 2, // Euro - 2 decimals
			models.CurrencyJPY: 0, // Japanese Yen - no decimals
			models.CurrencySGD: 2, // Singapore Dollar - 2 decimals
			models.CurrencyMYR: 2, // Malaysian Ringgit - 2 decimals
		},
		Separators: map[models.SupportedCurrency]CurrencySeparator{
			models.CurrencyIDR: {Decimal: ",", Thousand: "."},  // Indonesian format
			models.CurrencyUSD: {Decimal: ".", Thousand: ","},  // US format
			models.CurrencyCNY: {Decimal: ".", Thousand: ","},  // Chinese format
			models.CurrencyEUR: {Decimal: ",", Thousand: "."},  // European format
			models.CurrencyJPY: {Decimal: ".", Thousand: ","},  // Japanese format
			models.CurrencySGD: {Decimal: ".", Thousand: ","},  // Singapore format
			models.CurrencyMYR: {Decimal: ".", Thousand: ","},  // Malaysian format
		},
	}
}

// CurrencyAPIConfig represents configuration for external currency API
type CurrencyAPIConfig struct {
	BaseURL string
	APIKey  string
	Timeout time.Duration
}

// GetDefaultCurrencyAPIConfig returns default configuration for currency API
func GetDefaultCurrencyAPIConfig() *CurrencyAPIConfig {
	return &CurrencyAPIConfig{
		BaseURL: "https://api.exchangerate-api.com/v4/latest",
		Timeout: 30 * time.Second,
	}
}

// FetchExchangeRatesFromAPI fetches exchange rates from external API
func FetchExchangeRatesFromAPI(baseCurrency models.SupportedCurrency, config *CurrencyAPIConfig) (*ExchangeRateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	url := fmt.Sprintf("%s/%s", config.BaseURL, string(baseCurrency))
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+config.APIKey)
	}

	client := &http.Client{Timeout: config.Timeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exchange rates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var apiResponse ExchangeRateResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &apiResponse, nil
}

// ValidateCurrency checks if a currency is supported
func ValidateCurrency(currency models.SupportedCurrency) error {
	supportedCurrencies := []models.SupportedCurrency{
		models.CurrencyIDR,
		models.CurrencyUSD,
		models.CurrencyCNY,
		models.CurrencyEUR,
		models.CurrencyJPY,
		models.CurrencySGD,
		models.CurrencyMYR,
	}

	for _, supported := range supportedCurrencies {
		if currency == supported {
			return nil
		}
	}

	return fmt.Errorf("unsupported currency: %s", string(currency))
}

// ExchangeRateAPIResponse represents the response from ExchangeRate-API
type ExchangeRateAPIResponse struct {
	Result      string             `json:"result"`
	BaseCode    string             `json:"base_code"`
	TargetCode  string             `json:"target_code"`
	ConversionRate float64         `json:"conversion_rate"`
	ConversionResult float64       `json:"conversion_result"`
	TimeLastUpdateUTC string       `json:"time_last_update_utc"`
	TimeNextUpdateUTC string       `json:"time_next_update_utc"`
}

// LatestRatesResponse represents the latest rates response
type LatestRatesResponse struct {
	Result             string                       `json:"result"`
	Documentation      string                       `json:"documentation"`
	TermsOfUse         string                       `json:"terms_of_use"`
	TimeLastUpdateUTC  string                       `json:"time_last_update_utc"`
	TimeNextUpdateUTC  string                       `json:"time_next_update_utc"`
	BaseCode           string                       `json:"base_code"`
	ConversionRates    map[string]float64           `json:"conversion_rates"`
}

// CurrencyServiceImpl implements CurrencyService interface
type CurrencyServiceImpl struct {
	currencyRepo repository.CurrencyRepository
	httpClient   *http.Client
	apiKey       string // Optional API key for premium features
}

// NewCurrencyService creates a new currency service instance
func NewCurrencyService(currencyRepo repository.CurrencyRepository) CurrencyService {
	return &CurrencyServiceImpl{
		currencyRepo: currencyRepo,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		apiKey: "", // Free tier doesn't need API key
	}
}

// GetExchangeRate gets real-time exchange rate from external API
func (s *CurrencyServiceImpl) GetExchangeRate(ctx context.Context, from, to models.SupportedCurrency) (float64, error) {
	// Try to get from cache/database first
	if rate, err := s.currencyRepo.GetExchangeRate(ctx, from, to); err == nil && rate > 0 {
		return rate, nil
	}

	// Fetch from external API
	url := fmt.Sprintf("https://v6.exchangerate-api.com/v6/latest/%s", from)
	
	resp, err := s.httpClient.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch exchange rate: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var apiResp LatestRatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return 0, fmt.Errorf("failed to decode API response: %w", err)
	}

	if apiResp.Result != "success" {
		return 0, fmt.Errorf("API returned error result")
	}

	rate, exists := apiResp.ConversionRates[string(to)]
	if !exists {
		return 0, fmt.Errorf("currency %s not found in rates", to)
	}

	// Cache the rate in database
	go s.cacheExchangeRate(ctx, from, to, rate)

	return rate, nil
}

// ConvertAmount converts amount from one currency to another
func (s *CurrencyServiceImpl) ConvertAmount(ctx context.Context, amount float64, from, to models.SupportedCurrency) (float64, error) {
	if from == to {
		return amount, nil
	}

	rate, err := s.GetExchangeRate(ctx, from, to)
	if err != nil {
		return 0, err
	}

	return amount * rate, nil
}

// FormatCurrency formats amount according to currency rules
func (s *CurrencyServiceImpl) FormatCurrency(amount float64, currency models.SupportedCurrency) string {
	symbol := currency.GetSymbol()
	
	switch currency {
	case models.CurrencyIDR:
		// Indonesian format: Rp 25.000 (no decimals, dot thousands separator)
		return fmt.Sprintf("%s %.0f", symbol, amount)
	case models.CurrencyJPY:
		// Japanese format: ¥25 (no decimals)
		return fmt.Sprintf("%s%.0f", symbol, amount)
	default:
		// Default format: $25.00 (2 decimals)
		return fmt.Sprintf("%s%.2f", symbol, amount)
	}
}

// GetSupportedCurrencies returns list of supported currencies
func (s *CurrencyServiceImpl) GetSupportedCurrencies(ctx context.Context) ([]models.SupportedCurrency, error) {
	return []models.SupportedCurrency{
		models.CurrencyIDR,
		models.CurrencyUSD,
		models.CurrencyCNY,
		models.CurrencyEUR,
		models.CurrencyJPY,
		models.CurrencySGD,
		models.CurrencyMYR,
	}, nil
}

// UpdateExchangeRates fetches and updates all exchange rates
func (s *CurrencyServiceImpl) UpdateExchangeRates(ctx context.Context) error {
	supportedCurrencies, err := s.GetSupportedCurrencies(ctx)
	if err != nil {
		return err
	}

	// Use USD as base currency for updates
	baseCurrency := models.CurrencyUSD
	
	for _, targetCurrency := range supportedCurrencies {
		if targetCurrency == baseCurrency {
			continue
		}

		rate, err := s.GetExchangeRate(ctx, baseCurrency, targetCurrency)
		if err != nil {
			fmt.Printf("Failed to update rate for %s: %v\n", targetCurrency, err)
			continue
		}

		// Save to database
		err = s.currencyRepo.SaveExchangeRate(ctx, &models.CurrencyRate{
			FromCurrency: baseCurrency,
			ToCurrency:   targetCurrency,
			Rate:         rate,
			Source:       "exchangerate-api",
			IsActive:     true,
			LastUpdated:  time.Now().Unix(),
		})
		if err != nil {
			fmt.Printf("Failed to save rate for %s: %v\n", targetCurrency, err)
		}
	}

	return nil
}

// cacheExchangeRate saves exchange rate to database for caching
func (s *CurrencyServiceImpl) cacheExchangeRate(ctx context.Context, from, to models.SupportedCurrency, rate float64) {
	currencyRate := &models.CurrencyRate{
		FromCurrency: from,
		ToCurrency:   to,
		Rate:         rate,
		Source:       "exchangerate-api",
		IsActive:     true,
		LastUpdated:  time.Now().Unix(),
	}

	if err := s.currencyRepo.SaveExchangeRate(ctx, currencyRate); err != nil {
		fmt.Printf("Failed to cache exchange rate: %v\n", err)
	}
} 