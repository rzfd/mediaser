package repository

import (
	"github.com/rzfd/mediashar/internal/models"
)

type UserCacheRepository interface {
	Get(userID uint) (*models.UserCache, error)
	Set(userCache *models.UserCache) error
	Delete(userID uint) error
	GetExpiredCaches() ([]*models.UserCache, error)
	CleanupExpired() error
	GetBatch(userIDs []uint) ([]*models.UserCache, error)
	UpdateLastSync(userID uint) error
} 