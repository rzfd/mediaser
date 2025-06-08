package grpc

import (
	"context"

	"github.com/rzfd/mediashar/internal/adapter"
	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/internal/service/serviceImpl"
	"github.com/rzfd/mediashar/pkg/pb"
)

// MediaShareGRPCHandler implements the gRPC MediaShareService
type MediaShareGRPCHandler struct {
	pb.UnimplementedMediaShareServiceServer
	service   serviceImpl.MediaShareService
	converter *adapter.MediaShareConverter
}

// GetSettings retrieves streamer media share settings
func (s *MediaShareGRPCHandler) GetSettings(ctx context.Context, req *pb.GetSettingsRequest) (*pb.GetSettingsResponse, error) {
	settings, err := s.service.GetSettingsByStreamerID(uint(req.StreamerId))
	if err != nil {
		return nil, err
	}

	return &pb.GetSettingsResponse{
		Settings: s.converter.ToProtoSettings(settings),
	}, nil
}

// UpdateSettings updates streamer media share settings  
func (s *MediaShareGRPCHandler) UpdateSettings(ctx context.Context, req *pb.UpdateSettingsRequest) (*pb.UpdateSettingsResponse, error) {
	settings := s.converter.FromProtoSettings(req.Settings)
	settings.StreamerID = uint(req.StreamerId)

	err := s.service.UpdateSettings(settings)
	if err != nil {
		return &pb.UpdateSettingsResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.UpdateSettingsResponse{
		Success:  true,
		Message:  "Settings updated successfully",
		Settings: s.converter.ToProtoSettings(settings),
	}, nil
}

// SubmitMediaShare submits a new media share request
func (s *MediaShareGRPCHandler) SubmitMediaShare(ctx context.Context, req *pb.SubmitMediaShareRequest) (*pb.SubmitMediaShareResponse, error) {
	mediaReq := &models.MediaShareRequest{
		Type:           s.converter.FromProtoMediaType(req.MediaType),
		URL:            req.MediaUrl,
		Title:          req.CustomTitle,
		Message:        req.CustomDescription,
		DonationAmount: 0, // Will be set from donation record
	}

	response, err := s.service.SubmitMediaShare(
		uint(req.DonationId),
		uint(req.StreamerId),
		uint(req.DonatorId),
		mediaReq,
		req.DonatorName,
	)
	if err != nil {
		return &pb.SubmitMediaShareResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.SubmitMediaShareResponse{
		Success:       true,
		Message:       "Media share submitted successfully",
		MediaId:       uint32(response.ID),
		Status:        s.converter.ToProtoStatus(response.Status),
		QueuePosition: 1, // TODO: Calculate actual queue position
	}, nil
}

// GetMediaQueue retrieves the media queue for a streamer
func (s *MediaShareGRPCHandler) GetMediaQueue(ctx context.Context, req *pb.GetMediaQueueRequest) (*pb.GetMediaQueueResponse, error) {
	items, total, err := s.service.GetMediaQueue(
		uint(req.StreamerId),
		req.StatusFilter,
		int(req.Page),
		int(req.PageSize),
	)
	if err != nil {
		return nil, err
	}

	protoItems := make([]*pb.MediaShareItem, len(items))
	for i, item := range items {
		protoItems[i] = s.converter.ToProtoMediaItem(item)
	}

	totalPages := uint32((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	return &pb.GetMediaQueueResponse{
		Items:      protoItems,
		Total:      uint64(total),
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetMediaStats retrieves media share statistics
func (s *MediaShareGRPCHandler) GetMediaStats(ctx context.Context, req *pb.GetMediaStatsRequest) (*pb.GetMediaStatsResponse, error) {
	stats, err := s.service.GetMediaStats(uint(req.StreamerId))
	if err != nil {
		return nil, err
	}

	return &pb.GetMediaStatsResponse{
		TotalSubmissions:    uint32(stats["total_submissions"]),
		PendingCount:        uint32(stats["pending_count"]),
		ApprovedCount:       uint32(stats["approved_count"]),
		RejectedCount:       uint32(stats["rejected_count"]),
		PlayedCount:         uint32(stats["played_count"]),
		TotalDonationAmount: float64(stats["total_donation_amount"]),
		YoutubeCount:        uint32(stats["youtube_count"]),
		TiktokCount:         uint32(stats["tiktok_count"]),
	}, nil
}

// ApproveMedia approves a media share
func (s *MediaShareGRPCHandler) ApproveMedia(ctx context.Context, req *pb.ApproveMediaRequest) (*pb.ApproveMediaResponse, error) {
	err := s.service.ApproveMedia(uint(req.StreamerId), uint(req.MediaId))
	if err != nil {
		return &pb.ApproveMediaResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.ApproveMediaResponse{
		Success: true,
		Message: "Media approved successfully",
	}, nil
}

// RejectMedia rejects a media share
func (s *MediaShareGRPCHandler) RejectMedia(ctx context.Context, req *pb.RejectMediaRequest) (*pb.RejectMediaResponse, error) {
	err := s.service.RejectMedia(uint(req.StreamerId), uint(req.MediaId))
	if err != nil {
		return &pb.RejectMediaResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.RejectMediaResponse{
		Success: true,
		Message: "Media rejected successfully",
	}, nil
}

// StreamMediaQueue streams real-time media queue updates
func (s *MediaShareGRPCHandler) StreamMediaQueue(req *pb.StreamMediaQueueRequest, stream pb.MediaShareService_StreamMediaQueueServer) error {
	// TODO: Implement real-time streaming using channels and goroutines
	// This would typically push updates when queue changes occur
	return nil
} 