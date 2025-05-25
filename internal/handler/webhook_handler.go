package handler

import (
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/service"
	"github.com/rzfd/mediashar/pkg/utils"
)

type WebhookHandler struct {
	paymentService service.PaymentService
}

func NewWebhookHandler(paymentService service.PaymentService) *WebhookHandler {
	return &WebhookHandler{paymentService: paymentService}
}

// HandlePaypalWebhook handles PayPal payment webhooks
func (h *WebhookHandler) HandlePaypalWebhook(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to read request body", err))
	}

	transactionID, err := h.paymentService.ProcessWebhook(body, models.PaymentProviderPaypal)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to process PayPal webhook", err))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("PayPal webhook processed successfully", map[string]string{
		"transaction_id": transactionID,
	}))
}

// HandleStripeWebhook handles Stripe payment webhooks
func (h *WebhookHandler) HandleStripeWebhook(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to read request body", err))
	}

	transactionID, err := h.paymentService.ProcessWebhook(body, models.PaymentProviderStripe)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to process Stripe webhook", err))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Stripe webhook processed successfully", map[string]string{
		"transaction_id": transactionID,
	}))
}

// HandleCryptoWebhook handles cryptocurrency payment webhooks
func (h *WebhookHandler) HandleCryptoWebhook(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to read request body", err))
	}

	transactionID, err := h.paymentService.ProcessWebhook(body, models.PaymentProviderCrypto)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to process crypto webhook", err))
	}

	return c.JSON(http.StatusOK, utils.SuccessResponse("Crypto webhook processed successfully", map[string]string{
		"transaction_id": transactionID,
	}))
} 