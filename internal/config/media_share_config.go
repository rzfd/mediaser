package config

import (
	"fmt"
	"os"

	"github.com/rzfd/mediashar/pkg/logger"
)

// MediaShareConfig holds all configuration for the media share service
type MediaShareConfig struct {
	Database DatabaseConfig
	Server   ServerConfig
	Logger   logger.Config
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
}

// ServerConfig holds server configuration
type ServerConfig struct {
	GRPCPort string
}

// LoadMediaShareConfig loads configuration from environment variables
func LoadMediaShareConfig() *MediaShareConfig {
	return &MediaShareConfig{
		Database: DatabaseConfig{
			Host:     getEnv("MEDIA_SHARE_DB_HOST", "localhost"),
			Port:     getEnv("MEDIA_SHARE_DB_PORT", "5432"),
			Username: getEnv("MEDIA_SHARE_DB_USERNAME", "postgres"),
			Password: getEnv("MEDIA_SHARE_DB_PASSWORD", "password"),
			Name:     getEnv("MEDIA_SHARE_DB_NAME", "media_share_db"),
		},
		Server: ServerConfig{
			GRPCPort: getEnv("GRPC_PORT", "9094"),
		},
		Logger: logger.Config{
			Level:       getEnv("LOG_LEVEL", "info"),
			Output:      getEnv("LOG_OUTPUT", "stdout"),
			LogFile:     getEnv("LOG_FILE", "logs/media-share-service.log"),
			ServiceName: "media-share-service",
		},
	}
}

// getEnv gets environment variable with fallback to default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// DSN returns the database connection string
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		d.Host, d.Username, d.Password, d.Name, d.Port)
} 