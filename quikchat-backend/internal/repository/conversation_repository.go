package repository

import (
	"context"
	"quikchat/internal/domain"
)

type ConversationRepository interface {
	Save(ctx context.Context, conversation *domain.Conversation) error
	AddParticipant(ctx context.Context, conversationID, userID int64) error
	IsUserInConversation(ctx context.Context, userID, conversationID int64) (bool, error)
	FindParticipantIDs(ctx context.Context, conversationID int64) ([]int64, error)
}

