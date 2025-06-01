package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rzfd/mediashar/internal/service"
)

type MidtransHandler struct {
	midtransService service.MidtransService
	donationService service.DonationService
}

func NewMidtransHandler(midtransService service.MidtransService, donationService service.DonationService) *MidtransHandler {
	return &MidtransHandler{
		midtransService: midtransService,
		donationService: donationService,
	}
}

// CreatePayment creates a new Midtrans payment for donation
func (h *MidtransHandler) CreatePayment(c echo.Context) error {
	donationIDParam := c.Param("donationId")
	donationID, err := strconv.ParseUint(donationIDParam, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Invalid donation ID",
		})
	}

	// Get donation details
	donation, err := h.donationService.GetByID(uint(donationID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"status":  "error",
			"message": "Donation not found",
		})
	}

	// Create Midtrans payment
	response, err := h.midtransService.ProcessDonationPayment(donation)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Failed to create payment",
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   response,
	})
}

// HandleWebhook handles Midtrans payment notification webhook
func (h *MidtransHandler) HandleWebhook(c echo.Context) error {
	var notification service.MidtransNotification
	
	if err := c.Bind(&notification); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Invalid notification payload",
		})
	}

	// Process the notification
	err := h.midtransService.HandleNotification(&notification)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Failed to process notification",
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Notification processed successfully",
	})
}

// GetTransactionStatus gets the status of a Midtrans transaction
func (h *MidtransHandler) GetTransactionStatus(c echo.Context) error {
	orderID := c.Param("orderId")
	if orderID == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  "error",
			"message": "Order ID is required",
		})
	}

	status, err := h.midtransService.GetTransactionStatus(orderID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  "error",
			"message": "Failed to get transaction status",
			"details": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
		"data":   status,
	})
} 