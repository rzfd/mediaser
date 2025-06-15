package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/service"
	"github.com/rzfd/mediashar/pkg/metrics"
	"github.com/rzfd/mediashar/pkg/utils"
	"google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

type AuthHandler struct {
	userService service.UserService
	authService service.AuthService
}

func NewAuthHandler(userService service.UserService, authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		authService: authService,
	}
}

// Register creates a new user account
func (h *AuthHandler) Register(c echo.Context) error {
	var req struct {
		Username    string `json:"username" validate:"required"`
		Email       string `json:"email" validate:"required,email"`
		Password    string `json:"password" validate:"required,min=6"`
		FullName    string `json:"full_name"`
		IsStreamer  bool   `json:"is_streamer"`
		Description string `json:"description"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err))
	}

	// Create user
	user := &models.User{
		Username:    req.Username,
		Email:       req.Email,
		Password:    req.Password,
		FullName:    req.FullName,
		IsStreamer:  req.IsStreamer,
		Description: req.Description,
	}

	if err := h.userService.Create(user); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to create user", err))
	}

	// Generate JWT token
	token, err := h.authService.GenerateToken(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to generate token", err))
	}

	// Don't return password
	user.Password = ""

	return c.JSON(http.StatusCreated, utils.SuccessResponse("User registered successfully", map[string]interface{}{
		"user":  user,
		"token": token,
	}))
}

// Login authenticates a user and returns a JWT token
func (h *AuthHandler) Login(c echo.Context) error {
	var req struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		// Record failed login due to bad request
		metrics.GetMetrics().RecordUserLogin("api-gateway", "email", "failed_bad_request")
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err))
	}

	// Authenticate user
	user, err := h.userService.Authenticate(req.Email, req.Password)
	if err != nil {
		// Record failed login due to invalid credentials
		metrics.GetMetrics().RecordUserLogin("api-gateway", "email", "failed_invalid_credentials")
		// Record failed login activity in database (use dummy user ID 0 for failed attempts)
		h.recordFailedLoginActivity(req.Email, "login_invalid_credentials")
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid credentials", err))
	}

	// Generate JWT token
	token, err := h.authService.GenerateToken(user)
	if err != nil {
		// Record failed login due to token generation error
		metrics.GetMetrics().RecordUserLogin("api-gateway", "email", "failed_token_error")
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to generate token", err))
	}

	// Record successful login
	metrics.GetMetrics().RecordUserLogin("api-gateway", "email", "success")
	
	// Record login activity in database
	h.recordUserActivity(user.ID, "login_success")

	// Don't return password
	user.Password = ""

	return c.JSON(http.StatusOK, utils.SuccessResponse("Login successful", map[string]interface{}{
		"user":  user,
		"token": token,
	}))
}

// RefreshToken generates a new token from an existing valid token
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req struct {
		Token string `json:"token" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err))
	}

	// Generate new token
	newToken, err := h.authService.RefreshToken(req.Token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid or expired token", err))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Token refreshed successfully", map[string]interface{}{
		"token": newToken,
	}))
}

// GetProfile returns the current user's profile
func (h *AuthHandler) GetProfile(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated", nil))
	}

	user, err := h.userService.GetByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("User not found", err))
	}

	// Don't return password
	user.Password = ""

	return c.JSON(http.StatusOK, utils.SuccessResponse("Profile retrieved successfully", user))
}

// UpdateProfile updates the current user's profile
func (h *AuthHandler) UpdateProfile(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated", nil))
	}

	var req struct {
		FullName    string `json:"full_name"`
		Description string `json:"description"`
		ProfilePic  string `json:"profile_pic"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err))
	}

	// Get current user
	user, err := h.userService.GetByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("User not found", err))
	}

	// Update fields
	user.FullName = req.FullName
	user.Description = req.Description
	user.ProfilePic = req.ProfilePic

	if err := h.userService.Update(user); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update profile", err))
	}

	// Don't return password
	user.Password = ""

	return c.JSON(http.StatusOK, utils.SuccessResponse("Profile updated successfully", user))
}

// ChangePassword changes the current user's password
func (h *AuthHandler) ChangePassword(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User not authenticated", nil))
	}

	var req struct {
		CurrentPassword string `json:"current_password" validate:"required"`
		NewPassword     string `json:"new_password" validate:"required,min=6"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err))
	}

	// Get current user
	user, err := h.userService.GetByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("User not found", err))
	}

	// Verify current password
	_, err = h.userService.Authenticate(user.Email, req.CurrentPassword)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Current password is incorrect", err))
	}

	// Update password
	user.Password = req.NewPassword
	if err := h.userService.Update(user); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update password", err))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Password changed successfully", nil))
}

// Logout (client-side token invalidation)
func (h *AuthHandler) Logout(c echo.Context) error {
	// In a stateless JWT system, logout is typically handled client-side
	// by removing the token from storage. However, we can provide this endpoint
	// for consistency and future token blacklisting implementation.
	
	return c.JSON(http.StatusOK, utils.SuccessResponse("Logged out successfully", nil))
}

// GoogleLogin handles Google OAuth login
func (h *AuthHandler) GoogleLogin(c echo.Context) error {
	var req struct {
		Credential string `json:"credential" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err))
	}

	// Verify Google JWT token
	userInfo, err := h.verifyGoogleToken(req.Credential)
	if err != nil {
		// Record failed Google login
		metrics.GetMetrics().RecordUserLogin("api-gateway", "google", "failed_invalid_token")
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid Google token", err))
	}

	// Check if user exists
	user, err := h.userService.GetByEmail(userInfo.Email)
	if err != nil {
		// User doesn't exist, create new user
		user = &models.User{
			Username:   userInfo.Email, // Use email as username initially
			Email:      userInfo.Email,
			FullName:   userInfo.Name,
			IsStreamer: false, // Default to donator
			Password:   "", // No password for OAuth users
		}

		if err := h.userService.Create(user); err != nil {
			return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create user", err))
		}
	}

	// Generate JWT token
	token, err := h.authService.GenerateToken(user)
	if err != nil {
		// Record failed Google login due to token error
		metrics.GetMetrics().RecordUserLogin("api-gateway", "google", "failed_token_error")
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to generate token", err))
	}

	// Record successful Google login
	metrics.GetMetrics().RecordUserLogin("api-gateway", "google", "success")
	
	// Record Google login activity in database
	h.recordUserActivity(user.ID, "login_google_success")

	// Don't return password
	user.Password = ""

	return c.JSON(http.StatusOK, utils.SuccessResponse("Google login successful", map[string]interface{}{
		"user":  user,
		"token": token,
	}))
}

// GoogleUserInfo represents user info from Google
type GoogleUserInfo struct {
	Email string
	Name  string
	ID    string
}

// verifyGoogleToken verifies the Google JWT token and returns user info
func (h *AuthHandler) verifyGoogleToken(credential string) (*GoogleUserInfo, error) {
	ctx := context.Background()
	
	// Create OAuth2 service
	oauth2Service, err := oauth2.NewService(ctx, option.WithoutAuthentication())
	if err != nil {
		return nil, err
	}

	// Verify the token
	tokenInfo, err := oauth2Service.Tokeninfo().IdToken(credential).Do()
	if err != nil {
		return nil, err
	}

	// Get user info from tokeninfo
	userInfo := &GoogleUserInfo{
		Email: tokenInfo.Email,
		Name:  tokenInfo.Email, // Use email as name if full name not available
		ID:    tokenInfo.UserId,
	}

	return userInfo, nil
}

// recordUserActivity records user activity in the database
func (h *AuthHandler) recordUserActivity(userID uint, activityType string) {
	// Get database connection from user service
	db := h.userService.GetDB()
	if db == nil {
		return // Skip if no database connection
	}
	
	// Insert activity record
	query := `INSERT INTO user_activities (user_id, activity_type, created_at) VALUES (?, ?, NOW())`
	err := db.Exec(query, userID, activityType).Error
	if err != nil {
		// Log error but don't fail the request
		log.Printf("Failed to record user activity: %v", err)
	}
}

// recordFailedLoginActivity records failed login activity in the database
func (h *AuthHandler) recordFailedLoginActivity(email string, activityType string) {
	// Get database connection from user service
	db := h.userService.GetDB()
	if db == nil {
		return // Skip if no database connection
	}
	
	// Try to find user by email to get proper user_id
	user, err := h.userService.GetByEmail(email)
	if err != nil {
		// If user not found, skip recording (can't record without valid user_id)
		log.Printf("Cannot record failed login activity for non-existent user: %s", email)
		return
	}
	
	// Insert activity record with valid user_id
	query := `INSERT INTO user_activities (user_id, activity_type, created_at) VALUES (?, ?, NOW())`
	err = db.Exec(query, user.ID, activityType).Error
	if err != nil {
		// Log error but don't fail the request
		log.Printf("Failed to record failed login activity: %v", err)
	}
} 