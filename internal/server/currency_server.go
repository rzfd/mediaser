package server

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"github.com/rzfd/mediashar/internal/handler"
	"github.com/rzfd/mediashar/internal/utils"
	"github.com/rzfd/mediashar/pkg/logger"
)

type CurrencyServer struct {
	grpcServer *grpc.Server
	httpServer *http.Server
	handler    *handler.CurrencyHTTPHandler
	grpcPort   string
	httpPort   string
}

func NewCurrencyServer() (*CurrencyServer, error) {
	handler := handler.NewCurrencyHTTPHandler()
	
	return &CurrencyServer{
		handler:  handler,
		grpcPort: utils.GetEnv("GRPC_PORT", "8084"),
		httpPort: utils.GetEnv("HTTP_PORT", "8094"),
	}, nil
}

func (s *CurrencyServer) StartGRPC() error {
	lis, err := net.Listen("tcp", ":"+s.grpcPort)
	if err != nil {
		return fmt.Errorf("failed to listen on gRPC port %s: %w", s.grpcPort, err)
	}

	s.grpcServer = grpc.NewServer()
	// Register currency service here when protobuf definitions are ready

	appLogger := logger.GetLogger()
	appLogger.Info("Currency gRPC server starting", "port", s.grpcPort)
	
	return s.grpcServer.Serve(lis)
}

func (s *CurrencyServer) StartHTTP() error {
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

	// Currency endpoints
	api := r.Group("/api/currency")
	{
		api.GET("/list", s.handler.ListCurrencies)
		api.GET("/rate", s.handler.GetExchangeRate)
		api.POST("/convert", s.handler.ConvertCurrency)
	}

	s.httpServer = &http.Server{
		Addr:    ":" + s.httpPort,
		Handler: r,
	}

	appLogger := logger.GetLogger()
	appLogger.Info("Currency HTTP server starting", "port", s.httpPort)
	
	return s.httpServer.ListenAndServe()
}

func (s *CurrencyServer) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
	if s.httpServer != nil {
		s.httpServer.Close()
	}
}

func (s *CurrencyServer) GetGRPCPort() string {
	return s.grpcPort
}

func (s *CurrencyServer) GetHTTPPort() string {
	return s.httpPort
} 