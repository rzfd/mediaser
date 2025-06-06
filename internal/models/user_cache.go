package models

import "time"

// UserCache represents cached user data in donation service
// This is a denormalized copy of essential user data for performance
type UserCache struct {
	Base
	UserID      uint      `json:"user_id" gorm:"uniqueIndex;not null"`
	Username    string    `json:"username" gorm:"not null"`
	FullName    string    `json:"full_name"`
	IsStreamer  bool      `json:"is_streamer" gorm:"default:false"`
	ProfilePic  string    `json:"profile_pic"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	LastSyncAt  time.Time `json:"last_sync_at" gorm:"default:CURRENT_TIMESTAMP"`
	ExpiresAt   time.Time `json:"expires_at"` // Cache expiration
}

// TableName specifies the table name for UserCache
func (UserCache) TableName() string {
	return "user_cache"
}

// IsExpired checks if the cached user data has expired
func (uc *UserCache) IsExpired() bool {
	return time.Now().After(uc.ExpiresAt)
}

// ShouldRefresh checks if the cache should be refreshed (every 1 hour)
func (uc *UserCache) ShouldRefresh() bool {
	return time.Since(uc.LastSyncAt) > time.Hour
}

// ToUser converts UserCache to User model
func (uc *UserCache) ToUser() *User {
	return &User{
		Base: Base{
			ID:        uc.UserID,
			CreatedAt: uc.CreatedAt,
			UpdatedAt: uc.UpdatedAt,
		},
		Username:   uc.Username,
		FullName:   uc.FullName,
		IsStreamer: uc.IsStreamer,
		ProfilePic: uc.ProfilePic,
	}
}

// FromUser creates UserCache from User model
func NewUserCacheFromUser(user *User, cacheDuration time.Duration) *UserCache {
	now := time.Now()
	return &UserCache{
		UserID:     user.ID,
		Username:   user.Username,
		FullName:   user.FullName,
		IsStreamer: user.IsStreamer,
		ProfilePic: user.ProfilePic,
		IsActive:   true,
		LastSyncAt: now,
		ExpiresAt:  now.Add(cacheDuration),
	}
} 