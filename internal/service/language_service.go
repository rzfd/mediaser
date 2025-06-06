package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/repository"
)

// LanguageService defines the interface for language operations
type LanguageService interface {
	// External translation API operations
	TranslateText(ctx context.Context, text string, fromLang, toLang models.SupportedLanguage) (string, error)
	BulkTranslate(ctx context.Context, texts []string, fromLang, toLang models.SupportedLanguage) ([]string, error)
	DetectLanguage(ctx context.Context, text string) (models.SupportedLanguage, float64, error)
	
	// System translation operations
	GetTranslation(ctx context.Context, key string, language models.SupportedLanguage) (string, error)
	AddTranslation(ctx context.Context, key string, language models.SupportedLanguage, translation string) error
	
	// Language information
	GetSupportedLanguages(ctx context.Context) ([]models.SupportedLanguage, error)
}

// TranslationRequest represents a translation request
type TranslationRequest struct {
	Key        string                     `json:"key"`
	Language   models.SupportedLanguage   `json:"language"`
	Category   string                     `json:"category,omitempty"`
	Module     string                     `json:"module,omitempty"`
	Fallback   models.SupportedLanguage   `json:"fallback,omitempty"`
}

// TranslationResponse represents a translation response
type TranslationResponse struct {
	Key         string                     `json:"key"`
	Translation string                     `json:"translation"`
	Language    models.SupportedLanguage   `json:"language"`
	Category    string                     `json:"category"`
	Module      string                     `json:"module"`
	Found       bool                       `json:"found"`
}

// BulkTranslationRequest represents a bulk translation request
type BulkTranslationRequest struct {
	Keys       []string                   `json:"keys"`
	Language   models.SupportedLanguage   `json:"language"`
	Category   string                     `json:"category,omitempty"`
	Module     string                     `json:"module,omitempty"`
	Fallback   models.SupportedLanguage   `json:"fallback,omitempty"`
}

// BulkTranslationResponse represents a bulk translation response
type BulkTranslationResponse struct {
	Translations map[string]string         `json:"translations"`
	Language     models.SupportedLanguage  `json:"language"`
	Missing      []string                  `json:"missing,omitempty"`
}

// MessageFormatter provides message formatting utilities
type MessageFormatter struct {
	Placeholders map[string]string
	Functions    map[string]func(interface{}) string
}

// DefaultTranslations provides default translations for common keys
func GetDefaultTranslations() map[models.SupportedLanguage]map[string]string {
	return map[models.SupportedLanguage]map[string]string{
		models.LanguageIndonesian: {
			// Common UI
			"common.save":           "Simpan",
			"common.cancel":         "Batal",
			"common.submit":         "Kirim",
			"common.loading":        "Memuat...",
			"common.error":          "Terjadi kesalahan",
			"common.success":        "Berhasil",
			"common.confirm":        "Konfirmasi",
			"common.delete":         "Hapus",
			"common.edit":           "Edit",
			"common.view":           "Lihat",
			"common.close":          "Tutup",
			
			// Donation related
			"donation.title":        "Donasi",
			"donation.amount":       "Jumlah",
			"donation.message":      "Pesan",
			"donation.submit":       "Kirim Donasi",
			"donation.success":      "Donasi berhasil dikirim",
			"donation.failed":       "Donasi gagal",
			"donation.anonymous":    "Anonim",
			"donation.display_name": "Nama Pengirim",
			
			// Payment related
			"payment.processing":    "Memproses pembayaran...",
			"payment.success":       "Pembayaran berhasil",
			"payment.failed":        "Pembayaran gagal",
			"payment.cancelled":     "Pembayaran dibatalkan",
			
			// Authentication
			"auth.login":            "Masuk",
			"auth.register":         "Daftar",
			"auth.logout":           "Keluar",
			"auth.username":         "Nama Pengguna",
			"auth.email":            "Email",
			"auth.password":         "Kata Sandi",
			"auth.confirm_password": "Konfirmasi Kata Sandi",
		},
		
		models.LanguageEnglish: {
			// Common UI
			"common.save":           "Save",
			"common.cancel":         "Cancel",
			"common.submit":         "Submit",
			"common.loading":        "Loading...",
			"common.error":          "An error occurred",
			"common.success":        "Success",
			"common.confirm":        "Confirm",
			"common.delete":         "Delete",
			"common.edit":           "Edit",
			"common.view":           "View",
			"common.close":          "Close",
			
			// Donation related
			"donation.title":        "Donation",
			"donation.amount":       "Amount",
			"donation.message":      "Message",
			"donation.submit":       "Send Donation",
			"donation.success":      "Donation sent successfully",
			"donation.failed":       "Donation failed",
			"donation.anonymous":    "Anonymous",
			"donation.display_name": "Display Name",
			
			// Payment related
			"payment.processing":    "Processing payment...",
			"payment.success":       "Payment successful",
			"payment.failed":        "Payment failed",
			"payment.cancelled":     "Payment cancelled",
			
			// Authentication
			"auth.login":            "Login",
			"auth.register":         "Register",
			"auth.logout":           "Logout",
			"auth.username":         "Username",
			"auth.email":            "Email",
			"auth.password":         "Password",
			"auth.confirm_password": "Confirm Password",
		},
		
		models.LanguageMandarin: {
			// Common UI
			"common.save":           "保存",
			"common.cancel":         "取消",
			"common.submit":         "提交",
			"common.loading":        "加载中...",
			"common.error":          "发生错误",
			"common.success":        "成功",
			"common.confirm":        "确认",
			"common.delete":         "删除",
			"common.edit":           "编辑",
			"common.view":           "查看",
			"common.close":          "关闭",
			
			// Donation related
			"donation.title":        "捐赠",
			"donation.amount":       "金额",
			"donation.message":      "消息",
			"donation.submit":       "发送捐赠",
			"donation.success":      "捐赠发送成功",
			"donation.failed":       "捐赠失败",
			"donation.anonymous":    "匿名",
			"donation.display_name": "显示名称",
			
			// Payment related
			"payment.processing":    "处理付款中...",
			"payment.success":       "付款成功",
			"payment.failed":        "付款失败",
			"payment.cancelled":     "付款已取消",
			
			// Authentication
			"auth.login":            "登录",
			"auth.register":         "注册",
			"auth.logout":           "登出",
			"auth.username":         "用户名",
			"auth.email":            "电子邮件",
			"auth.password":         "密码",
			"auth.confirm_password": "确认密码",
		},
	}
}

// ValidateLanguage checks if a language is supported
func ValidateLanguage(language models.SupportedLanguage) error {
	supportedLanguages := []models.SupportedLanguage{
		models.LanguageIndonesian,
		models.LanguageEnglish,
		models.LanguageMandarin,
	}

	for _, supported := range supportedLanguages {
		if language == supported {
			return nil
		}
	}

	return fmt.Errorf("unsupported language: %s", string(language))
}

// FormatPlaceholders replaces placeholders in a translation string
func FormatPlaceholders(translation string, params map[string]interface{}) string {
	result := translation
	
	for key, value := range params {
		placeholder := fmt.Sprintf("{%s}", key)
		replacement := fmt.Sprintf("%v", value)
		result = strings.ReplaceAll(result, placeholder, replacement)
	}
	
	return result
}

// DetectLanguageFromHeader attempts to detect language from Accept-Language header
func DetectLanguageFromHeader(acceptLanguage string) models.SupportedLanguage {
	if acceptLanguage == "" {
		return models.LanguageIndonesian // Default to Indonesian
	}
	
	// Parse Accept-Language header (simplified)
	languages := strings.Split(acceptLanguage, ",")
	for _, lang := range languages {
		lang = strings.TrimSpace(strings.Split(lang, ";")[0])
		
		switch {
		case strings.HasPrefix(lang, "id"):
			return models.LanguageIndonesian
		case strings.HasPrefix(lang, "en"):
			return models.LanguageEnglish
		case strings.HasPrefix(lang, "zh"):
			return models.LanguageMandarin
		}
	}
	
	return models.LanguageIndonesian // Default fallback
}

// LibreTranslateRequest represents the request to LibreTranslate API
type LibreTranslateRequest struct {
	Q      string `json:"q"`
	Source string `json:"source"`
	Target string `json:"target"`
	Format string `json:"format,omitempty"`
}

// LibreTranslateResponse represents the response from LibreTranslate API
type LibreTranslateResponse struct {
	TranslatedText string `json:"translatedText"`
	DetectedLanguage *struct {
		Confidence float64 `json:"confidence"`
		Language   string  `json:"language"`
	} `json:"detectedLanguage,omitempty"`
}

// LanguageServiceImpl implements LanguageService interface
type LanguageServiceImpl struct {
	languageRepo repository.LanguageRepository
	httpClient   *http.Client
	apiBaseURL   string
	apiKey       string // Optional API key for premium instances
}

// NewLanguageService creates a new language service instance
func NewLanguageService(languageRepo repository.LanguageRepository) LanguageService {
	return &LanguageServiceImpl{
		languageRepo: languageRepo,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiBaseURL: "https://libretranslate.com", // Free public instance
		apiKey:     "",                           // No API key needed for free tier
	}
}

// TranslateText translates text using external API
func (s *LanguageServiceImpl) TranslateText(ctx context.Context, text string, fromLang, toLang models.SupportedLanguage) (string, error) {
	// Check if translation is cached first
	if cached, err := s.languageRepo.GetCachedTranslation(ctx, text, fromLang, toLang); err == nil && cached != "" {
		return cached, nil
	}

	// Prepare request to LibreTranslate API
	reqData := LibreTranslateRequest{
		Q:      text,
		Source: string(fromLang),
		Target: string(toLang),
		Format: "text",
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/translate", s.apiBaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if s.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+s.apiKey)
	}

	// Send request
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	// Decode response
	var apiResp LibreTranslateResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	// Cache the translation
	go s.cacheTranslation(ctx, text, fromLang, toLang, apiResp.TranslatedText)

	return apiResp.TranslatedText, nil
}

// GetTranslation gets translation from database or API
func (s *LanguageServiceImpl) GetTranslation(ctx context.Context, key string, language models.SupportedLanguage) (string, error) {
	// Try to get from database first
	translation, err := s.languageRepo.GetTranslation(ctx, key, language)
	if err == nil && translation != "" {
		return translation, nil
	}

	// If not found and not the default language, try to translate from English
	if language != models.LanguageEnglish {
		englishText, err := s.languageRepo.GetTranslation(ctx, key, models.LanguageEnglish)
		if err == nil && englishText != "" {
			translated, err := s.TranslateText(ctx, englishText, models.LanguageEnglish, language)
			if err == nil {
				// Save the translation
				go s.saveTranslation(ctx, key, language, translated)
				return translated, nil
			}
		}
	}

	return "", fmt.Errorf("translation not found for key: %s, language: %s", key, language)
}

// GetSupportedLanguages returns list of supported languages
func (s *LanguageServiceImpl) GetSupportedLanguages(ctx context.Context) ([]models.SupportedLanguage, error) {
	return []models.SupportedLanguage{
		models.LanguageIndonesian,
		models.LanguageEnglish,
		models.LanguageMandarin,
	}, nil
}

// AddTranslation adds a new translation to the system
func (s *LanguageServiceImpl) AddTranslation(ctx context.Context, key string, language models.SupportedLanguage, translation string) error {
	return s.languageRepo.SaveTranslation(ctx, &models.LanguageConfig{
		Language:    language,
		Key:         key,
		Translation: translation,
		Category:    "dynamic",
		Module:      "system",
	})
}

// BulkTranslate translates multiple texts at once
func (s *LanguageServiceImpl) BulkTranslate(ctx context.Context, texts []string, fromLang, toLang models.SupportedLanguage) ([]string, error) {
	if len(texts) == 0 {
		return []string{}, nil
	}

	var results []string
	for _, text := range texts {
		translated, err := s.TranslateText(ctx, text, fromLang, toLang)
		if err != nil {
			// If one translation fails, use original text
			results = append(results, text)
			continue
		}
		results = append(results, translated)
	}

	return results, nil
}

// DetectLanguage detects the language of given text
func (s *LanguageServiceImpl) DetectLanguage(ctx context.Context, text string) (models.SupportedLanguage, float64, error) {
	// Use LibreTranslate with auto-detect
	reqData := LibreTranslateRequest{
		Q:      text,
		Source: "auto", // Auto-detect
		Target: "en",   // Translate to English for detection
		Format: "text",
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return "", 0, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/translate", s.apiBaseURL)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", 0, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	var apiResp LibreTranslateResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return "", 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if apiResp.DetectedLanguage != nil {
		// Map detected language to our supported languages
		detectedLang := models.SupportedLanguage(apiResp.DetectedLanguage.Language)
		confidence := apiResp.DetectedLanguage.Confidence
		
		// Validate if it's one of our supported languages
		supported, _ := s.GetSupportedLanguages(ctx)
		for _, lang := range supported {
			if lang == detectedLang {
				return detectedLang, confidence, nil
			}
		}
	}

	// Default to English if detection fails
	return models.LanguageEnglish, 0.5, nil
}

// cacheTranslation saves translation to database for caching
func (s *LanguageServiceImpl) cacheTranslation(ctx context.Context, originalText string, fromLang, toLang models.SupportedLanguage, translation string) {
	// Create a cache key based on the original text and language pair
	cacheKey := fmt.Sprintf("cache_%s_%s_%s", originalText, fromLang, toLang)
	
	err := s.languageRepo.SaveTranslation(ctx, &models.LanguageConfig{
		Language:    toLang,
		Key:         cacheKey,
		Translation: translation,
		Category:    "cache",
		Module:      "translation",
	})
	if err != nil {
		fmt.Printf("Failed to cache translation: %v\n", err)
	}
}

// saveTranslation saves a system translation
func (s *LanguageServiceImpl) saveTranslation(ctx context.Context, key string, language models.SupportedLanguage, translation string) {
	err := s.languageRepo.SaveTranslation(ctx, &models.LanguageConfig{
		Language:    language,
		Key:         key,
		Translation: translation,
		Category:    "system",
		Module:      "auto",
	})
	if err != nil {
		fmt.Printf("Failed to save translation: %v\n", err)
	}
} 