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

// CurrencyService represents a basic currency service
type CurrencyService struct{}

func main() {
	// Get environment variables
	grpcPort := getEnv("GRPC_PORT", "8084")
	httpPort := getEnv("HTTP_PORT", "8094")

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
	// Register your currency service here when you have protobuf definitions
	// pb.RegisterCurrencyServiceServer(s, &CurrencyService{})

	log.Printf("Currency gRPC server starting on port %s", port)
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
			"service": "currency-service",
			"status":  "healthy",
			"version": "1.0.0",
			"time":    time.Now().Unix(),
		})
	})

	// Currency endpoints (basic implementation)
	api := r.Group("/api/currency")
	{
		api.GET("/list", func(c *gin.Context) {
			currencies := []string{"IDR", "USD", "CNY", "EUR", "JPY", "SGD", "MYR"}
			c.JSON(http.StatusOK, gin.H{
				"currencies": currencies,
				"count":      len(currencies),
			})
		})

		api.GET("/rate", func(c *gin.Context) {
			from := c.Query("from")
			to := c.Query("to")
			
			if from == "" || to == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "from and to parameters are required",
				})
				return
			}

			// Mock exchange rate (in real implementation, this would fetch from database or external API)
			rate := 15000.0 // USD to IDR
			c.JSON(http.StatusOK, gin.H{
				"from_currency": from,
				"to_currency":   to,
				"rate":          rate,
				"timestamp":     time.Now().Unix(),
			})
		})

		api.POST("/convert", func(c *gin.Context) {
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

			// Mock conversion (in real implementation, this would use actual exchange rates)
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
		})
	}

	// Start server
	log.Printf("Currency HTTP server starting on port %s", port)
	
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
		log.Println("Currency service shutting down...")
		srv.Close()
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
} 