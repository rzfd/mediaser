package models

// SupportedLanguage represents supported languages in the system
type SupportedLanguage string

const (
	LanguageIndonesian SupportedLanguage = "id" // Indonesian
	LanguageEnglish    SupportedLanguage = "en" // English
	LanguageMandarin   SupportedLanguage = "zh" // Chinese (Mandarin)
)

// LanguageConfig represents language configuration and translations
type LanguageConfig struct {
	Base
	Language    SupportedLanguage `json:"language" gorm:"not null;index"`
	Key         string            `json:"key" gorm:"not null;index"`
	Translation string            `json:"translation" gorm:"not null"`
	Category    string            `json:"category" gorm:"index"` // ui, messages, errors, etc
	Module      string            `json:"module" gorm:"index"`   // donation, payment, auth, etc
}

// LanguageInfo represents detailed language information
type LanguageInfo struct {
	Base
	Code         SupportedLanguage `json:"code" gorm:"unique;not null"`
	Name         string            `json:"name" gorm:"not null"`         // e.g., "Indonesian", "English"
	NativeName   string            `json:"native_name" gorm:"not null"`  // e.g., "Bahasa Indonesia", "中文"
	Flag         string            `json:"flag"`                         // emoji flag or image path
	IsActive     bool              `json:"is_active" gorm:"default:true"`
	IsRTL        bool              `json:"is_rtl" gorm:"default:false"`  // Right-to-Left languages
	DateFormat   string            `json:"date_format" gorm:"default:'DD/MM/YYYY'"`
	TimeFormat   string            `json:"time_format" gorm:"default:'HH:mm'"`
	NumberFormat string            `json:"number_format" gorm:"default:'1,234.56'"`
}

// UserLanguagePreference represents user's language preferences
type UserLanguagePreference struct {
	Base
	UserID          uint              `json:"user_id" gorm:"not null;unique;index"`
	PrimaryLanguage SupportedLanguage `json:"primary_language" gorm:"default:'id'"`
	FallbackLanguage SupportedLanguage `json:"fallback_language" gorm:"default:'en'"`
	AutoDetect      bool              `json:"auto_detect" gorm:"default:true"`
	Timezone        string            `json:"timezone" gorm:"default:'Asia/Jakarta'"`
}

// TableName specifies the table name for LanguageConfig
func (LanguageConfig) TableName() string {
	return "language_configs"
}

// TableName specifies the table name for LanguageInfo
func (LanguageInfo) TableName() string {
	return "language_info"
}

// TableName specifies the table name for UserLanguagePreference
func (UserLanguagePreference) TableName() string {
	return "user_language_preferences"
}

// LanguageMetadata holds language information for cleaner lookup
type LanguageMetadata struct {
	Name         string
	NativeName   string
	Flag         string
	IsRTL        bool
	DateFormat   string
	NumberFormat string
}

// languageMap provides O(1) lookup for language metadata
var languageMap = map[SupportedLanguage]LanguageMetadata{
	LanguageIndonesian: {
		Name:         "Indonesian",
		NativeName:   "Bahasa Indonesia",
		Flag:         "🇮🇩",
		IsRTL:        false,
		DateFormat:   "DD/MM/YYYY",
		NumberFormat: "1.234.567,89",
	},
	LanguageEnglish: {
		Name:         "English",
		NativeName:   "English",
		Flag:         "🇺🇸",
		IsRTL:        false,
		DateFormat:   "MM/DD/YYYY",
		NumberFormat: "1,234,567.89",
	},
	LanguageMandarin: {
		Name:         "Chinese (Mandarin)",
		NativeName:   "中文",
		Flag:         "🇨🇳",
		IsRTL:        false,
		DateFormat:   "YYYY/MM/DD",
		NumberFormat: "1,234,567.89",
	},
}

// GetName returns the English name for a given language
func (l SupportedLanguage) GetName() string {
	if metadata, exists := languageMap[l]; exists {
		return metadata.Name
	}
	return string(l)
}

// GetNativeName returns the native name for a given language
func (l SupportedLanguage) GetNativeName() string {
	if metadata, exists := languageMap[l]; exists {
		return metadata.NativeName
	}
	return string(l)
}

// GetFlag returns the flag emoji for a given language
func (l SupportedLanguage) GetFlag() string {
	if metadata, exists := languageMap[l]; exists {
		return metadata.Flag
	}
	return "🌐"
}

// IsRightToLeft returns if the language is written right-to-left
func (l SupportedLanguage) IsRightToLeft() bool {
	if metadata, exists := languageMap[l]; exists {
		return metadata.IsRTL
	}
	return false
}

// GetDateFormat returns the date format for a given language
func (l SupportedLanguage) GetDateFormat() string {
	if metadata, exists := languageMap[l]; exists {
		return metadata.DateFormat
	}
	return "DD/MM/YYYY"
}

// GetNumberFormat returns the number format for a given language
func (l SupportedLanguage) GetNumberFormat() string {
	if metadata, exists := languageMap[l]; exists {
		return metadata.NumberFormat
	}
	return "1,234.56"
} 