package adapter

import (
	"context"
	"time"

	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/service"
	"github.com/rzfd/mediashar/pkg/pb"
)

type DonationServiceAdapter struct {
	donationClient pb.DonationServiceClient
}

func NewDonationServiceAdapter(donationClient pb.DonationServiceClient) *DonationServiceAdapter {
	return &DonationServiceAdapter{
		donationClient: donationClient,
	}
}

func (d *DonationServiceAdapter) Create(donation *models.Donation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcReq := &pb.CreateDonationRequest{
		Amount:        donation.Amount,
		Currency:      string(donation.Currency),
		Message:       donation.Message,
		StreamerId:    uint32(donation.StreamerID),
		DisplayName:   donation.DisplayName,
		IsAnonymous:   donation.IsAnonymous,
		PaymentMethod: "qris",
	}

	if donation.DonatorID != 0 {
		grpcReq.DonatorId = uint32(donation.DonatorID)
	}

	resp, err := d.donationClient.CreateDonation(ctx, grpcReq)
	if err != nil {
		return err
	}

	donation.ID = uint(resp.DonationId)
	return nil
}

func (d *DonationServiceAdapter) CreateDonation(req *service.CreateDonationRequest) (*models.Donation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcReq := &pb.CreateDonationRequest{
		Amount:        req.Amount,
		Currency:      req.Currency,
		Message:       req.Message,
		StreamerId:    uint32(req.StreamerID),
		DisplayName:   req.DisplayName,
		IsAnonymous:   req.IsAnonymous,
		PaymentMethod: "qris",
	}

	if req.DonatorID != nil {
		grpcReq.DonatorId = uint32(*req.DonatorID)
	}

	resp, err := d.donationClient.CreateDonation(ctx, grpcReq)
	if err != nil {
		return nil, err
	}

	donation := &models.Donation{
		Amount:      req.Amount,
		Currency:    models.SupportedCurrency(req.Currency),
		Message:     req.Message,
		StreamerID:  req.StreamerID,
		DisplayName: req.DisplayName,
		IsAnonymous: req.IsAnonymous,
		Status:      models.PaymentPending,
	}
	donation.ID = uint(resp.DonationId)

	if req.DonatorID != nil {
		donation.DonatorID = *req.DonatorID
	}

	return donation, nil
}

func (d *DonationServiceAdapter) GetByID(id uint) (*models.Donation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := d.donationClient.GetDonation(ctx, &pb.GetDonationRequest{
		DonationId: uint32(id),
	})
	if err != nil {
		return nil, err
	}

	donation := &models.Donation{
		Amount:      resp.Donation.Amount,
		Currency:    models.SupportedCurrency(resp.Donation.Currency),
		Message:     resp.Donation.Message,
		StreamerID:  uint(resp.Donation.StreamerId),
		DonatorID:   uint(resp.Donation.DonatorId),
		DisplayName: resp.Donation.DisplayName,
		IsAnonymous: resp.Donation.IsAnonymous,
	}
	donation.ID = uint(resp.Donation.Id)

	return donation, nil
}

func (d *DonationServiceAdapter) GetByTransactionID(transactionID string) (*models.Donation, error) {
	return &models.Donation{}, nil
}

func (d *DonationServiceAdapter) List(page, pageSize int) ([]*models.Donation, error) {
	return []*models.Donation{}, nil
}

func (d *DonationServiceAdapter) GetByDonatorID(donatorID uint, page, pageSize int) ([]*models.Donation, error) {
	return []*models.Donation{}, nil
}

func (d *DonationServiceAdapter) GetByStreamerID(streamerID uint, page, pageSize int) ([]*models.Donation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := d.donationClient.GetDonationsByStreamer(ctx, &pb.GetDonationsByStreamerRequest{
		StreamerId: uint32(streamerID),
		Page:       int32(page),
		PageSize:   int32(pageSize),
	})
	if err != nil {
		return nil, err
	}

	var donations []*models.Donation
	for _, pbDonation := range resp.Donations {
		donation := &models.Donation{
			Amount:      pbDonation.Amount,
			Currency:    models.SupportedCurrency(pbDonation.Currency),
			Message:     pbDonation.Message,
			StreamerID:  uint(pbDonation.StreamerId),
			DonatorID:   uint(pbDonation.DonatorId),
			DisplayName: pbDonation.DisplayName,
			IsAnonymous: pbDonation.IsAnonymous,
		}
		donation.ID = uint(pbDonation.Id)
		donations = append(donations, donation)
	}

	return donations, nil
}

func (d *DonationServiceAdapter) UpdateStatus(id uint, status models.PaymentStatus) error {
	return nil
}

func (d *DonationServiceAdapter) ProcessPayment(donationID uint, transactionID string, provider models.PaymentProvider) error {
	return nil
}

func (d *DonationServiceAdapter) GetLatestDonations(limit int) ([]*models.Donation, error) {
	return []*models.Donation{}, nil
}

func (d *DonationServiceAdapter) GetTotalAmountByStreamer(streamerID uint) (float64, error) {
	return 0, nil
} 