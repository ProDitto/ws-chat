package repository

import (
	"context"
	"quikchat/internal/domain"
)

type MessageRepository interface {
	Save(ctx context.Context, message *domain.Message) error
	FindMessagesByConversationID(ctx context.Context, conversationID int64, cursor int64, limit int) ([]*domain.Message, error)
}

