package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/repository"
)

// UserAggregatorService combines local cache and external API calls
type UserAggregatorService interface {
	GetUser(ctx context.Context, userID uint) (*models.User, error)
	GetUsers(ctx context.Context, userIDs []uint) ([]*models.User, error)
	ValidateUser(ctx context.Context, userID uint, mustBeStreamer bool) error
	RefreshUserCache(ctx context.Context, userID uint) error
	SyncUserData(ctx context.Context, user *models.User) error
}

type userAggregatorService struct {
	cacheRepo   repository.UserCacheRepository
	userClient  UserServiceClient
	cacheDuration time.Duration
}

func NewUserAggregatorService(
	cacheRepo repository.UserCacheRepository, 
	userClient UserServiceClient,
) UserAggregatorService {
	return &userAggregatorService{
		cacheRepo:     cacheRepo,
		userClient:    userClient,
		cacheDuration: 24 * time.Hour, // Cache for 24 hours
	}
}

func (s *userAggregatorService) GetUser(ctx context.Context, userID uint) (*models.User, error) {
	// First, try to get from cache
	cachedUser, err := s.cacheRepo.Get(userID)
	if err == nil && !cachedUser.IsExpired() {
		fmt.Printf("Cache hit for user ID %d\n", userID)
		return cachedUser.ToUser(), nil
	}

	// Cache miss or expired, fetch from User Service
	fmt.Printf("Cache miss for user ID %d, fetching from User Service\n", userID)
	
	if s.userClient == nil {
		return nil, fmt.Errorf("user service client not available and no valid cache for user %d", userID)
	}

	user, err := s.userClient.GetUser(ctx, userID)
	if err != nil {
		// If we have expired cache, return it as fallback
		if cachedUser != nil {
			fmt.Printf("⚠️ Using expired cache as fallback for user ID %d\n", userID)
			return cachedUser.ToUser(), nil
		}
		return nil, fmt.Errorf("failed to fetch user %d: %v", userID, err)
	}

	// Update cache with fresh data
	s.SyncUserData(ctx, user)
	
	return user, nil
}

func (s *userAggregatorService) GetUsers(ctx context.Context, userIDs []uint) ([]*models.User, error) {
	// Get cached users
	cachedUsers, _ := s.cacheRepo.GetBatch(userIDs)
	userMap := make(map[uint]*models.User)
	var missingIDs []uint

	// Check which users are cached and valid
	for _, cached := range cachedUsers {
		if !cached.IsExpired() {
			userMap[cached.UserID] = cached.ToUser()
		} else {
			missingIDs = append(missingIDs, cached.UserID)
		}
	}

	// Find completely missing users
	for _, userID := range userIDs {
		if _, exists := userMap[userID]; !exists {
			missingIDs = append(missingIDs, userID)
		}
	}

	// Fetch missing users from User Service
	if len(missingIDs) > 0 && s.userClient != nil {
		fetchedUsers, err := s.userClient.GetUsers(ctx, missingIDs)
		if err == nil {
			for _, user := range fetchedUsers {
				userMap[user.ID] = user
				// Cache the fetched user
				s.SyncUserData(ctx, user)
			}
		}
	}

	// Build result in the same order as requested
	var result []*models.User
	for _, userID := range userIDs {
		if user, exists := userMap[userID]; exists {
			result = append(result, user)
		}
	}

	return result, nil
}

func (s *userAggregatorService) ValidateUser(ctx context.Context, userID uint, mustBeStreamer bool) error {
	user, err := s.GetUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("user %d not found: %v", userID, err)
	}

	if mustBeStreamer && !user.IsStreamer {
		return fmt.Errorf("user %d is not a streamer", userID)
	}

	return nil
}

func (s *userAggregatorService) RefreshUserCache(ctx context.Context, userID uint) error {
	if s.userClient == nil {
		return fmt.Errorf("user service client not available")
	}

	user, err := s.userClient.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	return s.SyncUserData(ctx, user)
}

func (s *userAggregatorService) SyncUserData(ctx context.Context, user *models.User) error {
	userCache := models.NewUserCacheFromUser(user, s.cacheDuration)
	return s.cacheRepo.Set(userCache)
} 