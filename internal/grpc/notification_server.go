package grpc

import (
	"context"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/rzfd/mediashar/pkg/pb"
)

// NotificationGRPCServer implements the gRPC NotificationService
type NotificationGRPCServer struct {
	pb.UnimplementedNotificationServiceServer
	notificationService NotificationService
	subscribers         map[uint32]chan *pb.DonationEvent
	mu                  sync.RWMutex
}

// NewNotificationGRPCServer creates a new notification gRPC server
func NewNotificationGRPCServer(notificationService NotificationService) *NotificationGRPCServer {
	return &NotificationGRPCServer{
		notificationService: notificationService,
		subscribers:         make(map[uint32]chan *pb.DonationEvent),
	}
}

// SendDonationNotification sends a notification about a donation
func (s *NotificationGRPCServer) SendDonationNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	// Convert request data to map[string]string
	data := make(map[string]string)
	for k, v := range req.Data {
		data[k] = v
	}

	// Send notification via the service
	err := s.notificationService.SendDonationNotification(
		ctx,
		uint(req.UserId),
		req.Title,
		req.Message,
		data,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to send notification: %v", err)
	}

	// Generate notification ID (in production, this would be a UUID)
	notificationID := generateNotificationID(req.UserId)

	return &pb.SendNotificationResponse{
		Success:        true,
		NotificationId: notificationID,
	}, nil
}

// SubscribeDonationEvents subscribes a user to donation events stream
func (s *NotificationGRPCServer) SubscribeDonationEvents(req *pb.SubscribeEventsRequest, stream pb.NotificationService_SubscribeDonationEventsServer) error {
	userID := req.UserId
	
	// Create a channel for this subscription
	eventChan := make(chan *pb.DonationEvent, 100)
	
	// Register subscriber
	s.mu.Lock()
	s.subscribers[userID] = eventChan
	s.mu.Unlock()
	
	// Clean up when stream ends
	defer func() {
		s.mu.Lock()
		delete(s.subscribers, userID)
		close(eventChan)
		s.mu.Unlock()
	}()

	// Send welcome event
	welcomeEvent := &pb.DonationEvent{
		Type:      pb.EventType_EVENT_TYPE_UNSPECIFIED,
		Timestamp: timestamppb.Now(),
		Metadata: map[string]string{
			"message": "Connected to donation events stream",
			"user_id": string(rune(userID)),
		},
	}
	
	if err := stream.Send(welcomeEvent); err != nil {
		return status.Errorf(codes.Internal, "failed to send welcome event: %v", err)
	}

	// Listen for events
	for {
		select {
		case event := <-eventChan:
			if event == nil {
				return nil // Channel closed
			}
			
			// Filter events based on requested types if specified
			if len(req.EventTypes) > 0 {
				allowed := false
				for _, eventType := range req.EventTypes {
					if event.Type == eventType {
						allowed = true
						break
					}
				}
				if !allowed {
					continue
				}
			}
			
			if err := stream.Send(event); err != nil {
				return status.Errorf(codes.Internal, "failed to send event: %v", err)
			}
			
		case <-stream.Context().Done():
			return stream.Context().Err()
		}
	}
}

// BroadcastDonationEvent broadcasts an event to all subscribers
func (s *NotificationGRPCServer) BroadcastDonationEvent(event *pb.DonationEvent) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	for userID, eventChan := range s.subscribers {
		select {
		case eventChan <- event:
			// Event sent successfully
		default:
			// Channel full, skip this subscriber
			// In production, you might want to log this or implement backpressure
			_ = userID // Avoid unused variable warning
		}
	}
}

// BroadcastDonationCreated broadcasts a donation created event
func (s *NotificationGRPCServer) BroadcastDonationCreated(donationID uint32, streamerID uint32, amount float64) {
	event := &pb.DonationEvent{
		Type:      pb.EventType_EVENT_TYPE_DONATION_CREATED,
		Timestamp: timestamppb.Now(),
		Metadata: map[string]string{
			"donation_id": string(rune(donationID)),
			"streamer_id": string(rune(streamerID)),
			"amount":      string(rune(int64(amount))),
		},
	}
	
	s.BroadcastDonationEvent(event)
}

// BroadcastPaymentCompleted broadcasts a payment completed event
func (s *NotificationGRPCServer) BroadcastPaymentCompleted(donationID uint32, transactionID string) {
	event := &pb.DonationEvent{
		Type:      pb.EventType_EVENT_TYPE_PAYMENT_VERIFIED,
		Timestamp: timestamppb.Now(),
		Metadata: map[string]string{
			"donation_id":    string(rune(donationID)),
			"transaction_id": transactionID,
		},
	}
	
	s.BroadcastDonationEvent(event)
}

// GetSubscriberCount returns the number of active subscribers
func (s *NotificationGRPCServer) GetSubscriberCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.subscribers)
}

// Helper functions

func generateNotificationID(userID uint32) string {
	// Generate a simple notification ID
	// In production, use UUID or other unique identifier
	return "notif_" + string(rune(userID)) + "_" + string(rune(time.Now().Unix()))
}

// MockNotificationService implements NotificationService for testing
type MockNotificationService struct{}

func NewMockNotificationService() *MockNotificationService {
	return &MockNotificationService{}
}

func (m *MockNotificationService) SendDonationNotification(ctx context.Context, userID uint, title, message string, data map[string]string) error {
	// Mock implementation - in production this would send push notifications,
	// emails, SMS, or other notification methods
	return nil
}

func (m *MockNotificationService) SubscribeEvents(ctx context.Context, userID uint, eventTypes []string) (<-chan *pb.DonationEvent, error) {
	// Mock implementation - in production this would connect to a message queue
	// or event streaming system
	eventChan := make(chan *pb.DonationEvent, 10)
	
	// Send a mock event
	go func() {
		defer close(eventChan)
		
		mockEvent := &pb.DonationEvent{
			Type:      pb.EventType_EVENT_TYPE_DONATION_CREATED,
			Timestamp: timestamppb.Now(),
			Metadata: map[string]string{
				"user_id": string(rune(userID)),
				"mock":    "true",
			},
		}
		
		select {
		case eventChan <- mockEvent:
		case <-ctx.Done():
		}
	}()
	
	return eventChan, nil
} 