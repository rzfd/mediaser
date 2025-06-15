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

type LanguageServer struct {
	grpcServer *grpc.Server
	httpServer *http.Server
	handler    *handler.LanguageHTTPHandler
	grpcPort   string
	httpPort   string
}

func NewLanguageServer() (*LanguageServer, error) {
	handler := handler.NewLanguageHTTPHandler()
	
	return &LanguageServer{
		handler:  handler,
		grpcPort: utils.GetEnv("GRPC_PORT", "8085"),
		httpPort: utils.GetEnv("HTTP_PORT", "8095"),
	}, nil
}

func (s *LanguageServer) StartGRPC() error {
	lis, err := net.Listen("tcp", ":"+s.grpcPort)
	if err != nil {
		return fmt.Errorf("failed to listen on gRPC port %s: %w", s.grpcPort, err)
	}

	s.grpcServer = grpc.NewServer()
	// Register language service here when protobuf definitions are ready

	appLogger := logger.GetLogger()
	appLogger.Info("Language gRPC server starting", "port", s.grpcPort)
	
	return s.grpcServer.Serve(lis)
}

func (s *LanguageServer) StartHTTP() error {
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

	// Language endpoints
	api := r.Group("/api/language")
	{
		api.GET("/list", s.handler.ListLanguages)
		api.POST("/translate", s.handler.TranslateText)
		api.POST("/detect", s.handler.DetectLanguage)
		api.POST("/bulk-translate", s.handler.BulkTranslate)
	}

	s.httpServer = &http.Server{
		Addr:    ":" + s.httpPort,
		Handler: r,
	}

	appLogger := logger.GetLogger()
	appLogger.Info("Language HTTP server starting", "port", s.httpPort)
	
	return s.httpServer.ListenAndServe()
}

func (s *LanguageServer) Stop() {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
	if s.httpServer != nil {
		s.httpServer.Close()
	}
}

func (s *LanguageServer) GetGRPCPort() string {
	return s.grpcPort
}

func (s *LanguageServer) GetHTTPPort() string {
	return s.httpPort
} 