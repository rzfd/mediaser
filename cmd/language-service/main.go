package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// LanguageService represents a basic language service
type LanguageService struct{}

func main() {
	// Get environment variables
	grpcPort := getEnv("GRPC_PORT", "8085")
	httpPort := getEnv("HTTP_PORT", "8095")

	// Start gRPC server in goroutine
	go startGRPCServer(grpcPort)

	// Start HTTP server for health checks
	startHTTPServer(httpPort)
}

func startGRPCServer(port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port %s: %v", port, err)
	}

	s := grpc.NewServer()
	// Register your language service here when you have protobuf definitions
	// pb.RegisterLanguageServiceServer(s, &LanguageService{})

	log.Printf("Language gRPC server starting on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

func startHTTPServer(port string) {
	// Create Gin router
	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "language-service",
			"status":  "healthy",
			"version": "1.0.0",
			"time":    time.Now().Unix(),
		})
	})

	// Language endpoints (basic implementation)
	api := r.Group("/api/language")
	{
		api.GET("/list", func(c *gin.Context) {
			languages := []map[string]string{
				{"code": "id", "name": "Indonesian", "native": "Bahasa Indonesia"},
				{"code": "en", "name": "English", "native": "English"},
				{"code": "zh", "name": "Chinese", "native": "中文"},
			}
			c.JSON(http.StatusOK, gin.H{
				"languages": languages,
				"count":     len(languages),
			})
		})

		api.POST("/translate", func(c *gin.Context) {
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

			// Mock translation (in real implementation, this would use LibreTranslate API)
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
		})

		api.POST("/detect", func(c *gin.Context) {
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
			if containsIndonesian(req.Text) {
				detectedLang = "id"
				confidence = 0.95
			} else if containsChinese(req.Text) {
				detectedLang = "zh"
				confidence = 0.92
			}

			c.JSON(http.StatusOK, gin.H{
				"text":              req.Text,
				"detected_language": detectedLang,
				"confidence":        confidence,
				"timestamp":         time.Now().Unix(),
			})
		})

		api.POST("/bulk-translate", func(c *gin.Context) {
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
		})
	}

	// Start server
	log.Printf("Language HTTP server starting on port %s", port)
	
	// Graceful shutdown
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Handle shutdown signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Language service shutting down...")
		srv.Close()
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

// Helper functions for simple language detection
func containsIndonesian(text string) bool {
	indonesianWords := []string{"selamat", "terima", "kasih", "pagi", "siang", "malam", "adalah", "dan", "atau", "yang"}
	for _, word := range indonesianWords {
		if contains(text, word) {
			return true
		}
	}
	return false
}

func containsChinese(text string) bool {
	// Simple check for Chinese characters
	for _, r := range text {
		if r >= 0x4e00 && r <= 0x9fff {
			return true
		}
	}
	return false
}

func contains(text, substr string) bool {
	return len(text) >= len(substr) && 
		   (text == substr || 
		    contains(text[1:], substr) || 
		    (len(text) > len(substr) && text[:len(substr)] == substr))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}