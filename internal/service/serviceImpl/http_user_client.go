package serviceImpl

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/service"
)

type httpUserServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPUserServiceClient(baseURL string) service.UserServiceClient {
	return &httpUserServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *httpUserServiceClient) GetUser(ctx context.Context, userID uint) (*models.User, error) {
	url := fmt.Sprintf("%s/api/users/%d", c.baseURL, userID)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("user %d not found", userID)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user %d: status %d", userID, resp.StatusCode)
	}

	var response struct {
		Status string       `json:"status"`
		Data   models.User  `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return &response.Data, nil
}

func (c *httpUserServiceClient) ValidateUser(ctx context.Context, userID uint, mustBeStreamer bool) error {
	user, err := c.GetUser(ctx, userID)
	if err != nil {
		return err
	}

	if mustBeStreamer && !user.IsStreamer {
		return fmt.Errorf("user %d is not a streamer", userID)
	}

	return nil
}

func (c *httpUserServiceClient) GetUsers(ctx context.Context, userIDs []uint) ([]*models.User, error) {
	// For batch requests, we could implement a POST endpoint
	// For now, we'll make multiple individual requests
	var users []*models.User
	
	for _, userID := range userIDs {
		user, err := c.GetUser(ctx, userID)
		if err != nil {
			// Log error but continue with other users
			fmt.Printf("Warning: Could not fetch user %d: %v\n", userID, err)
			continue
		}
		users = append(users, user)
	}

	return users, nil
} 