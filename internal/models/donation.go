package models

import "time"

// PaymentStatus represents the status of a donation payment
type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "pending"
	PaymentCompleted PaymentStatus = "completed"
	PaymentFailed    PaymentStatus = "failed"
	PaymentRefunded  PaymentStatus = "refunded"
)

// PaymentProvider represents the payment method used for donation
type PaymentProvider string

const (
	PaymentProviderPaypal   PaymentProvider = "paypal"
	PaymentProviderStripe   PaymentProvider = "stripe"
	PaymentProviderCrypto   PaymentProvider = "crypto"
	PaymentProviderMidtrans PaymentProvider = "midtrans"
)

// Donation represents a donation from a donator to a streamer
type Donation struct {
	Base
	Amount          float64         `json:"amount" gorm:"not null"`
	Currency        SupportedCurrency `json:"currency" gorm:"default:'IDR'"`
	Message         string          `json:"message"`
	DonatorID       uint            `json:"donator_id" gorm:"not null;index"` // No foreign key constraint for microservices
	Donator         User            `json:"donator,omitempty" gorm:"-"` // Excluded from database, handled at application level
	StreamerID      uint            `json:"streamer_id" gorm:"not null;index"` // No foreign key constraint for microservices
	Streamer        User            `json:"streamer,omitempty" gorm:"-"` // Excluded from database, handled at application level
	Status          PaymentStatus   `json:"status" gorm:"default:'pending'"`
	PaymentProvider PaymentProvider `json:"payment_provider"`
	TransactionID   string          `json:"transaction_id"`
	PaymentTime     *time.Time      `json:"payment_time"`
	DisplayName     string          `json:"display_name"` // Name to display (might be different from user name)
	IsAnonymous     bool            `json:"is_anonymous" gorm:"default:false"`
} 