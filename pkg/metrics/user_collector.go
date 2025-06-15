package metrics

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

// UserMetricsCollector collects user metrics from database
type UserMetricsCollector struct {
	db      *sql.DB
	metrics *Metrics
	ticker  *time.Ticker
	done    chan bool
}

// NewUserMetricsCollector creates a new user metrics collector
func NewUserMetricsCollector(db *sql.DB, metrics *Metrics) *UserMetricsCollector {
	return &UserMetricsCollector{
		db:      db,
		metrics: metrics,
		ticker:  time.NewTicker(30 * time.Second), // Update every 30 seconds
		done:    make(chan bool),
	}
}

// Start starts the metrics collection
func (c *UserMetricsCollector) Start() {
	log.Printf("Starting user metrics collector...")
	go func() {
		// Initial collection
		log.Printf("Running initial user metrics collection...")
		c.collectMetrics()

		for {
			select {
			case <-c.ticker.C:
				log.Printf("Running scheduled user metrics collection...")
				c.collectMetrics()
			case <-c.done:
				log.Printf("Stopping user metrics collector...")
				return
			}
		}
	}()
}

// Stop stops the metrics collection
func (c *UserMetricsCollector) Stop() {
	c.ticker.Stop()
	c.done <- true
}

// collectMetrics collects all user metrics from database
func (c *UserMetricsCollector) collectMetrics() {
	ctx := context.Background()

	// Collect total registered users
	totalUsers := c.getTotalUsers(ctx)
	c.metrics.TotalUsersRegistered.Set(float64(totalUsers))

	// Collect active users for different periods
	activeUsers24h := c.getActiveUsers(ctx, 24)
	activeUsers7d := c.getActiveUsers(ctx, 24*7)
	activeUsers30d := c.getActiveUsers(ctx, 24*30)
	
	// Collect today's active users (midnight to now)
	activeTodayUsers := c.getActiveTodayUsers(ctx)

	c.metrics.ActiveUsers24h.Set(float64(activeUsers24h))
	c.metrics.ActiveUsers7d.Set(float64(activeUsers7d))
	c.metrics.ActiveUsers30d.Set(float64(activeUsers30d))
	c.metrics.ActiveUsersToday.Set(float64(activeTodayUsers))

	// Collect current active/online users
	currentActiveUsers := c.getCurrentActiveUsers(ctx)
	onlineUsers := c.getOnlineUsers(ctx)

	c.metrics.ActiveUsersTotal.Set(float64(currentActiveUsers))
	c.metrics.OnlineUsersCurrent.Set(float64(onlineUsers))

	// Collect active sessions
	activeSessions := c.getActiveSessions(ctx)
	c.metrics.ActiveSessionsTotal.Set(float64(activeSessions))

	// Collect login rate metrics
	c.collectLoginRateMetrics(ctx)

	log.Printf("User metrics updated: total=%d, active_24h=%d, active_7d=%d, active_30d=%d, active_today=%d, online=%d",
		totalUsers, activeUsers24h, activeUsers7d, activeUsers30d, activeTodayUsers, onlineUsers)
}

// getTotalUsers gets total number of registered users
func (c *UserMetricsCollector) getTotalUsers(ctx context.Context) int {
	var count int
	query := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL`
	
	err := c.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		log.Printf("Error getting total users: %v", err)
		return 0
	}
	
	return count
}

// getActiveUsers gets active users within specified hours
func (c *UserMetricsCollector) getActiveUsers(ctx context.Context, hours int) int {
	var count int
	query := `
		SELECT COUNT(DISTINCT user_id) 
		FROM user_activities 
		WHERE created_at >= NOW() - INTERVAL '%d hours'
	`
	
	err := c.db.QueryRowContext(ctx, fmt.Sprintf(query, hours)).Scan(&count)
	if err != nil {
		log.Printf("Error getting active users for %d hours: %v", hours, err)
		return 0
	}
	
	return count
}

// getCurrentActiveUsers gets currently active users (logged in recently)
func (c *UserMetricsCollector) getCurrentActiveUsers(ctx context.Context) int {
	var count int
	query := `
		SELECT COUNT(DISTINCT user_id) 
		FROM user_sessions 
		WHERE is_active = true AND last_activity >= NOW() - INTERVAL '1 hour'
	`
	
	err := c.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		log.Printf("Error getting current active users: %v", err)
		return 0
	}
	
	return count
}

// getOnlineUsers gets currently online users
func (c *UserMetricsCollector) getOnlineUsers(ctx context.Context) int {
	var count int
	query := `
		SELECT COUNT(DISTINCT user_id) 
		FROM user_sessions 
		WHERE is_active = true AND last_activity >= NOW() - INTERVAL '5 minutes'
	`
	
	err := c.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		log.Printf("Error getting online users: %v", err)
		return 0
	}
	
	return count
}

// getActiveSessions gets number of active sessions
func (c *UserMetricsCollector) getActiveSessions(ctx context.Context) int {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM user_sessions 
		WHERE is_active = true AND expires_at > NOW()
	`
	
	err := c.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		log.Printf("Error getting active sessions: %v", err)
		return 0
	}
	
	return count
}

// getActiveTodayUsers gets active users today (midnight to now)
func (c *UserMetricsCollector) getActiveTodayUsers(ctx context.Context) int {
	var count int
	query := `
		SELECT COUNT(DISTINCT user_id) 
		FROM user_activities 
		WHERE DATE(created_at) = CURRENT_DATE
	`
	
	err := c.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		log.Printf("Error getting today's active users: %v", err)
		return 0
	}
	
	return count
}

// collectLoginRateMetrics collects login success/failure rates
func (c *UserMetricsCollector) collectLoginRateMetrics(ctx context.Context) {
	// Get login attempts in the last hour
	loginAttemptsLastHour := c.getLoginAttempts(ctx, 1)
	successfulLoginsLastHour := c.getSuccessfulLogins(ctx, 1)
	failedLoginsLastHour := c.getFailedLogins(ctx, 1)
	
	// Get login attempts in the last 24 hours
	loginAttemptsLast24h := c.getLoginAttempts(ctx, 24)
	successfulLoginsLast24h := c.getSuccessfulLogins(ctx, 24)
	failedLoginsLast24h := c.getFailedLogins(ctx, 24)
	
	// Calculate success rate (percentage)
	var successRateLastHour, successRateLast24h float64
	if loginAttemptsLastHour > 0 {
		successRateLastHour = (float64(successfulLoginsLastHour) / float64(loginAttemptsLastHour)) * 100
	}
	if loginAttemptsLast24h > 0 {
		successRateLast24h = (float64(successfulLoginsLast24h) / float64(loginAttemptsLast24h)) * 100
	}
	
	// Update metrics
	c.metrics.LoginAttemptsLastHour.Set(float64(loginAttemptsLastHour))
	c.metrics.LoginSuccessLastHour.Set(float64(successfulLoginsLastHour))
	c.metrics.LoginFailedLastHour.Set(float64(failedLoginsLastHour))
	c.metrics.LoginSuccessRateLastHour.Set(successRateLastHour)
	
	c.metrics.LoginAttemptsLast24h.Set(float64(loginAttemptsLast24h))
	c.metrics.LoginSuccessLast24h.Set(float64(successfulLoginsLast24h))
	c.metrics.LoginFailedLast24h.Set(float64(failedLoginsLast24h))
	c.metrics.LoginSuccessRateLast24h.Set(successRateLast24h)
	
	log.Printf("Login rate metrics: attempts_1h=%d, success_1h=%d, failed_1h=%d, success_rate_1h=%.2f%%, attempts_24h=%d, success_24h=%d, failed_24h=%d, success_rate_24h=%.2f%%",
		loginAttemptsLastHour, successfulLoginsLastHour, failedLoginsLastHour, successRateLastHour,
		loginAttemptsLast24h, successfulLoginsLast24h, failedLoginsLast24h, successRateLast24h)
}

// getLoginAttempts gets total login attempts in the last N hours
func (c *UserMetricsCollector) getLoginAttempts(ctx context.Context, hours int) int {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM user_activities 
		WHERE activity_type IN ('login', 'login_success', 'login_failed', 'login_google_success', 'login_google_failed', 'login_invalid_credentials')
		AND created_at >= NOW() - INTERVAL '%d hours'
	`
	
	err := c.db.QueryRowContext(ctx, fmt.Sprintf(query, hours)).Scan(&count)
	if err != nil {
		log.Printf("Error getting login attempts for %d hours: %v", hours, err)
		return 0
	}
	
	return count
}

// getSuccessfulLogins gets successful login attempts in the last N hours
func (c *UserMetricsCollector) getSuccessfulLogins(ctx context.Context, hours int) int {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM user_activities 
		WHERE activity_type IN ('login_success', 'login_google_success', 'login')
		AND created_at >= NOW() - INTERVAL '%d hours'
	`
	
	err := c.db.QueryRowContext(ctx, fmt.Sprintf(query, hours)).Scan(&count)
	if err != nil {
		log.Printf("Error getting successful logins for %d hours: %v", hours, err)
		return 0
	}
	
	return count
}

// getFailedLogins gets failed login attempts in the last N hours
func (c *UserMetricsCollector) getFailedLogins(ctx context.Context, hours int) int {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM user_activities 
		WHERE activity_type IN ('login_failed', 'login_google_failed', 'login_invalid_credentials')
		AND created_at >= NOW() - INTERVAL '%d hours'
	`
	
	err := c.db.QueryRowContext(ctx, fmt.Sprintf(query, hours)).Scan(&count)
	if err != nil {
		log.Printf("Error getting failed logins for %d hours: %v", hours, err)
		return 0
	}
	
	return count
}

// Global user metrics collector
var GlobalUserCollector *UserMetricsCollector

// InitUserMetricsCollector initializes the global user metrics collector
func InitUserMetricsCollector(db *sql.DB, metrics *Metrics) {
	GlobalUserCollector = NewUserMetricsCollector(db, metrics)
	GlobalUserCollector.Start()
}

// StopUserMetricsCollector stops the global user metrics collector
func StopUserMetricsCollector() {
	if GlobalUserCollector != nil {
		GlobalUserCollector.Stop()
	}
} 