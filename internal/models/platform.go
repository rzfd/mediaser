package models

import (
	"time"
)

// StreamingPlatform represents a connected social media platform
type StreamingPlatform struct {
	Base
	UserID           uint   `json:"user_id" gorm:"not null;index"`
	PlatformType     string `json:"platform_type" gorm:"type:varchar(20);not null;check:platform_type IN ('youtube','tiktok','twitch')"`
	PlatformUserID   string `json:"platform_user_id" gorm:"type:varchar(255);not null"`
	PlatformUsername string `json:"platform_username" gorm:"type:varchar(255);not null"`
	ChannelURL       string `json:"channel_url" gorm:"type:text;not null"`
	ChannelName      string `json:"channel_name" gorm:"type:varchar(255)"`
	ProfileImageURL  string `json:"profile_image_url" gorm:"type:text"`
	FollowerCount    int    `json:"follower_count" gorm:"default:0"`
	IsVerified       bool   `json:"is_verified" gorm:"default:false"`
	IsActive         bool   `json:"is_active" gorm:"default:true"`

	// Relationships
	User             User               `json:"user,omitempty" gorm:"foreignKey:UserID"`
	StreamingContent []StreamingContent `json:"streaming_content,omitempty" gorm:"foreignKey:PlatformID"`
}

// TableName specifies the table name for StreamingPlatform
func (StreamingPlatform) TableName() string {
	return "streaming_platforms"
}

// StreamingContent represents content from streaming platforms
type StreamingContent struct {
	Base
	PlatformID   uint      `json:"platform_id" gorm:"not null;index"`
	ContentType  string    `json:"content_type" gorm:"type:varchar(20);not null;check:content_type IN ('live','video','short')"`
	ContentID    string    `json:"content_id" gorm:"type:varchar(255);not null"`
	ContentURL   string    `json:"content_url" gorm:"type:text;not null"`
	Title        string    `json:"title" gorm:"type:varchar(500)"`
	Description  string    `json:"description" gorm:"type:text"`
	ThumbnailURL string    `json:"thumbnail_url" gorm:"type:text"`
	Duration     *int      `json:"duration"` // dalam detik, NULL untuk live stream
	ViewCount    int       `json:"view_count" gorm:"default:0"`
	LikeCount    int       `json:"like_count" gorm:"default:0"`
	IsLive       bool      `json:"is_live" gorm:"default:false"`
	StartedAt    *time.Time `json:"started_at"`
	EndedAt      *time.Time `json:"ended_at"`

	// Relationships
	Platform         StreamingPlatform `json:"platform,omitempty" gorm:"foreignKey:PlatformID"`
	ContentDonations []ContentDonation `json:"content_donations,omitempty" gorm:"foreignKey:ContentID"`
}

// TableName specifies the table name for StreamingContent
func (StreamingContent) TableName() string {
	return "streaming_content"
}

// ContentDonation represents donations linked to specific content
type ContentDonation struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	DonationID   uint      `json:"donation_id" gorm:"not null;index"`
	ContentID    *uint     `json:"content_id" gorm:"index"` // nullable, no foreign key constraint
	PlatformType string    `json:"platform_type" gorm:"type:varchar(20);not null"`
	ContentURL   string    `json:"content_url" gorm:"type:text;not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relationships (disabled foreign key constraints to avoid circular dependency)
	Donation Donation         `json:"donation,omitempty" gorm:"-"`
	Content  *StreamingContent `json:"content,omitempty" gorm:"-"`
}

// TableName specifies the table name for ContentDonation
func (ContentDonation) TableName() string {
	return "content_donations"
}

// PlatformValidationResult represents the result of URL validation
type PlatformValidationResult struct {
	IsValid     bool                   `json:"is_valid"`
	Platform    string                 `json:"platform"`
	ContentType string                 `json:"content_type"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// PlatformMetadata represents metadata extracted from platform URLs
type PlatformMetadata struct {
	Title       string `json:"title"`
	Creator     string `json:"creator"`
	Thumbnail   string `json:"thumbnail"`
	Duration    *int   `json:"duration"`
	IsLive      bool   `json:"is_live"`
	ViewCount   int    `json:"view_count"`
	LikeCount   int    `json:"like_count"`
	ChannelID   string `json:"channel_id"`
	VideoID     string `json:"video_id"`
	Username    string `json:"username"`
} 