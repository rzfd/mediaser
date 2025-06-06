package repositoryImpl

import (
	"time"
	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/repository"
	"gorm.io/gorm"
)

type userCacheRepository struct {
	db *gorm.DB
}

func NewUserCacheRepository(db *gorm.DB) repository.UserCacheRepository {
	return &userCacheRepository{db: db}
}

func (r *userCacheRepository) Get(userID uint) (*models.UserCache, error) {
	var userCache models.UserCache
	err := r.db.Where("user_id = ? AND is_active = ?", userID, true).First(&userCache).Error
	if err != nil {
		return nil, err
	}
	return &userCache, nil
}

func (r *userCacheRepository) Set(userCache *models.UserCache) error {
	// Use UPSERT: insert or update if exists
	return r.db.Save(userCache).Error
}

func (r *userCacheRepository) Delete(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.UserCache{}).Error
}

func (r *userCacheRepository) GetExpiredCaches() ([]*models.UserCache, error) {
	var caches []*models.UserCache
	err := r.db.Where("expires_at < ? AND is_active = ?", time.Now(), true).Find(&caches).Error
	return caches, err
}

func (r *userCacheRepository) CleanupExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Update("is_active", false).Error
}

func (r *userCacheRepository) GetBatch(userIDs []uint) ([]*models.UserCache, error) {
	var caches []*models.UserCache
	err := r.db.Where("user_id IN ? AND is_active = ?", userIDs, true).Find(&caches).Error
	return caches, err
}

func (r *userCacheRepository) UpdateLastSync(userID uint) error {
	return r.db.Model(&models.UserCache{}).
		Where("user_id = ?", userID).
		Update("last_sync_at", time.Now()).Error
} 