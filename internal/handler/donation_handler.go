package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/service"
	"github.com/rzfd/mediashar/pkg/utils"
)

type DonationHandler struct {
	donationService service.DonationService
}

func NewDonationHandler(donationService service.DonationService) *DonationHandler {
	return &DonationHandler{donationService: donationService}
}

// CreateDonation creates a new donation
func (h *DonationHandler) CreateDonation(c echo.Context) error {
	var req service.CreateDonationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err))
	}

	// Get donator ID from JWT token
	userID, ok := c.Get("user_id").(uint)
	if ok {
		req.DonatorID = &userID
	}

	// Create donation using service method
	donation, err := h.donationService.CreateDonation(&req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to create donation", err))
	}

	return c.JSON(http.StatusCreated, utils.SuccessResponse("Donation created successfully", donation))
}

// GetDonation gets a donation by ID
func (h *DonationHandler) GetDonation(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid donation ID", err))
	}

	donation, err := h.donationService.GetByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("Donation not found", err))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Donation found", donation))
}

// ListDonations lists all donations with pagination
func (h *DonationHandler) ListDonations(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	if pageSize <= 0 {
		pageSize = 10
	}

	donations, err := h.donationService.List(page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch donations", err))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Donations fetched successfully", donations))
}

// GetStreamerDonations gets all donations for a specific streamer
func (h *DonationHandler) GetStreamerDonations(c echo.Context) error {
	streamerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid streamer ID", err))
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page <= 0 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(c.QueryParam("pageSize"))
	if pageSize <= 0 {
		pageSize = 10
	}

	donations, err := h.donationService.GetByStreamerID(uint(streamerID), page, pageSize)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch streamer donations", err))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Streamer donations fetched successfully", donations))
}

// ProcessPayment updates a donation with payment information
func (h *DonationHandler) ProcessPayment(c echo.Context) error {
	var paymentData struct {
		DonationID     uint                 `json:"donation_id"`
		TransactionID  string               `json:"transaction_id"`
		PaymentProvider models.PaymentProvider `json:"payment_provider"`
	}

	if err := c.Bind(&paymentData); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err))
	}

	err := h.donationService.ProcessPayment(
		paymentData.DonationID,
		paymentData.TransactionID,
		paymentData.PaymentProvider,
	)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to process payment", err))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Payment processed successfully", nil))
}

// GetLatestDonations gets the latest donations for display
func (h *DonationHandler) GetLatestDonations(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 10
	}

	donations, err := h.donationService.GetLatestDonations(limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch latest donations", err))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Latest donations fetched successfully", donations))
}

// GetTotalDonations gets the total donation amount for a streamer
func (h *DonationHandler) GetTotalDonations(c echo.Context) error {
	streamerID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid streamer ID", err))
	}

	total, err := h.donationService.GetTotalAmountByStreamer(uint(streamerID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch total donations", err))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Total donations fetched successfully", map[string]float64{
		"total": total,
	}))
} 