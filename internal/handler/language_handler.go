package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/service"
)

// LanguageHandler handles language-related HTTP requests using Echo framework
type LanguageHandler struct {
	languageService service.LanguageService
}

// NewLanguageHandler creates a new language handler
func NewLanguageHandler(languageService service.LanguageService) *LanguageHandler {
	return &LanguageHandler{
		languageService: languageService,
	}
}

// TranslateTextRequest represents text translation request
type TranslateTextRequest struct {
	Text     string                    `json:"text" validate:"required"`
	FromLang models.SupportedLanguage  `json:"from_lang" validate:"required"`
	ToLang   models.SupportedLanguage  `json:"to_lang" validate:"required"`
}

// TranslateTextResponse represents text translation response
type TranslateTextResponse struct {
	Success         bool                     `json:"success"`
	OriginalText    string                   `json:"original_text"`
	TranslatedText  string                   `json:"translated_text"`
	FromLang        models.SupportedLanguage `json:"from_lang"`
	ToLang          models.SupportedLanguage `json:"to_lang"`
	Message         string                   `json:"message,omitempty"`
}

// BulkTranslateRequest represents bulk translation request
type BulkTranslateRequest struct {
	Texts    []string                  `json:"texts" validate:"required"`
	FromLang models.SupportedLanguage  `json:"from_lang" validate:"required"`
	ToLang   models.SupportedLanguage  `json:"to_lang" validate:"required"`
}

// BulkTranslateResponse represents bulk translation response
type BulkTranslateResponse struct {
	Success          bool                     `json:"success"`
	OriginalTexts    []string                 `json:"original_texts"`
	TranslatedTexts  []string                 `json:"translated_texts"`
	FromLang         models.SupportedLanguage `json:"from_lang"`
	ToLang           models.SupportedLanguage `json:"to_lang"`
	ProcessedCount   int                      `json:"processed_count"`
	Message          string                   `json:"message,omitempty"`
}

// DetectLanguageRequest represents language detection request
type DetectLanguageRequest struct {
	Text string `json:"text" validate:"required"`
}

// DetectLanguageResponse represents language detection response
type DetectLanguageResponse struct {
	Success          bool                     `json:"success"`
	Text             string                   `json:"text"`
	DetectedLanguage models.SupportedLanguage `json:"detected_language"`
	Confidence       float64                  `json:"confidence"`
	Message          string                   `json:"message,omitempty"`
}

// SupportedLanguagesResponse represents supported languages response
type SupportedLanguagesResponse struct {
	Success   bool `json:"success"`
	Languages []struct {
		Code       models.SupportedLanguage `json:"code"`
		Name       string                   `json:"name"`
		NativeName string                   `json:"native_name"`
		Flag       string                   `json:"flag"`
	} `json:"languages"`
	Message string `json:"message,omitempty"`
}

// UserLanguagePreferenceRequest represents user language preference request
type UserLanguagePreferenceRequest struct {
	PrimaryLanguage  models.SupportedLanguage `json:"primary_language" validate:"required"`
	FallbackLanguage models.SupportedLanguage `json:"fallback_language"`
	AutoDetect       bool                     `json:"auto_detect"`
	Timezone         string                   `json:"timezone"`
}

// TranslateText handles text translation
func (h *LanguageHandler) TranslateText(c echo.Context) error {
	var req TranslateTextRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, TranslateTextResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
	}

	ctx := context.Background()
	
	translatedText, err := h.languageService.TranslateText(ctx, req.Text, req.FromLang, req.ToLang)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, TranslateTextResponse{
			Success: false,
			Message: "Failed to translate text: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, TranslateTextResponse{
		Success:        true,
		OriginalText:   req.Text,
		TranslatedText: translatedText,
		FromLang:       req.FromLang,
		ToLang:         req.ToLang,
		Message:        "Text translated successfully",
	})
}

// BulkTranslate handles bulk text translation
func (h *LanguageHandler) BulkTranslate(c echo.Context) error {
	var req BulkTranslateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, BulkTranslateResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
	}

	if len(req.Texts) > 100 {
		return c.JSON(http.StatusBadRequest, BulkTranslateResponse{
			Success: false,
			Message: "Maximum 100 texts allowed per request",
		})
	}

	ctx := context.Background()
	
	translatedTexts, err := h.languageService.BulkTranslate(ctx, req.Texts, req.FromLang, req.ToLang)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, BulkTranslateResponse{
			Success: false,
			Message: "Failed to translate texts: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, BulkTranslateResponse{
		Success:         true,
		OriginalTexts:   req.Texts,
		TranslatedTexts: translatedTexts,
		FromLang:        req.FromLang,
		ToLang:          req.ToLang,
		ProcessedCount:  len(translatedTexts),
		Message:         "Texts translated successfully",
	})
}

// DetectLanguage handles language detection
func (h *LanguageHandler) DetectLanguage(c echo.Context) error {
	var req DetectLanguageRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, DetectLanguageResponse{
			Success: false,
			Message: "Invalid request format: " + err.Error(),
		})
	}

	ctx := context.Background()
	
	detectedLang, confidence, err := h.languageService.DetectLanguage(ctx, req.Text)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, DetectLanguageResponse{
			Success: false,
			Message: "Failed to detect language: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, DetectLanguageResponse{
		Success:          true,
		Text:             req.Text,
		DetectedLanguage: detectedLang,
		Confidence:       confidence,
		Message:          "Language detected successfully",
	})
}

// GetSupportedLanguages returns list of supported languages
func (h *LanguageHandler) GetSupportedLanguages(c echo.Context) error {
	ctx := context.Background()
	languages, err := h.languageService.GetSupportedLanguages(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to get supported languages: " + err.Error(),
		})
	}

	var languageList []map[string]interface{}
	for _, language := range languages {
		languageList = append(languageList, map[string]interface{}{
			"code":        language,
			"name":        language.GetName(),
			"native_name": language.GetNativeName(),
			"flag":        language.GetFlag(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":   true,
		"languages": languageList,
		"message":   "Supported languages retrieved successfully",
	})
}

// GetTranslation gets system translation for a key
func (h *LanguageHandler) GetTranslation(c echo.Context) error {
	key := c.QueryParam("key")
	lang := models.SupportedLanguage(c.QueryParam("lang"))

	if key == "" || lang == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Both 'key' and 'lang' parameters are required",
		})
	}

	ctx := context.Background()
	translation, err := h.languageService.GetTranslation(ctx, key, lang)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to get translation: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":     true,
		"key":         key,
		"language":    lang,
		"translation": translation,
		"message":     "Translation retrieved successfully",
	})
}

// AddTranslation adds a new system translation
func (h *LanguageHandler) AddTranslation(c echo.Context) error {
	var req map[string]interface{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request format: " + err.Error(),
		})
	}

	key, ok := req["key"].(string)
	if !ok || key == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Key is required",
		})
	}

	langStr, ok := req["language"].(string)
	if !ok || langStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Language is required",
		})
	}

	translation, ok := req["translation"].(string)
	if !ok || translation == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Translation is required",
		})
	}

	ctx := context.Background()
	lang := models.SupportedLanguage(langStr)
	
	err := h.languageService.AddTranslation(ctx, key, lang, translation)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"message": "Failed to add translation: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":     true,
		"key":         key,
		"language":    lang,
		"translation": translation,
		"message":     "Translation added successfully",
	})
}

// GetUserLanguagePreference gets user's language preference
func (h *LanguageHandler) GetUserLanguagePreference(c echo.Context) error {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid user ID",
		})
	}

	// TODO: Implement get user language preference from service
	// For now, return a mock response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"user_id": userID,
		"primary_language": "id",
		"fallback_language": "en",
		"auto_detect": false,
		"timezone": "Asia/Jakarta",
		"message": "User language preference retrieved successfully",
	})
}

// SetUserLanguagePreference sets user's language preference
func (h *LanguageHandler) SetUserLanguagePreference(c echo.Context) error {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid user ID",
		})
	}

	var req UserLanguagePreferenceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"message": "Invalid request format: " + err.Error(),
		})
	}

	// TODO: Implement set user language preference in service
	// For now, return a mock response
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"user_id": userID,
		"primary_language": req.PrimaryLanguage,
		"fallback_language": req.FallbackLanguage,
		"auto_detect": req.AutoDetect,
		"timezone": req.Timezone,
		"message": "User language preference updated successfully",
	})
} 