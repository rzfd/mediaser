package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rzfd/mediashar/internal/service"
	"github.com/rzfd/mediashar/pkg/utils"
)

type QRISHandler struct {
	qrisService     service.QRISService
	donationService service.DonationService
}

func NewQRISHandler(qrisService service.QRISService, donationService service.DonationService) *QRISHandler {
	return &QRISHandler{
		qrisService:     qrisService,
		donationService: donationService,
	}
}

// GenerateQRIS generates QRIS QR code for a donation
func (h *QRISHandler) GenerateQRIS(c echo.Context) error {
	donationID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid donation ID", err))
	}

	// Get donation details
	donation, err := h.donationService.GetByID(uint(donationID))
	if err != nil {
		return c.JSON(http.StatusNotFound, utils.ErrorResponse("Donation not found", err))
	}

	// Check if donation is still pending
	if donation.Status != "pending" {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Donation is not in pending status", nil))
	}

	// Generate QRIS
	qrisResponse, err := h.qrisService.GenerateQRIS(donation)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to generate QRIS", err))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("QRIS generated successfully", qrisResponse))
}

// CheckQRISStatus checks the payment status of a QRIS transaction
func (h *QRISHandler) CheckQRISStatus(c echo.Context) error {
	transactionID := c.Param("transaction_id")
	if transactionID == "" {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Transaction ID is required", nil))
	}

	// Check payment status
	status, err := h.qrisService.ValidateQRISPayment(transactionID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to check payment status", err))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Payment status retrieved", status))
}

// QRISCallback handles payment notification from QRIS provider
func (h *QRISHandler) QRISCallback(c echo.Context) error {
	// Get raw body
	body := c.Request().Body
	defer body.Close()

	// Read body content
	payload := make([]byte, c.Request().ContentLength)
	_, err := body.Read(payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to read callback payload", err))
	}

	// Process callback
	if err := h.qrisService.ProcessQRISCallback(payload); err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to process callback", err))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Callback processed successfully", nil))
}

// CreateQRISDonation creates a donation and immediately generates QRIS
func (h *QRISHandler) CreateQRISDonation(c echo.Context) error {
	var req struct {
		Amount      float64 `json:"amount" validate:"required,min=1000"`
		Currency    string  `json:"currency" validate:"required"`
		Message     string  `json:"message"`
		StreamerID  uint    `json:"streamer_id" validate:"required"`
		DisplayName string  `json:"display_name"`
		IsAnonymous bool    `json:"is_anonymous"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", err))
	}

	// Get donator ID from JWT (optional for anonymous donations)
	var donatorID *uint
	if userID, ok := c.Get("user_id").(uint); ok && !req.IsAnonymous {
		donatorID = &userID
	}

	// Create donation
	donation, err := h.donationService.CreateDonation(&service.CreateDonationRequest{
		Amount:      req.Amount,
		Currency:    req.Currency,
		Message:     req.Message,
		StreamerID:  req.StreamerID,
		DonatorID:   donatorID,
		DisplayName: req.DisplayName,
		IsAnonymous: req.IsAnonymous,
	})
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to create donation", err))
	}

	// Generate QRIS immediately
	qrisResponse, err := h.qrisService.GenerateQRIS(donation)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to generate QRIS", err))
	}

	// Return both donation and QRIS data
	response := map[string]interface{}{
		"donation": donation,
		"qris":     qrisResponse,
	}

	return c.JSON(http.StatusCreated, utils.SuccessResponse("Donation created with QRIS", response))
} 