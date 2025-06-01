package service

import (
	"github.com/rzfd/mediashar/internal/models"
)

type MidtransPaymentRequest struct {
	OrderID      string  `json:"order_id"`
	Amount       float64 `json:"amount"`
	Currency     string  `json:"currency"`
	CustomerName string  `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
	Description  string  `json:"description"`
	CallbackURL  string  `json:"callback_url"`
}

type MidtransPaymentResponse struct {
	Token       string `json:"token"`
	RedirectURL string `json:"redirect_url"`
	OrderID     string `json:"order_id"`
}

type MidtransNotification struct {
	TransactionStatus string `json:"transaction_status"`
	StatusCode        string `json:"status_code"`
	TransactionID     string `json:"transaction_id"`
	OrderID           string `json:"order_id"`
	GrossAmount       string `json:"gross_amount"`
	PaymentType       string `json:"payment_type"`
	SignatureKey      string `json:"signature_key"`
}

type MidtransService interface {
	CreateSnapTransaction(req *MidtransPaymentRequest) (*MidtransPaymentResponse, error)
	HandleNotification(notification *MidtransNotification) error
	VerifySignature(notification *MidtransNotification) bool
	GetTransactionStatus(orderID string) (*MidtransNotification, error)
	ProcessDonationPayment(donation *models.Donation) (*MidtransPaymentResponse, error)
} 