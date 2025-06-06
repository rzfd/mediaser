package middleware

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rzfd/mediashar/pkg/metrics"
)

// MetricsMiddleware creates a middleware that records metrics for each HTTP request
func MetricsMiddleware(serviceName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			
			// Process the request
			err := next(c)
			
			// Record metrics
			duration := time.Since(start)
			statusCode := c.Response().Status
			
			// Get metrics instance
			m := metrics.GetMetrics()
			if m != nil {
				m.RecordHTTPRequest(
					serviceName,
					c.Request().Method,
					c.Request().URL.Path,
					statusCode,
					duration,
					int(c.Response().Size),
				)
			}
			
			return err
		}
	}
}

// GRPCMetricsInterceptor creates a gRPC interceptor for metrics collection
func GRPCMetricsInterceptor(serviceName string) func(interface{}, error) {
	return func(resp interface{}, err error) {
		status := "success"
		if err != nil {
			status = "error"
		}
		
		m := metrics.GetMetrics()
		if m != nil {
			m.RecordGRPCRequest(serviceName, "unknown", status, 0)
		}
	}
} 