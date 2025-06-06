package repositoryImpl

import (
	"context"
	"fmt"

	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/repository"
	"gorm.io/gorm"
)

// Use the languageMap from models package
var languageMap = map[models.SupportedLanguage]models.LanguageMetadata{
	models.LanguageIndonesian: {
		Name:         "Indonesian",
		NativeName:   "Bahasa Indonesia",
		Flag:         "🇮🇩",
		IsRTL:        false,
		DateFormat:   "DD/MM/YYYY",
		NumberFormat: "1.234.567,89",
	},
	models.LanguageEnglish: {
		Name:         "English",
		NativeName:   "English",
		Flag:         "🇺🇸",
		IsRTL:        false,
		DateFormat:   "MM/DD/YYYY",
		NumberFormat: "1,234,567.89",
	},
	models.LanguageMandarin: {
		Name:         "Chinese (Mandarin)",
		NativeName:   "中文",
		Flag:         "🇨🇳",
		IsRTL:        false,
		DateFormat:   "YYYY/MM/DD",
		NumberFormat: "1,234,567.89",
	},
}

// LanguageRepositoryImpl implements LanguageRepository using GORM
type LanguageRepositoryImpl struct {
	db *gorm.DB
}

// NewLanguageRepository creates a new language repository instance
func NewLanguageRepository(db *gorm.DB) repository.LanguageRepository {
	return &LanguageRepositoryImpl{db: db}
}

// GetTranslation retrieves translation by key and language
func (r *LanguageRepositoryImpl) GetTranslation(ctx context.Context, key string, language models.SupportedLanguage) (string, error) {
	var config models.LanguageConfig
	
	err := r.db.WithContext(ctx).Where(
		"language = ? AND key = ?",
		language, key,
	).First(&config).Error
	
	if err != nil {
		return "", err
	}
	
	return config.Translation, nil
}

// SaveTranslation saves translation to database
func (r *LanguageRepositoryImpl) SaveTranslation(ctx context.Context, config *models.LanguageConfig) error {
	var existing models.LanguageConfig
	err := r.db.WithContext(ctx).Where(
		"language = ? AND key = ?",
		config.Language, config.Key,
	).First(&existing).Error
	
	if err == gorm.ErrRecordNotFound {
		return r.db.WithContext(ctx).Create(config).Error
	} else if err != nil {
		return err
	}
	
	// Update existing translation
	existing.Translation = config.Translation
	existing.Category = config.Category
	existing.Module = config.Module
	
	return r.db.WithContext(ctx).Save(&existing).Error
}

// GetCachedTranslation retrieves cached translation
func (r *LanguageRepositoryImpl) GetCachedTranslation(ctx context.Context, originalText string, fromLang, toLang models.SupportedLanguage) (string, error) {
	// Create cache key
	cacheKey := fmt.Sprintf("cache_%s_%s_%s", originalText, fromLang, toLang)
	
	var config models.LanguageConfig
	err := r.db.WithContext(ctx).Where(
		"language = ? AND key = ? AND category = ?",
		toLang, cacheKey, "cache",
	).First(&config).Error
	
	if err != nil {
		return "", err
	}
	
	return config.Translation, nil
}

// GetAllTranslations retrieves all translations for a language
func (r *LanguageRepositoryImpl) GetAllTranslations(ctx context.Context, language models.SupportedLanguage) ([]*models.LanguageConfig, error) {
	var configs []*models.LanguageConfig
	
	err := r.db.WithContext(ctx).Where("language = ?", language).
		Order("category, module, key").Find(&configs).Error
	
	return configs, err
}

// GetLanguageInfo retrieves language information
func (r *LanguageRepositoryImpl) GetLanguageInfo(ctx context.Context, language models.SupportedLanguage) (*models.LanguageInfo, error) {
	var info models.LanguageInfo
	
	err := r.db.WithContext(ctx).Where("code = ? AND is_active = ?", language, true).
		First(&info).Error
	
	if err == gorm.ErrRecordNotFound {
		// Create default language info if not found
		metadata := models.LanguageMetadata{}
		if langMeta, exists := languageMap[language]; exists {
			metadata = langMeta
		}
		
		info = models.LanguageInfo{
			Code:         language,
			Name:         language.GetName(),
			NativeName:   language.GetNativeName(),
			Flag:         language.GetFlag(),
			IsActive:     true,
			IsRTL:        metadata.IsRTL,
			DateFormat:   metadata.DateFormat,
			TimeFormat:   "HH:mm",
			NumberFormat: metadata.NumberFormat,
		}
		
		if err := r.db.WithContext(ctx).Create(&info).Error; err != nil {
			return nil, err
		}
		
		return &info, nil
	}
	
	return &info, err
}

// SaveLanguageInfo saves language information
func (r *LanguageRepositoryImpl) SaveLanguageInfo(ctx context.Context, info *models.LanguageInfo) error {
	var existing models.LanguageInfo
	err := r.db.WithContext(ctx).Where("code = ?", info.Code).First(&existing).Error
	
	if err == gorm.ErrRecordNotFound {
		return r.db.WithContext(ctx).Create(info).Error
	} else if err != nil {
		return err
	}
	
	// Update existing
	existing.Name = info.Name
	existing.NativeName = info.NativeName
	existing.Flag = info.Flag
	existing.IsActive = info.IsActive
	existing.IsRTL = info.IsRTL
	existing.DateFormat = info.DateFormat
	existing.TimeFormat = info.TimeFormat
	existing.NumberFormat = info.NumberFormat
	
	return r.db.WithContext(ctx).Save(&existing).Error
}

// GetUserLanguagePreference retrieves user's language preference
func (r *LanguageRepositoryImpl) GetUserLanguagePreference(ctx context.Context, userID uint) (*models.UserLanguagePreference, error) {
	var pref models.UserLanguagePreference
	
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&pref).Error
	
	if err == gorm.ErrRecordNotFound {
		// Create default preference
		pref = models.UserLanguagePreference{
			UserID:           userID,
			PrimaryLanguage:  models.LanguageIndonesian,
			FallbackLanguage: models.LanguageEnglish,
			AutoDetect:       true,
			Timezone:         "Asia/Jakarta",
		}
		
		if err := r.db.WithContext(ctx).Create(&pref).Error; err != nil {
			return nil, err
		}
		
		return &pref, nil
	}
	
	return &pref, err
}

// SaveUserLanguagePreference saves user's language preference
func (r *LanguageRepositoryImpl) SaveUserLanguagePreference(ctx context.Context, pref *models.UserLanguagePreference) error {
	var existing models.UserLanguagePreference
	err := r.db.WithContext(ctx).Where("user_id = ?", pref.UserID).First(&existing).Error
	
	if err == gorm.ErrRecordNotFound {
		return r.db.WithContext(ctx).Create(pref).Error
	} else if err != nil {
		return err
	}
	
	// Update existing
	existing.PrimaryLanguage = pref.PrimaryLanguage
	existing.FallbackLanguage = pref.FallbackLanguage
	existing.AutoDetect = pref.AutoDetect
	existing.Timezone = pref.Timezone
	
	return r.db.WithContext(ctx).Save(&existing).Error
}

// BulkSaveTranslations saves multiple translations at once
func (r *LanguageRepositoryImpl) BulkSaveTranslations(ctx context.Context, configs []*models.LanguageConfig) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, config := range configs {
			var existing models.LanguageConfig
			err := tx.Where(
				"language = ? AND key = ?",
				config.Language, config.Key,
			).First(&existing).Error
			
			if err == gorm.ErrRecordNotFound {
				if err := tx.Create(config).Error; err != nil {
					return err
				}
			} else if err != nil {
				return err
			} else {
				// Update existing
				existing.Translation = config.Translation
				existing.Category = config.Category
				existing.Module = config.Module
				
				if err := tx.Save(&existing).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
} 