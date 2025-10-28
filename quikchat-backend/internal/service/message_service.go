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

	"github.com/google/uuid"
)

type messageService struct {
	msgRepo             repository.MessageRepository
	convoRepo           repository.ConversationRepository
	userRepo            repository.UserRepository
	notificationService usecase.NotificationUseCase
	s3Client            s3.S3Client
	hub                 *ws.Hub
	s3Bucket            string
}

func NewMessageService(msgRepo repository.MessageRepository, convoRepo repository.ConversationRepository, userRepo repository.UserRepository, notificationService usecase.NotificationUseCase, s3Client s3.S3Client, hub *ws.Hub) usecase.MessageUseCase {
	return &messageService{
		msgRepo:             msgRepo,
		convoRepo:           convoRepo,
		userRepo:            userRepo,
		notificationService: notificationService,
		s3Client:            s3Client,
		hub:                 hub,
		s3Bucket:            "quikchat-media", // This should come from config
	}
}

func (s *messageService) SendMessage(ctx context.Context, senderID, conversationID int64, msgType domain.MessageType, content string) (*domain.Message, error) {
	isParticipant, err := s.convoRepo.IsUserInConversation(ctx, senderID, conversationID)
	if err != nil {
		return nil, err
	}
	if !isParticipant {
		return nil, errors.New("user is not a participant in this conversation")
	}

	message := &domain.Message{
		ConversationID: conversationID,
		SenderID:       senderID,
		Type:           msgType,
		Content:        content,
	}

	if err := s.msgRepo.Save(ctx, message); err != nil {
		return nil, err
	}

	sender, err := s.userRepo.FindByID(ctx, senderID)
	if err != nil {
		return nil, err
	}
	sender.PasswordHash = ""
	message.Sender = sender

	participantIDs, err := s.convoRepo.FindParticipantIDs(ctx, conversationID)
	if err != nil {
		return nil, err
	}

	// Broadcast message via WebSocket
	payload := map[string]interface{}{
		"type":    "new_message",
		"payload": message,
	}
	s.hub.BroadcastToUsers(payload, participantIDs)

	// Create notifications for other participants
	notificationMessage := fmt.Sprintf("%s sent you a message.", sender.Username)
	actorID := senderID
	objectID := message.ID
	for _, userID := range participantIDs {
		if userID != senderID {
			_, _ = s.notificationService.CreateNotification(context.Background(), userID, domain.NotificationNewMessage, notificationMessage, &actorID, &objectID)
		}
	}

	return message, nil
}

func (s *messageService) GetMessagesByConversationID(ctx context.Context, userID, conversationID int64, cursor int64, limit int) ([]*domain.Message, error) {
	isParticipant, err := s.convoRepo.IsUserInConversation(ctx, userID, conversationID)
	if err != nil {
		return nil, err
	}
	if !isParticipant {
		return nil, errors.New("user is not a participant in this conversation")
	}

	if limit <= 0 || limit > 100 {
		limit = 50
	}

	return s.msgRepo.FindMessagesByConversationID(ctx, conversationID, cursor, limit)
}

func (s *messageService) GetPresignedUploadURL(ctx context.Context, userID int64, filename string) (string, string, error) {
	ext := filepath.Ext(filename)
	key := fmt.Sprintf("uploads/%d/%s%s", userID, uuid.New().String(), ext)

	uploadURL, err := s.s3Client.GeneratePresignedUploadURL(ctx, s.s3Bucket, key, 0)
	if err != nil {
		return "", "", err
	}

	// In a real scenario, you'd return the final object URL or just the key
	// For simplicity, we'll just return the key as the content for the message
	return uploadURL, key, nil
}
