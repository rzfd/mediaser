package service

import (
	"github.com/rzfd/mediashar/internal/models"
	"gorm.io/gorm"
)

type UserService interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	List(page, pageSize int) ([]*models.User, error)
	GetStreamers(page, pageSize int) ([]*models.User, error)
	Authenticate(email, password string) (*models.User, error)
	GetDB() *gorm.DB
} 