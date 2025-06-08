package models

import (
	"time"
	"gorm.io/gorm"
)

// MediaShareSettings represents streamer's media share configuration
type MediaShareSettings struct {
	gorm.Model
	StreamerID           uint    `gorm:"uniqueIndex" json:"streamer_id"`
	MediaShareEnabled    bool    `json:"media_share_enabled" gorm:"default:true"`
	MinDonationAmount    float64 `json:"min_donation_amount" gorm:"default:5000"`
	Currency             string  `json:"currency" gorm:"default:'IDR'"`
	AllowYoutube         bool    `json:"allow_youtube" gorm:"default:true"`
	AllowTiktok          bool    `json:"allow_tiktok" gorm:"default:true"`
	AutoApprove          bool    `json:"auto_approve" gorm:"default:false"`
	MaxDurationYoutube   int     `json:"max_duration_youtube" gorm:"default:300"` // seconds
	MaxDurationTiktok    int     `json:"max_duration_tiktok" gorm:"default:180"`  // seconds
	WelcomeMessage       string  `json:"welcome_message" gorm:"type:text"`
	
	// Relations
	Streamer *User `gorm:"foreignKey:StreamerID" json:"streamer,omitempty"`
}

// MediaShareStatus represents the status of a shared media
type MediaShareStatus string

const (
	MediaShareStatusPending  MediaShareStatus = "pending"
	MediaShareStatusApproved MediaShareStatus = "approved"
	MediaShareStatusRejected MediaShareStatus = "rejected"
)

// MediaShareType represents the type of media platform
type MediaShareType string

const (
	MediaShareTypeYoutube MediaShareType = "youtube"
	MediaShareTypeTiktok  MediaShareType = "tiktok"
)

// MediaShare represents a media shared by a donator
type MediaShare struct {
	gorm.Model
	DonationID       uint             `json:"donation_id" gorm:"index"`
	StreamerID       uint             `json:"streamer_id" gorm:"index"`
	DonatorID        uint             `json:"donator_id" gorm:"index"`
	Type             MediaShareType   `json:"type" gorm:"type:varchar(20)"`
	URL              string           `json:"url" gorm:"type:text"`
	Title            string           `json:"title" gorm:"type:varchar(255)"`
	Message          string           `json:"message" gorm:"type:text"`
	Status           MediaShareStatus `json:"status" gorm:"type:varchar(20);default:'pending'"`
	DonationAmount   float64          `json:"donation_amount"`
	Currency         string           `json:"currency" gorm:"default:'IDR'"`
	DonatorName      string           `json:"donator_name" gorm:"type:varchar(100)"`
	Thumbnail        string           `json:"thumbnail" gorm:"type:text"`
	Duration         int              `json:"duration"` // in seconds
	ProcessedAt      *time.Time       `json:"processed_at"`
	
	// Relations
	Donation *Donation `gorm:"foreignKey:DonationID" json:"donation,omitempty"`
	Streamer *User     `gorm:"foreignKey:StreamerID" json:"streamer,omitempty"`
	Donator  *User     `gorm:"foreignKey:DonatorID" json:"donator,omitempty"`
}

// MediaShareRequest represents the request to share media
type MediaShareRequest struct {
	Type           MediaShareType `json:"type" validate:"required,oneof=youtube tiktok"`
	URL            string         `json:"url" validate:"required,url"`
	Title          string         `json:"title" validate:"max=255"`
	Message        string         `json:"message" validate:"max=1000"`
	DonationAmount float64        `json:"donation_amount" validate:"required,min=0"`
}

// MediaShareResponse represents the response after sharing media
type MediaShareResponse struct {
	ID        uint             `json:"id"`
	Type      MediaShareType   `json:"type"`
	URL       string           `json:"url"`
	Title     string           `json:"title"`
	Message   string           `json:"message"`
	Status    MediaShareStatus `json:"status"`
	CreatedAt time.Time        `json:"created_at"`
}

// MediaQueueItem represents an item in the media queue for streamers
type MediaQueueItem struct {
	ID             uint             `json:"id"`
	Type           MediaShareType   `json:"type"`
	URL            string           `json:"url"`
	Title          string           `json:"title"`
	Message        string           `json:"message"`
	Status         MediaShareStatus `json:"status"`
	DonatorName    string           `json:"donator_name"`
	DonationAmount float64          `json:"donation_amount"`
	Currency       string           `json:"currency"`
	Thumbnail      string           `json:"thumbnail"`
	SubmittedAt    time.Time        `json:"submitted_at"`
	ProcessedAt    *time.Time       `json:"processed_at"`
} 