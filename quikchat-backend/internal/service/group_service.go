package service

import (
	"context"
	"errors"
	"fmt"
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
	ErrOwnerCannotLeave   = errors.New("group owner cannot leave the group")
	ErrInvalidRoleUpdate  = errors.New("invalid role update")
	ErrCannotRemoveSelf   = errors.New("cannot remove yourself from a group")
	ErrCannotUpdateSelf   = errors.New("cannot update your own role")
)

type groupService struct {
	groupRepo           repository.GroupRepository
	userRepo            repository.UserRepository
	convoRepo           repository.ConversationRepository
	notificationService usecase.NotificationUseCase
}

func NewGroupService(groupRepo repository.GroupRepository, userRepo repository.UserRepository, convoRepo repository.ConversationRepository, notificationService usecase.NotificationUseCase) usecase.GroupUseCase {
	return &groupService{
		groupRepo:           groupRepo,
		userRepo:            userRepo,
		convoRepo:           convoRepo,
		notificationService: notificationService,
	}
}

func (s *groupService) CreateGroup(ctx context.Context, creatorID int64, name, tag string, description *string) (*domain.Group, error) {
	// In a real application, this should be a transaction
	conversation := &domain.Conversation{Type: domain.ConversationTypeGroup}
	if err := s.convoRepo.Save(ctx, conversation); err != nil {
		return nil, err
	}

	group := &domain.Group{
		Name:           name,
		Tag:            tag,
		Description:    description,
		CreatedBy:      creatorID,
		ConversationID: conversation.ID,
	}
	if err := s.groupRepo.Create(ctx, group); err != nil {
		return nil, err
	}

	if err := s.convoRepo.AddParticipant(ctx, conversation.ID, creatorID); err != nil {
		return nil, err
	}

	ownerMember := &domain.GroupMember{
		GroupID: group.ID,
		UserID:  creatorID,
		Role:    domain.GroupRoleOwner,
	}
	if err := s.groupRepo.AddMember(ctx, ownerMember); err != nil {
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
		limit = 10
	}
	return s.groupRepo.SearchByTag(ctx, tag, limit)
}

func (s *groupService) JoinGroup(ctx context.Context, userID, groupID int64) error {
	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return err
	}
	if group == nil {
		return ErrGroupNotFound
	}

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

	return s.convoRepo.AddParticipant(ctx, group.ConversationID, userID)
}

func (s *groupService) LeaveGroup(ctx context.Context, userID, groupID int64) error {
	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return err
	}
	if group == nil {
		return ErrGroupNotFound
	}

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

	// TODO: remove from conversation participants as well
	return s.groupRepo.RemoveMember(ctx, groupID, userID)
}

func (s *groupService) InviteUser(ctx context.Context, inviterID, groupID int64, inviteeUsername string) error {
	inviterRole, err := s.groupRepo.FindUserRoleInGroup(ctx, inviterID, groupID)
	if err != nil {
		return err
	}
	if inviterRole == nil {
		return ErrPermissionDenied // Inviter is not in the group
	}

	invitee, err := s.userRepo.FindByUsername(ctx, inviteeUsername)
	if err != nil || invitee == nil {
		return ErrUserNotFound
	}

	inviteeRole, err := s.groupRepo.FindUserRoleInGroup(ctx, invitee.ID, groupID)
	if err != nil {
		return err
	}
	if inviteeRole != nil {
		return ErrAlreadyInGroup
	}

	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil || group == nil {
		return ErrGroupNotFound
	}

	inviter, err := s.userRepo.FindByID(ctx, inviterID)
	if err != nil || inviter == nil {
		return ErrUserNotFound
	}

	// Create notification
	message := fmt.Sprintf("%s invited you to join the group '%s'.", inviter.Username, group.Name)
	actorID := inviterID
	objectID := groupID
	_, _ = s.notificationService.CreateNotification(ctx, invitee.ID, domain.NotificationGroupInvite, message, &actorID, &objectID)

	return nil
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
	if err != nil || targetRole == nil {
		return ErrNotInGroup
	}

	// Permission checks
	if *actorRole == domain.GroupRoleMember {
		return ErrPermissionDenied // Members cannot remove anyone
	}
	if *actorRole == domain.GroupRoleAdmin && (*targetRole == domain.GroupRoleAdmin || *targetRole == domain.GroupRoleOwner) {
		return ErrPermissionDenied // Admins cannot remove other admins or the owner
	}

	return s.groupRepo.RemoveMember(ctx, groupID, targetUserID)
}

func (s *groupService) UpdateUserRole(ctx context.Context, actorID, groupID, targetUserID int64, newRole domain.GroupRole) error {
	if actorID == targetUserID {
		return ErrCannotUpdateSelf
	}

	actorRole, err := s.groupRepo.FindUserRoleInGroup(ctx, actorID, groupID)
	if err != nil || actorRole == nil {
		return ErrPermissionDenied
	}

	if *actorRole != domain.GroupRoleOwner {
		return ErrPermissionDenied // Only owner can change roles
	}

	targetRole, err := s.groupRepo.FindUserRoleInGroup(ctx, targetUserID, groupID)
	if err != nil || targetRole == nil {
		return ErrNotInGroup
	}

	if *targetRole == domain.GroupRoleOwner {
		return ErrInvalidRoleUpdate // Cannot change owner's role
	}
	if newRole == domain.GroupRoleOwner {
		return ErrInvalidRoleUpdate // Cannot promote to owner (requires separate transfer ownership flow)
	}

	return s.groupRepo.UpdateMemberRole(ctx, groupID, targetUserID, newRole)
}
