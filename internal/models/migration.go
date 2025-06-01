package models

import (
	"gorm.io/gorm"
)

// MigrateDB performs database migrations
func MigrateDB(db *gorm.DB) error {
	// Auto migrate the schemas in dependency order
	// Base tables first
	err := db.AutoMigrate(&User{}, &Donation{})
	if err != nil {
		return err
	}
	
	// Platform tables (depends on User)
	err = db.AutoMigrate(&StreamingPlatform{})
	if err != nil {
		return err
	}
	
	// Content tables (depends on StreamingPlatform)
	err = db.AutoMigrate(&StreamingContent{})
	if err != nil {
		return err
	}
	
	// Junction tables (depends on Donation and StreamingContent)
	return db.AutoMigrate(&ContentDonation{})
} 