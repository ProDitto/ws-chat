package repository

import (
	"context"
	"quikchat/internal/domain"
)

type FriendRepository interface {
	SaveRequest(ctx context.Context, req *domain.FriendRequest) error
	FindRequestByID(ctx context.Context, requestID int64) (*domain.FriendRequest, error)
	FindRequestBySenderAndReceiver(ctx context.Context, senderID, receiverID int64) (*domain.FriendRequest, error)
	UpdateRequestStatus(ctx context.Context, requestID int64, status domain.FriendRequestStatus) error
	FindPendingRequestsByReceiverID(ctx context.Context, userID int64) ([]*domain.FriendRequest, error)
	FindFriendsByUserID(ctx context.Context, userID int64) ([]*domain.User, error)
	DeleteFriendship(ctx context.Context, userID1, userID2 int64) error
}

