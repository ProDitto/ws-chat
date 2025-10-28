package repository

import (
	"context"
	"quikchat/internal/domain"
)

type GroupRepository interface {
	Create(ctx context.Context, group *domain.Group) error
	FindByID(ctx context.Context, groupID int64) (*domain.Group, error)
	SearchByTag(ctx context.Context, tag string, limit int) ([]*domain.Group, error)
	AddMember(ctx context.Context, member *domain.GroupMember) error
	UpdateMemberRole(ctx context.Context, groupID, userID int64, role domain.GroupRole) error
	RemoveMember(ctx context.Context, groupID, userID int64) error
	FindMembersByGroupID(ctx context.Context, groupID int64) ([]*domain.GroupMember, error)
	FindUserRoleInGroup(ctx context.Context, userID, groupID int64) (*domain.GroupRole, error)
}

