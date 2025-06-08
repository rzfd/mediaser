package adapter

import (
	"github.com/rzfd/mediashar/internal/models"
	"github.com/rzfd/mediashar/pkg/pb"
)

// MediaShareConverter handles conversions between internal models and protobuf messages
type MediaShareConverter struct{}

// NewMediaShareConverter creates a new converter instance
func NewMediaShareConverter() *MediaShareConverter {
	return &MediaShareConverter{}
}

// ToProtoSettings converts MediaShareSettings to protobuf
func (c *MediaShareConverter) ToProtoSettings(settings *models.MediaShareSettings) *pb.MediaShareSettings {
	return &pb.MediaShareSettings{
		Id:                 uint32(settings.ID),
		StreamerId:         uint32(settings.StreamerID),
		Enabled:            settings.MediaShareEnabled,
		MinDonationAmount:  settings.MinDonationAmount,
		YoutubeEnabled:     settings.AllowYoutube,
		TiktokEnabled:      settings.AllowTiktok,
		AutoApprove:        settings.AutoApprove,
		MaxDurationSeconds: uint32(settings.MaxDurationYoutube), // Use YouTube as default
		WelcomeMessage:     settings.WelcomeMessage,
	}
}

// FromProtoSettings converts protobuf to MediaShareSettings
func (c *MediaShareConverter) FromProtoSettings(proto *pb.MediaShareSettings) *models.MediaShareSettings {
	return &models.MediaShareSettings{
		MediaShareEnabled:  proto.Enabled,
		MinDonationAmount:  proto.MinDonationAmount,
		AllowYoutube:       proto.YoutubeEnabled,
		AllowTiktok:        proto.TiktokEnabled,
		AutoApprove:        proto.AutoApprove,
		MaxDurationYoutube: int(proto.MaxDurationSeconds),
		MaxDurationTiktok:  int(proto.MaxDurationSeconds),
		WelcomeMessage:     proto.WelcomeMessage,
	}
}

// FromProtoMediaType converts protobuf MediaType to model MediaShareType
func (c *MediaShareConverter) FromProtoMediaType(protoType pb.MediaType) models.MediaShareType {
	switch protoType {
	case pb.MediaType_MEDIA_TYPE_YOUTUBE:
		return models.MediaShareTypeYoutube
	case pb.MediaType_MEDIA_TYPE_TIKTOK:
		return models.MediaShareTypeTiktok
	default:
		return models.MediaShareTypeYoutube
	}
}

// ToProtoMediaType converts model MediaShareType to protobuf MediaType
func (c *MediaShareConverter) ToProtoMediaType(mediaType models.MediaShareType) pb.MediaType {
	switch mediaType {
	case models.MediaShareTypeYoutube:
		return pb.MediaType_MEDIA_TYPE_YOUTUBE
	case models.MediaShareTypeTiktok:
		return pb.MediaType_MEDIA_TYPE_TIKTOK
	default:
		return pb.MediaType_MEDIA_TYPE_YOUTUBE
	}
}

// ToProtoStatus converts model MediaShareStatus to protobuf MediaStatus
func (c *MediaShareConverter) ToProtoStatus(status models.MediaShareStatus) pb.MediaStatus {
	switch status {
	case models.MediaShareStatusPending:
		return pb.MediaStatus_MEDIA_STATUS_PENDING
	case models.MediaShareStatusApproved:
		return pb.MediaStatus_MEDIA_STATUS_APPROVED
	case models.MediaShareStatusRejected:
		return pb.MediaStatus_MEDIA_STATUS_REJECTED
	default:
		return pb.MediaStatus_MEDIA_STATUS_PENDING
	}
}

// ToProtoMediaItem converts MediaQueueItem to protobuf MediaShareItem
func (c *MediaShareConverter) ToProtoMediaItem(item *models.MediaQueueItem) *pb.MediaShareItem {
	return &pb.MediaShareItem{
		Id:                uint32(item.ID),
		DonationId:        0, // TODO: Add DonationID to MediaQueueItem
		StreamerId:        0, // TODO: Add StreamerID to MediaQueueItem
		DonatorId:         0, // TODO: Add DonatorID to MediaQueueItem
		DonatorName:       item.DonatorName,
		MediaType:         c.ToProtoMediaType(item.Type),
		MediaUrl:          item.URL,
		CustomTitle:       item.Title,
		CustomDescription: item.Message,
		ThumbnailUrl:      item.Thumbnail,
		DurationSeconds:   0, // TODO: Add Duration to MediaQueueItem
		StartTime:         0,
		EndTime:           0,
		Status:            c.ToProtoStatus(item.Status),
		DonationAmount:    item.DonationAmount,
	}
} 