package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type LanguageHTTPHandler struct{}

func NewLanguageHTTPHandler() *LanguageHTTPHandler {
	return &LanguageHTTPHandler{}
}

func (h *LanguageHTTPHandler) ListLanguages(c *gin.Context) {
	languages := []map[string]string{
		{"code": "id", "name": "Indonesian", "native": "Bahasa Indonesia"},
		{"code": "en", "name": "English", "native": "English"},
		{"code": "zh", "name": "Chinese", "native": "中文"},
	}
	c.JSON(http.StatusOK, gin.H{
		"languages": languages,
		"count":     len(languages),
	})
}

func (h *LanguageHTTPHandler) TranslateText(c *gin.Context) {
	var req struct {
		Text         string `json:"text"`
		FromLanguage string `json:"from_language"`
		ToLanguage   string `json:"to_language"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Mock translation
	translatedText := req.Text
	if req.FromLanguage == "en" && req.ToLanguage == "id" {
		switch req.Text {
		case "Hello":
			translatedText = "Halo"
		case "Thank you":
			translatedText = "Terima kasih"
		case "Good morning":
			translatedText = "Selamat pagi"
		case "Hello World":
			translatedText = "Halo Dunia"
		default:
			translatedText = "Terjemahan: " + req.Text
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"original_text":    req.Text,
		"translated_text":  translatedText,
		"from_language":    req.FromLanguage,
		"to_language":      req.ToLanguage,
		"confidence":       0.95,
		"timestamp":        time.Now().Unix(),
	})
}

func (h *LanguageHTTPHandler) DetectLanguage(c *gin.Context) {
	var req struct {
		Text string `json:"text"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Mock language detection
	detectedLang := "en"
	confidence := 0.90

	// Simple detection based on common words
	if h.containsIndonesian(req.Text) {
		detectedLang = "id"
		confidence = 0.95
	} else if h.containsChinese(req.Text) {
		detectedLang = "zh"
		confidence = 0.92
	}

	c.JSON(http.StatusOK, gin.H{
		"text":              req.Text,
		"detected_language": detectedLang,
		"confidence":        confidence,
		"timestamp":         time.Now().Unix(),
	})
}

func (h *LanguageHTTPHandler) BulkTranslate(c *gin.Context) {
	var req struct {
		Texts        []string `json:"texts"`
		FromLanguage string   `json:"from_language"`
		ToLanguage   string   `json:"to_language"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Mock bulk translation
	translations := make([]string, len(req.Texts))
	for i, text := range req.Texts {
		if req.FromLanguage == "en" && req.ToLanguage == "id" {
			switch text {
			case "Hello":
				translations[i] = "Halo"
			case "Thank you":
				translations[i] = "Terima kasih"
			case "Good morning":
				translations[i] = "Selamat pagi"
			default:
				translations[i] = "Terjemahan: " + text
			}
		} else {
			translations[i] = text
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"original_texts":   req.Texts,
		"translated_texts": translations,
		"from_language":    req.FromLanguage,
		"to_language":      req.ToLanguage,
		"count":            len(translations),
		"timestamp":        time.Now().Unix(),
	})
}

func (h *LanguageHTTPHandler) containsIndonesian(text string) bool {
	indonesianWords := []string{"dan", "atau", "dengan", "untuk", "dari", "ke", "di", "pada", "yang", "adalah"}
	lowerText := strings.ToLower(text)
	for _, word := range indonesianWords {
		if strings.Contains(lowerText, word) {
			return true
		}
	}
	return false
}

func (h *LanguageHTTPHandler) containsChinese(text string) bool {
	for _, r := range text {
		if r >= '\u4e00' && r <= '\u9fff' {
			return true
		}
	}
	return false
} 