package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CurrencyHTTPHandler struct{}

func NewCurrencyHTTPHandler() *CurrencyHTTPHandler {
	return &CurrencyHTTPHandler{}
}

func (h *CurrencyHTTPHandler) ListCurrencies(c *gin.Context) {
	currencies := []string{"IDR", "USD", "CNY", "EUR", "JPY", "SGD", "MYR"}
	c.JSON(http.StatusOK, gin.H{
		"currencies": currencies,
		"count":      len(currencies),
	})
}

func (h *CurrencyHTTPHandler) GetExchangeRate(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	
	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "from and to parameters are required",
		})
		return
	}

	// Mock exchange rate
	rate := 15000.0 // USD to IDR
	c.JSON(http.StatusOK, gin.H{
		"from_currency": from,
		"to_currency":   to,
		"rate":          rate,
		"timestamp":     time.Now().Unix(),
	})
}

func (h *CurrencyHTTPHandler) ConvertCurrency(c *gin.Context) {
	var req struct {
		Amount       float64 `json:"amount"`
		FromCurrency string  `json:"from_currency"`
		ToCurrency   string  `json:"to_currency"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Mock conversion
	rate := 15000.0 // USD to IDR
	convertedAmount := req.Amount * rate

	c.JSON(http.StatusOK, gin.H{
		"original_amount":   req.Amount,
		"from_currency":     req.FromCurrency,
		"converted_amount":  convertedAmount,
		"to_currency":       req.ToCurrency,
		"exchange_rate":     rate,
		"timestamp":         time.Now().Unix(),
	})
} 