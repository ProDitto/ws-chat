package service

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"quikchat/internal/adapter/external/s3"
	"quikchat/internal/adapter/ws"
	"quikchat/internal/domain"
	"quikchat/internal/repository"
	"quikchat/internal/usecase"
	"time"

	"github.com/google/uuid"
)

type messageService struct {
	msgRepo    repository.MessageRepository
	convoRepo  repository.ConversationRepository // To check if user is in conversation
	userRepo   repository.UserRepository
	s3Client   s3.S3Client
	hub        *ws.Hub
	s3Bucket   string
}

func NewMessageService(msgRepo repository.MessageRepository, convoRepo repository.ConversationRepository, userRepo repository.UserRepository, s3Client s3.S3Client, hub *ws.Hub) usecase.MessageUseCase {
	return &messageService{
		msgRepo:   msgRepo,
		convoRepo: convoRepo,
		userRepo:  userRepo,
		s3Client:  s3Client,
		hub:       hub,
		s3Bucket:  "quikchat-media", // This should come from config
	}
}

func (s *messageService) SendMessage(ctx context.Context, senderID, conversationID int64, msgType domain.MessageType, content string) (*domain.Message, error) {
	// 1. Validate that the user is part of the conversation
	isParticipant, err := s.convoRepo.IsUserInConversation(ctx, senderID, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to check conversation participation: %w", err)
	}
	if !isParticipant {
		return nil, errors.New("user is not a participant in this conversation")
	}

	// 2. Create and save the message
	message := &domain.Message{
		ConversationID: conversationID,
		SenderID:       senderID,
		Type:           msgType,
		Content:        content,
		CreatedAt:      time.Now(),
	}

	if err := s.msgRepo.Save(ctx, message); err != nil {
		return nil, fmt.Errorf("failed to save message: %w", err)
	}

	// 3. Populate sender info for broadcast
	sender, err := s.userRepo.FindByID(ctx, senderID)
	if err != nil {
		return nil, fmt.Errorf("failed to find sender: %w", err)
	}
	// Don't send password hash over the wire
	sender.PasswordHash = ""
	message.Sender = sender

	// 4. Broadcast the message to other participants
	participantIDs, err := s.convoRepo.FindParticipantIDs(ctx, conversationID)
	if err != nil {
		// Log the error but don't fail the whole operation, the message is already saved.
		fmt.Printf("failed to get participant IDs for broadcast: %v\n", err)
	} else {
		s.hub.BroadcastToUsers(message, participantIDs)
	}

	return message, nil
}

func (s *messageService) GetMessagesByConversationID(ctx context.Context, userID, conversationID int64, cursor int64, limit int) ([]*domain.Message, error) {
	isParticipant, err := s.convoRepo.IsUserInConversation(ctx, userID, conversationID)
	if err != nil {
		return nil, fmt.Errorf("failed to check conversation participation: %w", err)
	}
	if !isParticipant {
		return nil, errors.New("user is not a participant in this conversation")
	}

	if limit <= 0 || limit > 100 {
		limit = 50 // Default/max limit
	}

	return s.msgRepo.FindMessagesByConversationID(ctx, conversationID, cursor, limit)
}

func (s *messageService) GetPresignedUploadURL(ctx context.Context, userID int64, filename string) (string, string, error) {
	ext := filepath.Ext(filename)
	key := fmt.Sprintf("uploads/%d/%s%s", userID, uuid.New().String(), ext)

	url, err := s.s3Client.GeneratePresignedUploadURL(ctx, s.s3Bucket, key, 15*time.Minute)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url, key, nil
}

