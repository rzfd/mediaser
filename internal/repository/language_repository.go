package repository

import (
	"context"
	"github.com/rzfd/mediashar/internal/models"
)

// LanguageRepository defines the interface for language data operations
type LanguageRepository interface {
	// GetTranslation retrieves translation by key and language
	GetTranslation(ctx context.Context, key string, language models.SupportedLanguage) (string, error)
	
	// SaveTranslation saves translation to database
	SaveTranslation(ctx context.Context, config *models.LanguageConfig) error
	
	// GetCachedTranslation retrieves cached translation
	GetCachedTranslation(ctx context.Context, originalText string, fromLang, toLang models.SupportedLanguage) (string, error)
	
	// GetAllTranslations retrieves all translations for a language
	GetAllTranslations(ctx context.Context, language models.SupportedLanguage) ([]*models.LanguageConfig, error)
	
	// GetLanguageInfo retrieves language information
	GetLanguageInfo(ctx context.Context, language models.SupportedLanguage) (*models.LanguageInfo, error)
	
	// SaveLanguageInfo saves language information
	SaveLanguageInfo(ctx context.Context, info *models.LanguageInfo) error
	
	// GetUserLanguagePreference retrieves user's language preference
	GetUserLanguagePreference(ctx context.Context, userID uint) (*models.UserLanguagePreference, error)
	
	// SaveUserLanguagePreference saves user's language preference
	SaveUserLanguagePreference(ctx context.Context, pref *models.UserLanguagePreference) error
	
	// BulkSaveTranslations saves multiple translations at once
	BulkSaveTranslations(ctx context.Context, configs []*models.LanguageConfig) error
} 