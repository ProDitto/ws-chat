package usecase

import (
	"context"
	"quikchat/internal/domain"
)

type NotificationUseCase interface {
	CreateNotification(ctx context.Context, userID int64, nType domain.NotificationType, message string, actorID, objectID *int64) (*domain.Notification, error)
	GetNotifications(ctx context.Context, userID int64, limit, offset int) ([]*domain.Notification, error)
	MarkNotificationAsRead(ctx context.Context, userID, notificationID int64) error
	MarkAllNotificationsAsRead(ctx context.Context, userID int64) error
	GetUnreadCount(ctx context.Context, userID int64) (int, error)
}

