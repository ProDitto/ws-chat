package usecase

import (
	"context"
	"quikchat/internal/domain"
)

type FriendUseCase interface {
	SendRequest(ctx context.Context, senderID int64, receiverUsername string) (*domain.FriendRequest, error)
	RespondToRequest(ctx context.Context, userID, requestID int64, action string) error
	ListIncomingRequests(ctx context.Context, userID int64) ([]*domain.FriendRequest, error)
	ListFriends(ctx context.Context, userID int64) ([]*domain.User, error)
	Unfriend(ctx context.Context, userID int64, friendUsername string) error
}

