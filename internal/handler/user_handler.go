package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/service"
	"github.com/rzfd/mediashar/pkg/utils"
)

type UserHandler struct {
	userService     service.UserService
	donationService service.DonationService
}

func NewUserHandler(userService service.UserService, donationService service.DonationService) *UserHandler {
	return &UserHandler{
		userService:     userService,
		donationService: donationService,
	}
}

// GetProfile gets the current user profile from JWT token
func (h *UserHandler) GetProfile(c echo.Context) error {
	userID, ok := c.Get("user_id").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid token", nil))
	}

	user, err := h.userService.GetByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("User not found", err))
	}

	// Don't return password
	user.Password = ""

	return c.JSON(http.StatusOK, utils.SuccessResponse("Profile retrieved successfully", user))
}

// CreateUser creates a new user (admin only or public registration)
func (h *UserHandler) CreateUser(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err))
	}

	if err := h.userService.Create(&user); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to create user", err))
	}

	// Don't return the password hash
	user.Password = ""

	return c.JSON(http.StatusCreated, utils.SuccessResponse("User created successfully", user))
}

// GetUser gets a user by ID
func (h *UserHandler) GetUser(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid user ID", err))
	}

	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("User not found", err))
	}

	// Don't return password
	user.Password = ""

	return c.JSON(http.StatusOK, utils.SuccessResponse("User found", user))
}

// UpdateUser updates a user's information (admin or self only)
func (h *UserHandler) UpdateUser(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid user ID", err))
	}

	// Check if user is updating their own profile or is admin
	userID, ok := c.Get("user_id").(uint)
	if ok && userID != uint(id) {
		// TODO: Add admin check here
		return c.JSON(http.StatusForbidden, utils.ErrorResponse("Access denied", nil))
	}

	var userData models.User
	if err := c.Bind(&userData); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err))
	}

	// Ensure the ID in the URL matches the user object
	userData.ID = uint(id)

	if err := h.userService.Update(&userData); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to update user", err))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("User updated successfully", nil))
}

// GetUserDonations gets all donations by a specific user (as donator)
func (h *UserHandler) GetUserDonations(c echo.Context) error {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid user ID", err))
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	if pageSize <= 0 {
		pageSize = 10
	}

	donations, err := h.donationService.GetByDonatorID(uint(userID), page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch user donations", err))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("User donations fetched successfully", donations))
}

// ListStreamers lists all streamer users
func (h *UserHandler) ListStreamers(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	if pageSize <= 0 {
		pageSize = 10
	}

	streamers, err := h.userService.GetStreamers(page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch streamers", err))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Streamers fetched successfully", streamers))
} 