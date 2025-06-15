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

	c.metrics.ActiveUsers24h.Set(float64(activeUsers24h))
	c.metrics.ActiveUsers7d.Set(float64(activeUsers7d))
	c.metrics.ActiveUsers30d.Set(float64(activeUsers30d))

	// Collect current active/online users
	currentActiveUsers := c.getCurrentActiveUsers(ctx)
	onlineUsers := c.getOnlineUsers(ctx)

	c.metrics.ActiveUsersTotal.Set(float64(currentActiveUsers))
	c.metrics.OnlineUsersCurrent.Set(float64(onlineUsers))

	// Collect active sessions
	activeSessions := c.getActiveSessions(ctx)
	c.metrics.ActiveSessionsTotal.Set(float64(activeSessions))

	log.Printf("User metrics updated: total=%d, active_24h=%d, active_7d=%d, active_30d=%d, online=%d",
		totalUsers, activeUsers24h, activeUsers7d, activeUsers30d, onlineUsers)
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