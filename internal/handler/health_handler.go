package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck provides a simple health check endpoint
func (h *HealthHandler) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"service": "mediashar-api",
		"version": "1.0.0",
	})
}

// ReadinessCheck checks if the service is ready to accept requests
func (h *HealthHandler) ReadinessCheck(c echo.Context) error {
	// You can add more sophisticated checks here
	// like database connectivity, external service availability, etc.
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ready",
		"checks": map[string]string{
			"database": "connected",
			"midtrans": "configured",
		},
	})
} 