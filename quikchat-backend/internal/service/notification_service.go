package service

import (
	"context"
	"quikchat/internal/adapter/ws"
	"quikchat/internal/domain"
	"quikchat/internal/repository"
	"quikchat/internal/usecase"
)

type notificationService struct {
	notificationRepo repository.NotificationRepository
	userRepo         repository.UserRepository
	hub              *ws.Hub
}

func NewNotificationService(notificationRepo repository.NotificationRepository, userRepo repository.UserRepository, hub *ws.Hub) usecase.NotificationUseCase {
	return &notificationService{
		notificationRepo: notificationRepo,
		userRepo:         userRepo,
		hub:              hub,
	}
}

func (s *notificationService) CreateNotification(ctx context.Context, userID int64, nType domain.NotificationType, message string, actorID, objectID *int64) (*domain.Notification, error) {
	notification := &domain.Notification{
		UserID:   userID,
		Type:     nType,
		Message:  message,
		ActorID:  actorID,
		ObjectID: objectID,
	}

	if err := s.notificationRepo.Save(ctx, notification); err != nil {
		return nil, err
	}

	// Populate actor details for real-time push
	if actorID != nil {
		actor, err := s.userRepo.FindByID(ctx, *actorID)
		if err == nil && actor != nil {
			actor.PasswordHash = "" // Ensure password hash is not exposed
			notification.Actor = actor
		}
	}

	// Push notification via WebSocket
	payload := map[string]interface{}{
		"type":    "new_notification",
		"payload": notification,
	}
	s.hub.BroadcastToUsers(payload, []int64{userID})

	return notification, nil
}

func (s *notificationService) GetNotifications(ctx context.Context, userID int64, limit, offset int) ([]*domain.Notification, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.notificationRepo.FindByUserID(ctx, userID, limit, offset)
}

func (s *notificationService) MarkNotificationAsRead(ctx context.Context, userID, notificationID int64) error {
	return s.notificationRepo.MarkAsRead(ctx, notificationID, userID)
}

func (s *notificationService) MarkAllNotificationsAsRead(ctx context.Context, userID int64) error {
	return s.notificationRepo.MarkAllAsRead(ctx, userID)
}

func (s *notificationService) GetUnreadCount(ctx context.Context, userID int64) (int, error) {
	return s.notificationRepo.CountUnread(ctx, userID)
}

