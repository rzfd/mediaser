package repository

import (
	"github.com/rzfd/mediashar/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	List(offset, limit int) ([]*models.User, error)
	GetStreamers(offset, limit int) ([]*models.User, error)
	GetDB() *gorm.DB
} 