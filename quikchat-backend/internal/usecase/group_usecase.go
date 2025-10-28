package usecase

import (
	"context"
	"quikchat/internal/domain"
)

type GroupUseCase interface {
	CreateGroup(ctx context.Context, creatorID int64, name, tag string, description *string) (*domain.Group, error)
	GetGroupDetails(ctx context.Context, userID, groupID int64) (*domain.Group, []*domain.GroupMember, error)
	SearchGroups(ctx context.Context, tag string, limit int) ([]*domain.Group, error)
	JoinGroup(ctx context.Context, userID, groupID int64) error
	LeaveGroup(ctx context.Context, userID, groupID int64) error
	InviteUser(ctx context.Context, inviterID, groupID int64, inviteeUsername string) error
	RemoveUser(ctx context.Context, actorID, groupID, targetUserID int64) error
	UpdateUserRole(ctx context.Context, actorID, groupID, targetUserID int64, newRole domain.GroupRole) error
}

