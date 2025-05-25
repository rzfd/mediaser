package models

import (
	"gorm.io/gorm"
)

// MigrateDB performs database migrations
func MigrateDB(db *gorm.DB) error {
	// Auto migrate the schemas
	return db.AutoMigrate(
		&User{},
		&Donation{},
	)
} 