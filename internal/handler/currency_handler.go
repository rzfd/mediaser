package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/service"
)

// CurrencyHandler handles currency-related HTTP requests using Echo framework
type CurrencyHandler struct {
	currencyService service.CurrencyService
}

// NewCurrencyHandler creates a new currency handler
func NewCurrencyHandler(currencyService service.CurrencyService) *CurrencyHandler {
	return &CurrencyHandler{
		currencyService: currencyService,
	}
}

// ConvertCurrencyRequest represents currency conversion request
type ConvertCurrencyRequest struct {
	Amount float64                   `json:"amount" validate:"required,gt=0"`
	From   models.SupportedCurrency `json:"from" validate:"required"`
	To     models.SupportedCurrency `json:"to" validate:"required"`
}

// ConvertCurrencyResponse represents currency conversion response
type ConvertCurrencyResponse struct {
	Success        bool                     `json:"success"`
	Amount         float64                  `json:"amount"`
	From           models.SupportedCurrency `json:"from"`
	To             models.SupportedCurrency `json:"to"`
	ConvertedAmount float64                 `json:"converted_amount"`
	FormattedAmount string                  `json:"formatted_amount"`
	ExchangeRate   float64                  `json:"exchange_rate"`
	Message        string                   `json:"message,omitempty"`
}

// CurrencyListResponse represents supported currencies response
type CurrencyListResponse struct {
	Success    bool `json:"success"`
	Currencies []struct {
		Code    models.SupportedCurrency `json:"code"`
		Name    string                   `json:"name"`
		Symbol  string                   `json:"symbol"`
		Country string                   `json:"country"`
		Region  string                   `json:"region"`
	} `json:"currencies"`
	Message string `json:"message,omitempty"`
}

// UserCurrencyPreferenceRequest represents user preference request
type UserCurrencyPreferenceRequest struct {
	PrimaryCurrency    models.SupportedCurrency `json:"primary_currency" validate:"required"`
	SecondaryCurrency  models.SupportedCurrency `json:"secondary_currency"`
	AutoConvert        bool                     `json:"auto_convert"`
	ShowBothCurrencies bool                     `json:"show_both_currencies"`
}

// ConvertCurrency handles currency conversion
func (h *CurrencyHandler) ConvertCurrency(c echo.Context) error {
	var req ConvertCurrencyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ConvertCurrencyResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
	}

	ctx := context.Background()
	
	// Get exchange rate
	rate, err := h.currencyService.GetExchangeRate(ctx, req.From, req.To)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ConvertCurrencyResponse{
			Success: false,
			Message: "Failed to get exchange rate: " + err.Error(),
		})
	}

	// Convert amount
	convertedAmount, err := h.currencyService.ConvertAmount(ctx, req.Amount, req.From, req.To)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ConvertCurrencyResponse{
			Success: false,
			Message: "Failed to convert currency: " + err.Error(),
		})
	}

	// Format the converted amount
	formattedAmount := h.currencyService.FormatCurrency(convertedAmount, req.To)

	return c.JSON(http.StatusOK, ConvertCurrencyResponse{
		Success:         true,
		Amount:          req.Amount,
		From:            req.From,
		To:              req.To,
		ConvertedAmount: convertedAmount,
		FormattedAmount: formattedAmount,
		ExchangeRate:    rate,
		Message:         "Currency converted successfully",
	})
}

// GetExchangeRate gets exchange rate between two currencies
func (h *CurrencyHandler) GetExchangeRate(c echo.Context) error {
	from := models.SupportedCurrency(c.QueryParam("from"))
	to := models.SupportedCurrency(c.QueryParam("to"))

	if from == "" || to == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Both 'from' and 'to' currency codes are required",
		})
	}

	ctx := context.Background()
	rate, err := h.currencyService.GetExchangeRate(ctx, from, to)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to get exchange rate: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":       true,
		"from":          from,
		"to":            to,
		"exchange_rate": rate,
		"message":       "Exchange rate retrieved successfully",
	})
}

// GetSupportedCurrencies returns list of supported currencies
func (h *CurrencyHandler) GetSupportedCurrencies(c echo.Context) error {
	ctx := context.Background()
	currencies, err := h.currencyService.GetSupportedCurrencies(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, CurrencyListResponse{
			Success: false,
			Message: "Failed to get supported currencies: " + err.Error(),
		})
	}

	var currencyList []struct {
		Code    models.SupportedCurrency `json:"code"`
		Name    string                   `json:"name"`
		Symbol  string                   `json:"symbol"`
		Country string                   `json:"country"`
		Region  string                   `json:"region"`
	}

	for _, currency := range currencies {
		currencyList = append(currencyList, struct {
			Code    models.SupportedCurrency `json:"code"`
			Name    string                   `json:"name"`
			Symbol  string                   `json:"symbol"`
			Country string                   `json:"country"`
			Region  string                   `json:"region"`
		}{
			Code:    currency,
			Name:    currency.GetName(),
			Symbol:  currency.GetSymbol(),
			Country: currency.GetCountry(),
			Region:  currency.GetRegion(),
		})
	}

	return c.JSON(http.StatusOK, CurrencyListResponse{
		Success:    true,
		Currencies: currencyList,
		Message:    "Supported currencies retrieved successfully",
	})
}

// UpdateExchangeRates manually updates all exchange rates
func (h *CurrencyHandler) UpdateExchangeRates(c echo.Context) error {
	ctx := context.Background()
	err := h.currencyService.UpdateExchangeRates(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to update exchange rates: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Exchange rates updated successfully",
	})
}

// GetUserCurrencyPreference gets user's currency preference
func (h *CurrencyHandler) GetUserCurrencyPreference(c echo.Context) error {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid user ID",
		})
	}

	// TODO: Implement get user currency preference from service
	// For now, return a mock response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"user_id": userID,
		"primary_currency": "IDR",
		"secondary_currency": "USD",
		"auto_convert": false,
		"show_both_currencies": false,
		"message": "User currency preference retrieved successfully",
	})
}

// SetUserCurrencyPreference sets user's currency preference
func (h *CurrencyHandler) SetUserCurrencyPreference(c echo.Context) error {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid user ID",
		})
	}

	var req UserCurrencyPreferenceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request format: " + err.Error(),
		})
	}

	// TODO: Implement set user currency preference in service
	// For now, return a mock response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"user_id": userID,
		"primary_currency": req.PrimaryCurrency,
		"secondary_currency": req.SecondaryCurrency,
		"auto_convert": req.AutoConvert,
		"show_both_currencies": req.ShowBothCurrencies,
		"message": "User currency preference updated successfully",
	})
} 