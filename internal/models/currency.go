package models

// SupportedCurrency represents supported currencies in the system
type SupportedCurrency string

const (
	CurrencyIDR SupportedCurrency = "IDR" // Indonesian Rupiah
	CurrencyUSD SupportedCurrency = "USD" // US Dollar
	CurrencyCNY SupportedCurrency = "CNY" // Chinese Yuan
	CurrencyEUR SupportedCurrency = "EUR" // Euro
	CurrencyJPY SupportedCurrency = "JPY" // Japanese Yen
	CurrencySGD SupportedCurrency = "SGD" // Singapore Dollar
	CurrencyMYR SupportedCurrency = "MYR" // Malaysian Ringgit
)

// CurrencyRate represents currency exchange rates
type CurrencyRate struct {
	Base
	FromCurrency SupportedCurrency `json:"from_currency" gorm:"not null;index"`
	ToCurrency   SupportedCurrency `json:"to_currency" gorm:"not null;index"`
	Rate         float64           `json:"rate" gorm:"not null"`
	Source       string            `json:"source" gorm:"default:'manual'"` // manual, api, coinbase, etc
	IsActive     bool              `json:"is_active" gorm:"default:true"`
	LastUpdated  int64             `json:"last_updated" gorm:"autoUpdateTime"`
}

// CurrencyInfo represents detailed currency information
type CurrencyInfo struct {
	Base
	Code        SupportedCurrency `json:"code" gorm:"unique;not null"`
	Name        string            `json:"name" gorm:"not null"` // e.g., "Indonesian Rupiah"
	Symbol      string            `json:"symbol"`               // e.g., "Rp", "$", "¥"
	DecimalUnit int               `json:"decimal_unit" gorm:"default:2"`
	IsActive    bool              `json:"is_active" gorm:"default:true"`
	Country     string            `json:"country"`
	Region      string            `json:"region"`
}

// UserCurrencyPreference represents user's currency preferences
type UserCurrencyPreference struct {
	Base
	UserID            uint              `json:"user_id" gorm:"not null;unique;index"`
	PrimaryCurrency   SupportedCurrency `json:"primary_currency" gorm:"default:'IDR'"`
	SecondaryCurrency SupportedCurrency `json:"secondary_currency"`
	AutoConvert       bool              `json:"auto_convert" gorm:"default:false"`
	ShowBothCurrencies bool             `json:"show_both_currencies" gorm:"default:false"`
}

// TableName specifies the table name for CurrencyRate
func (CurrencyRate) TableName() string {
	return "currency_rates"
}

// TableName specifies the table name for CurrencyInfo
func (CurrencyInfo) TableName() string {
	return "currency_info"
}

// TableName specifies the table name for UserCurrencyPreference
func (UserCurrencyPreference) TableName() string {
	return "user_currency_preferences"
}

// CurrencyMetadata holds currency information for cleaner lookup
type CurrencyMetadata struct {
	Symbol  string
	Name    string
	Country string
	Region  string
}

// currencyMap provides O(1) lookup for currency metadata
var currencyMap = map[SupportedCurrency]CurrencyMetadata{
	CurrencyIDR: {"Rp", "Indonesian Rupiah", "Indonesia", "Southeast Asia"},
	CurrencyUSD: {"$", "US Dollar", "United States", "North America"},
	CurrencyCNY: {"¥", "Chinese Yuan", "China", "East Asia"},
	CurrencyEUR: {"€", "Euro", "European Union", "Europe"},
	CurrencyJPY: {"¥", "Japanese Yen", "Japan", "East Asia"},
	CurrencySGD: {"S$", "Singapore Dollar", "Singapore", "Southeast Asia"},
	CurrencyMYR: {"RM", "Malaysian Ringgit", "Malaysia", "Southeast Asia"},
}

// GetSymbol returns the symbol for a given currency
func (c SupportedCurrency) GetSymbol() string {
	if metadata, exists := currencyMap[c]; exists {
		return metadata.Symbol
	}
	return string(c)
}

// GetName returns the full name for a given currency
func (c SupportedCurrency) GetName() string {
	if metadata, exists := currencyMap[c]; exists {
		return metadata.Name
	}
	return string(c)
}

// GetCountry returns the country for a given currency
func (c SupportedCurrency) GetCountry() string {
	if metadata, exists := currencyMap[c]; exists {
		return metadata.Country
	}
	return ""
}

// GetRegion returns the region for a given currency
func (c SupportedCurrency) GetRegion() string {
	if metadata, exists := currencyMap[c]; exists {
		return metadata.Region
	}
	return ""
} 