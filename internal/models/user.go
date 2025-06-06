package models

// User represents a user in the system (streamer or donator)
type User struct {
	Base
	Username    string `json:"username" gorm:"unique;not null"`
	Email       string `json:"email" gorm:"unique;not null"`
	Password    string `json:"-" gorm:"not null"`
	FullName    string `json:"full_name"`
	IsStreamer  bool   `json:"is_streamer" gorm:"default:false"`
	ProfilePic  string `json:"profile_pic"`
	Description string `json:"description"`
	// Removed foreign key relationships to prevent constraints in microservices architecture
	// Donations   []Donation `json:"donations,omitempty" gorm:"foreignKey:DonatorID"`
	// Received    []Donation `json:"received,omitempty" gorm:"foreignKey:StreamerID"`
} 