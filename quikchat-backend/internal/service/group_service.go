package service

import (
	"context"
	"errors"
	"quikchat/internal/domain"
	"quikchat/internal/repository"
	"quikchat/internal/usecase"
)

var (
	ErrGroupNotFound      = errors.New("group not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrPermissionDenied   = errors.New("permission denied")
	ErrAlreadyInGroup     = errors.New("user is already in the group")
	ErrNotInGroup         = errors.New("user is not in the group")
	ErrOwnerCannotLeave   = errors.New("owner cannot leave the group")
	ErrInvalidRoleUpdate  = errors.New("invalid role update")
	ErrCannotRemoveSelf   = errors.New("cannot remove yourself from a group")
	ErrCannotUpdateSelf   = errors.New("cannot update your own role")
)

type groupService struct {
	groupRepo repository.GroupRepository
	userRepo  repository.UserRepository
	convoRepo repository.ConversationRepository
}

func NewGroupService(groupRepo repository.GroupRepository, userRepo repository.UserRepository, convoRepo repository.ConversationRepository) usecase.GroupUseCase {
	return &groupService{
		groupRepo: groupRepo,
		userRepo:  userRepo,
		convoRepo: convoRepo,
	}
}

func (s *groupService) CreateGroup(ctx context.Context, creatorID int64, name, tag string, description *string) (*domain.Group, error) {
	// In a real app, you'd use a transaction from a UoW pattern
	// For simplicity, we'll do it step-by-step
	convo := &domain.Conversation{Type: domain.ConversationTypeGroup}
	if err := s.convoRepo.Save(ctx, convo); err != nil {
		return nil, err
	}

	group := &domain.Group{
		Name:           name,
		Tag:            tag,
		Description:    description,
		CreatedBy:      creatorID,
		ConversationID: convo.ID,
	}
	if err := s.groupRepo.Create(ctx, group); err != nil {
		// TODO: Rollback conversation creation
		return nil, err
	}

	if err := s.convoRepo.AddParticipant(ctx, convo.ID, creatorID); err != nil {
		// TODO: Rollback
		return nil, err
	}

	ownerMember := &domain.GroupMember{
		GroupID: group.ID,
		UserID:  creatorID,
		Role:    domain.GroupRoleOwner,
	}
	if err := s.groupRepo.AddMember(ctx, ownerMember); err != nil {
		// TODO: Rollback
		return nil, err
	}

	return group, nil
}

func (s *groupService) GetGroupDetails(ctx context.Context, userID, groupID int64) (*domain.Group, []*domain.GroupMember, error) {
	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return nil, nil, err
	}
	if group == nil {
		return nil, nil, ErrGroupNotFound
	}

	members, err := s.groupRepo.FindMembersByGroupID(ctx, groupID)
	if err != nil {
		return nil, nil, err
	}

	return group, members, nil
}

func (s *groupService) SearchGroups(ctx context.Context, tag string, limit int) ([]*domain.Group, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}
	return s.groupRepo.SearchByTag(ctx, tag, limit)
}

func (s *groupService) JoinGroup(ctx context.Context, userID, groupID int64) error {
	role, err := s.groupRepo.FindUserRoleInGroup(ctx, userID, groupID)
	if err != nil {
		return err
	}
	if role != nil {
		return ErrAlreadyInGroup
	}

	member := &domain.GroupMember{
		GroupID: groupID,
		UserID:  userID,
		Role:    domain.GroupRoleMember,
	}
	if err := s.groupRepo.AddMember(ctx, member); err != nil {
		return err
	}

	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil || group == nil {
		// Rollback would be needed here in a real app
		return ErrGroupNotFound
	}

	return s.convoRepo.AddParticipant(ctx, group.ConversationID, userID)
}

func (s *groupService) LeaveGroup(ctx context.Context, userID, groupID int64) error {
	role, err := s.groupRepo.FindUserRoleInGroup(ctx, userID, groupID)
	if err != nil {
		return err
	}
	if role == nil {
		return ErrNotInGroup
	}
	if *role == domain.GroupRoleOwner {
		return ErrOwnerCannotLeave
	}

	return s.groupRepo.RemoveMember(ctx, groupID, userID)
	// Note: We are not removing them from conversation_participants,
	// so they can still see old messages. This is a design choice.
}

func (s *groupService) InviteUser(ctx context.Context, inviterID, groupID int64, inviteeUsername string) error {
	// For now, any member can invite.
	inviterRole, err := s.groupRepo.FindUserRoleInGroup(ctx, inviterID, groupID)
	if err != nil {
		return err
	}
	if inviterRole == nil {
		return ErrPermissionDenied
	}

	invitee, err := s.userRepo.FindByUsername(ctx, inviteeUsername)
	if err != nil {
		return err
	}
	if invitee == nil {
		return ErrUserNotFound
	}

	return s.JoinGroup(ctx, invitee.ID, groupID)
}

func (s *groupService) RemoveUser(ctx context.Context, actorID, groupID, targetUserID int64) error {
	if actorID == targetUserID {
		return ErrCannotRemoveSelf
	}

	actorRole, err := s.groupRepo.FindUserRoleInGroup(ctx, actorID, groupID)
	if err != nil || actorRole == nil {
		return ErrPermissionDenied
	}

	targetRole, err := s.groupRepo.FindUserRoleInGroup(ctx, targetUserID, groupID)
	if err != nil {
		return err
	}
	if targetRole == nil {
		return ErrUserNotFound // User not in group
	}

	// Permission check
	if *actorRole == domain.GroupRoleMember {
		return ErrPermissionDenied
	}
	if *actorRole == domain.GroupRoleAdmin && (*targetRole == domain.GroupRoleAdmin || *targetRole == domain.GroupRoleOwner) {
		return ErrPermissionDenied
	}

	return s.groupRepo.RemoveMember(ctx, groupID, targetUserID)
}

func (s *groupService) UpdateUserRole(ctx context.Context, actorID, groupID, targetUserID int64, newRole domain.GroupRole) error {
	if actorID == targetUserID {
		return ErrCannotUpdateSelf
	}
	if newRole == domain.GroupRoleOwner {
		return ErrInvalidRoleUpdate // Ownership transfer should be a separate, explicit action
	}

	actorRole, err := s.groupRepo.FindUserRoleInGroup(ctx, actorID, groupID)
	if err != nil || actorRole == nil {
		return ErrPermissionDenied
	}

	targetRole, err := s.groupRepo.FindUserRoleInGroup(ctx, targetUserID, groupID)
	if err != nil || targetRole == nil {
		return ErrUserNotFound // User not in group
	}

	// Permission check
	if *actorRole != domain.GroupRoleOwner {
		return ErrPermissionDenied // Only owner can promote/demote
	}
	if *targetRole == domain.GroupRoleOwner {
		return ErrPermissionDenied // Cannot change owner's role
	}

	return s.groupRepo.UpdateMemberRole(ctx, groupID, targetUserID, newRole)
}

