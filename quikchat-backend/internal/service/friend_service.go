package service

import (
	"context"
	"errors"
	"quikchat/internal/domain"
	"quikchat/internal/repository"
	"quikchat/internal/usecase"
)

type friendService struct {
	friendRepo repository.FriendRepository
	userRepo   repository.UserRepository
	blockRepo  repository.BlockRepository
}

func NewFriendService(friendRepo repository.FriendRepository, userRepo repository.UserRepository, blockRepo repository.BlockRepository) usecase.FriendUseCase {
	return &friendService{
		friendRepo: friendRepo,
		userRepo:   userRepo,
		blockRepo:  blockRepo,
	}
}

func (s *friendService) SendRequest(ctx context.Context, senderID int64, receiverUsername string) (*domain.FriendRequest, error) {
	receiver, err := s.userRepo.FindByUsername(ctx, receiverUsername)
	if err != nil {
		return nil, err
	}
	if receiver == nil {
		return nil, errors.New("receiver not found")
	}
	if senderID == receiver.ID {
		return nil, errors.New("cannot send friend request to yourself")
	}

	isBlocked, err := s.blockRepo.IsBlocked(ctx, senderID, receiver.ID)
	if err != nil {
		return nil, err
	}
	if isBlocked {
		return nil, errors.New("cannot send friend request to a blocked user")
	}

	existing, err := s.friendRepo.FindRequestBySenderAndReceiver(ctx, senderID, receiver.ID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("friend request already exists")
	}

	req := &domain.FriendRequest{
		SenderID:   senderID,
		ReceiverID: receiver.ID,
		Status:     domain.FriendRequestStatusPending,
	}

	if err := s.friendRepo.SaveRequest(ctx, req); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *friendService) RespondToRequest(ctx context.Context, userID, requestID int64, action string) error {
	req, err := s.friendRepo.FindRequestByID(ctx, requestID)
	if err != nil {
		return err
	}
	if req == nil {
		return errors.New("request not found")
	}
	if req.ReceiverID != userID {
		return errors.New("not authorized to respond to this request")
	}
	if req.Status != domain.FriendRequestStatusPending {
		return errors.New("request already handled")
	}

	var newStatus domain.FriendRequestStatus
	switch action {
	case "accept":
		newStatus = domain.FriendRequestStatusAccepted
	case "decline":
		newStatus = domain.FriendRequestStatusDeclined
	default:
		return errors.New("invalid action")
	}

	return s.friendRepo.UpdateRequestStatus(ctx, requestID, newStatus)
}

func (s *friendService) ListIncomingRequests(ctx context.Context, userID int64) ([]*domain.FriendRequest, error) {
	return s.friendRepo.FindPendingRequestsByReceiverID(ctx, userID)
}

func (s *friendService) ListFriends(ctx context.Context, userID int64) ([]*domain.User, error) {
	return s.friendRepo.FindFriendsByUserID(ctx, userID)
}

func (s *friendService) Unfriend(ctx context.Context, userID int64, friendUsername string) error {
	friend, err := s.userRepo.FindByUsername(ctx, friendUsername)
	if err != nil {
		return err
	}
	if friend == nil {
		return errors.New("friend not found")
	}
	return s.friendRepo.DeleteFriendship(ctx, userID, friend.ID)
}

