package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds all Prometheus metrics
type Metrics struct {
	// HTTP metrics
	HTTPRequestsTotal     *prometheus.CounterVec
	HTTPRequestDuration   *prometheus.HistogramVec
	HTTPResponseSize      *prometheus.HistogramVec
	
	// gRPC metrics
	GRPCRequestsTotal     *prometheus.CounterVec
	GRPCRequestDuration   *prometheus.HistogramVec
	
	// Database metrics
	DBConnectionsActive   prometheus.Gauge
	DBQueryDuration       *prometheus.HistogramVec
	DBQueriesTotal        *prometheus.CounterVec
	
	// Business metrics
	DonationsTotal        *prometheus.CounterVec
	DonationAmount        *prometheus.HistogramVec
	PaymentsProcessed     *prometheus.CounterVec
	
	// User metrics
	UserRegistrationsTotal *prometheus.CounterVec
	TotalUsersRegistered   prometheus.Gauge
	ActiveUsersTotal       prometheus.Gauge
	OnlineUsersCurrent     prometheus.Gauge
	ActiveUsers24h         prometheus.Gauge
	ActiveUsers7d          prometheus.Gauge
	ActiveUsers30d         prometheus.Gauge
	UserLoginTotal         *prometheus.CounterVec
	UserLogoutTotal        *prometheus.CounterVec
	UserActivityTotal      *prometheus.CounterVec
	ActiveSessionsTotal    prometheus.Gauge
	SessionDurationSeconds *prometheus.HistogramVec
	
	// System metrics
	GoRoutines            prometheus.Gauge
	MemoryUsage          prometheus.Gauge
	CPUUsage             prometheus.Gauge
}

// NewMetrics creates a new metrics instance
func NewMetrics(serviceName string) *Metrics {
	m := &Metrics{
		// HTTP metrics
		HTTPRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"service", "method", "path", "status_code"},
		),
		HTTPRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"service", "method", "path"},
		),
		HTTPResponseSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_response_size_bytes",
				Help:    "Size of HTTP responses in bytes",
				Buckets: []float64{100, 1000, 10000, 100000, 1000000},
			},
			[]string{"service", "method", "path"},
		),
		
		// gRPC metrics
		GRPCRequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "grpc_requests_total",
				Help: "Total number of gRPC requests",
			},
			[]string{"service", "method", "status"},
		),
		GRPCRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "grpc_request_duration_seconds",
				Help:    "Duration of gRPC requests in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"service", "method"},
		),
		
		// Database metrics
		DBConnectionsActive: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "db_connections_active",
				Help: "Number of active database connections",
			},
		),
		DBQueryDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "db_query_duration_seconds",
				Help:    "Duration of database queries in seconds",
				Buckets: []float64{0.001, 0.01, 0.1, 1, 10},
			},
			[]string{"service", "query_type"},
		),
		DBQueriesTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "db_queries_total",
				Help: "Total number of database queries",
			},
			[]string{"service", "query_type", "status"},
		),
		
		// Business metrics
		DonationsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "donations_total",
				Help: "Total number of donations",
			},
			[]string{"service", "currency", "status"},
		),
		DonationAmount: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "donation_amount",
				Help:    "Amount of donations",
				Buckets: []float64{10, 50, 100, 500, 1000, 5000, 10000},
			},
			[]string{"service", "currency"},
		),
		PaymentsProcessed: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "payments_processed_total",
				Help: "Total number of payments processed",
			},
			[]string{"service", "provider", "status"},
		),
		
		// User metrics
		UserRegistrationsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "user_registrations_total",
				Help: "Total number of user registrations",
			},
			[]string{"service", "platform", "status"},
		),
		TotalUsersRegistered: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "total_users_registered",
				Help: "Total number of registered users",
			},
		),
		ActiveUsersTotal: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "active_users_total",
				Help: "Total number of active users",
			},
		),
		OnlineUsersCurrent: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "online_users_current",
				Help: "Current number of online users",
			},
		),
		ActiveUsers24h: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "active_users_24h",
				Help: "Number of active users in the last 24 hours",
			},
		),
		ActiveUsers7d: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "active_users_7d",
				Help: "Number of active users in the last 7 days",
			},
		),
		ActiveUsers30d: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "active_users_30d",
				Help: "Number of active users in the last 30 days",
			},
		),
		UserLoginTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "user_login_total",
				Help: "Total number of user logins",
			},
			[]string{"service", "method", "status"},
		),
		UserLogoutTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "user_logout_total",
				Help: "Total number of user logouts",
			},
			[]string{"service"},
		),
		UserActivityTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "user_activity_total",
				Help: "Total user activities",
			},
			[]string{"service", "activity_type"},
		),
		ActiveSessionsTotal: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "active_sessions_total",
				Help: "Total number of active user sessions",
			},
		),
		SessionDurationSeconds: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "session_duration_seconds",
				Help:    "Duration of user sessions in seconds",
				Buckets: []float64{60, 300, 900, 1800, 3600, 7200, 14400}, // 1min to 4hours
			},
			[]string{"service"},
		),
		
		// System metrics
		GoRoutines: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "go_routines",
				Help: "Number of goroutines",
			},
		),
		MemoryUsage: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "memory_usage_bytes",
				Help: "Memory usage in bytes",
			},
		),
		CPUUsage: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "cpu_usage_percent",
				Help: "CPU usage percentage",
			},
		),
	}

	// Register all metrics
	prometheus.MustRegister(
		m.HTTPRequestsTotal,
		m.HTTPRequestDuration,
		m.HTTPResponseSize,
		m.GRPCRequestsTotal,
		m.GRPCRequestDuration,
		m.DBConnectionsActive,
		m.DBQueryDuration,
		m.DBQueriesTotal,
		m.DonationsTotal,
		m.DonationAmount,
		m.PaymentsProcessed,
		m.UserRegistrationsTotal,
		m.TotalUsersRegistered,
		m.ActiveUsersTotal,
		m.OnlineUsersCurrent,
		m.ActiveUsers24h,
		m.ActiveUsers7d,
		m.ActiveUsers30d,
		m.UserLoginTotal,
		m.UserLogoutTotal,
		m.UserActivityTotal,
		m.ActiveSessionsTotal,
		m.SessionDurationSeconds,
		m.GoRoutines,
		m.MemoryUsage,
		m.CPUUsage,
	)

	return m
}

// RecordHTTPRequest records HTTP request metrics
func (m *Metrics) RecordHTTPRequest(serviceName, method, path string, statusCode int, duration time.Duration, responseSize int) {
	m.HTTPRequestsTotal.WithLabelValues(serviceName, method, path, strconv.Itoa(statusCode)).Inc()
	m.HTTPRequestDuration.WithLabelValues(serviceName, method, path).Observe(duration.Seconds())
	m.HTTPResponseSize.WithLabelValues(serviceName, method, path).Observe(float64(responseSize))
}

// RecordGRPCRequest records gRPC request metrics
func (m *Metrics) RecordGRPCRequest(serviceName, method, status string, duration time.Duration) {
	m.GRPCRequestsTotal.WithLabelValues(serviceName, method, status).Inc()
	m.GRPCRequestDuration.WithLabelValues(serviceName, method).Observe(duration.Seconds())
}

// RecordDBQuery records database query metrics
func (m *Metrics) RecordDBQuery(serviceName, queryType, status string, duration time.Duration) {
	m.DBQueriesTotal.WithLabelValues(serviceName, queryType, status).Inc()
	m.DBQueryDuration.WithLabelValues(serviceName, queryType).Observe(duration.Seconds())
}

// RecordDonation records donation metrics
func (m *Metrics) RecordDonation(serviceName, currency, status string, amount float64) {
	m.DonationsTotal.WithLabelValues(serviceName, currency, status).Inc()
	m.DonationAmount.WithLabelValues(serviceName, currency).Observe(amount)
}

// RecordPayment records payment metrics
func (m *Metrics) RecordPayment(serviceName, provider, status string) {
	m.PaymentsProcessed.WithLabelValues(serviceName, provider, status).Inc()
}

// RecordUserRegistration records user registration metrics
func (m *Metrics) RecordUserRegistration(serviceName, platform, status string) {
	m.UserRegistrationsTotal.WithLabelValues(serviceName, platform, status).Inc()
}

// RecordUserLogin records user login metrics
func (m *Metrics) RecordUserLogin(serviceName, method, status string) {
	m.UserLoginTotal.WithLabelValues(serviceName, method, status).Inc()
}

// RecordUserLogout records user logout metrics
func (m *Metrics) RecordUserLogout(serviceName string) {
	m.UserLogoutTotal.WithLabelValues(serviceName).Inc()
}

// RecordUserActivity records user activity metrics
func (m *Metrics) RecordUserActivity(serviceName, activityType string) {
	m.UserActivityTotal.WithLabelValues(serviceName, activityType).Inc()
}

// RecordSessionDuration records session duration
func (m *Metrics) RecordSessionDuration(serviceName string, duration float64) {
	m.SessionDurationSeconds.WithLabelValues(serviceName).Observe(duration)
}

// UpdateUserMetrics updates user count metrics
func (m *Metrics) UpdateUserMetrics(totalUsers, activeUsers, onlineUsers float64) {
	m.TotalUsersRegistered.Set(totalUsers)
	m.ActiveUsersTotal.Set(activeUsers)
	m.OnlineUsersCurrent.Set(onlineUsers)
}

// UpdateActiveUsersMetrics updates active users metrics for different periods
func (m *Metrics) UpdateActiveUsersMetrics(users24h, users7d, users30d, activeSessions float64) {
	m.ActiveUsers24h.Set(users24h)
	m.ActiveUsers7d.Set(users7d)
	m.ActiveUsers30d.Set(users30d)
	m.ActiveSessionsTotal.Set(activeSessions)
}

// UpdateSystemMetrics updates system metrics
func (m *Metrics) UpdateSystemMetrics(goroutines int, memoryUsage, cpuUsage float64) {
	m.GoRoutines.Set(float64(goroutines))
	m.MemoryUsage.Set(memoryUsage)
	m.CPUUsage.Set(cpuUsage)
}

// HTTPMiddleware returns HTTP middleware for metrics collection
func (m *Metrics) HTTPMiddleware(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Create response writer wrapper to capture status code and size
			wrapper := &responseWriter{ResponseWriter: w, statusCode: 200}
			
			next.ServeHTTP(wrapper, r)
			
			duration := time.Since(start)
			m.RecordHTTPRequest(serviceName, r.Method, r.URL.Path, wrapper.statusCode, duration, wrapper.size)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture metrics
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

// MetricsHandler returns HTTP handler for Prometheus metrics endpoint
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

// Global metrics instance
var GlobalMetrics *Metrics

// Init initializes global metrics
func Init(serviceName string) {
	GlobalMetrics = NewMetrics(serviceName)
}

// GetMetrics returns global metrics instance
func GetMetrics() *Metrics {
	return GlobalMetrics
} 