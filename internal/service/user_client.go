package service

import (
	"context"
	"github.com/rzfd/mediashar/internal/models"
)

// UserServiceClient defines interface for calling User Service
type UserServiceClient interface {
	GetUser(ctx context.Context, userID uint) (*models.User, error)
	ValidateUser(ctx context.Context, userID uint, mustBeStreamer bool) error
	GetUsers(ctx context.Context, userIDs []uint) ([]*models.User, error)
}

// HTTPUserServiceClient implements UserServiceClient using HTTP REST calls
type HTTPUserServiceClient struct {
	baseURL string
}

// GRPCUserServiceClient implements UserServiceClient using gRPC calls
type GRPCUserServiceClient struct {
	grpcClient interface{} // Will use pb.UserServiceClient when available
}

func NewHTTPUserServiceClient(baseURL string) UserServiceClient {
	return &HTTPUserServiceClient{
		baseURL: baseURL,
	}
}

func (c *HTTPUserServiceClient) GetUser(ctx context.Context, userID uint) (*models.User, error) {
	// TODO: Implement HTTP call to User Service
	// GET /api/users/{userID}
	return nil, nil
}

func (c *HTTPUserServiceClient) ValidateUser(ctx context.Context, userID uint, mustBeStreamer bool) error {
	// TODO: Implement HTTP call to validate user
	// GET /api/users/{userID}/validate?streamer={mustBeStreamer}
	return nil
}

func (c *HTTPUserServiceClient) GetUsers(ctx context.Context, userIDs []uint) ([]*models.User, error) {
	// TODO: Implement batch HTTP call to get multiple users
	// POST /api/users/batch with JSON body containing user IDs
	return nil, nil
} 