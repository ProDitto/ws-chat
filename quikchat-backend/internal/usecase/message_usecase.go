package usecase

import (
	"context"
	"quikchat/internal/domain"
)

type MessageUseCase interface {
	SendMessage(ctx context.Context, senderID, conversationID int64, msgType domain.MessageType, content string) (*domain.Message, error)
	GetMessagesByConversationID(ctx context.Context, userID, conversationID int64, cursor int64, limit int) ([]*domain.Message, error)
	GetPresignedUploadURL(ctx context.Context, userID int64, filename string) (string, string, error)
}

