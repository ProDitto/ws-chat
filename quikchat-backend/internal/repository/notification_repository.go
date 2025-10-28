package repository

import (
	"context"
	"quikchat/internal/domain"
)

type NotificationRepository interface {
	Save(ctx context.Context, notification *domain.Notification) error
	FindByUserID(ctx context.Context, userID int64, limit, offset int) ([]*domain.Notification, error)
	MarkAsRead(ctx context.Context, notificationID, userID int64) error
	MarkAllAsRead(ctx context.Context, userID int64) error
	CountUnread(ctx context.Context, userID int64) (int, error)
}

