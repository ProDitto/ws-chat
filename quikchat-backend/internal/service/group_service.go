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
	ErrGroupFull        = errors.New("group is full")
)

type groupService struct {
	groupRepo           repository.GroupRepository
	userRepo            repository.UserRepository
	convoRepo           repository.ConversationRepository
	notificationService usecase.NotificationUseCase
	maxGroupMembers   int
}

func NewGroupService(
	groupRepo repository.GroupRepository,
	userRepo repository.UserRepository,
	convoRepo repository.ConversationRepository,
	notificationService usecase.NotificationUseCase,
	maxGroupMembers int,
) usecase.GroupUseCase {
	return &groupService{
		groupRepo:           groupRepo,
		userRepo:            userRepo,
		convoRepo:           convoRepo,
		notificationService: notificationService,
		maxGroupMembers:   maxGroupMembers,
	}
}

func (s *groupService) CreateGroup(ctx context.Context, creatorID int64, name, tag string, description *string) (*domain.Group, error) {
	// 1. Create a conversation for the group
	conversation := &domain.Conversation{
		Type: domain.ConversationTypeGroup,
	}
	if err := s.convoRepo.Save(ctx, conversation); err != nil {
		return nil, fmt.Errorf("failed to create conversation for group: %w", err)
	}

	// 2. Create the group
	group := &domain.Group{
		Name:           name,
		Tag:            tag,
		Description:    description,
		CreatedBy:      creatorID,
		ConversationID: conversation.ID,
	}
	if err := s.groupRepo.Create(ctx, group); err != nil {
		// TODO: Add logic to clean up the created conversation if group creation fails
		return nil, fmt.Errorf("failed to create group: %w", err)
	}

	// 3. Add creator as a participant in the conversation
	if err := s.convoRepo.AddParticipant(ctx, conversation.ID, creatorID); err != nil {
		// TODO: Cleanup logic
		return nil, fmt.Errorf("failed to add creator to conversation: %w", err)
	}

	// 4. Add creator as the owner of the group
	ownerMember := &domain.GroupMember{
		GroupID: group.ID,
		UserID:  creatorID,
		Role:    domain.GroupRoleOwner,
	}
	if err := s.groupRepo.AddMember(ctx, ownerMember); err != nil {
		// TODO: Cleanup logic
		return nil, fmt.Errorf("failed to add owner to group: %w", err)
	}

	return group, nil
}

func (s *groupService) GetGroupDetails(ctx context.Context, userID, groupID int64) (*domain.Group, []*domain.GroupMember, error) {
	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find group: %w", err)
	}
	if group == nil {
		return nil, nil, ErrGroupNotFound
	}

	members, err := s.groupRepo.FindMembersByGroupID(ctx, groupID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to find group members: %w", err)
	}

	return group, members, nil
}

func (s *groupService) SearchGroups(ctx context.Context, tag string, limit int) ([]*domain.Group, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return s.groupRepo.SearchByTag(ctx, tag, limit)
}

func (s *groupService) JoinGroup(ctx context.Context, userID, groupID int64) error {
	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return fmt.Errorf("failed to find group: %w", err)
	}
	if group == nil {
		return ErrGroupNotFound
	}

	role, err := s.groupRepo.FindUserRoleInGroup(ctx, userID, groupID)
	if err != nil {
		return fmt.Errorf("failed to check user role: %w", err)
	}
	if role != nil {
		return ErrAlreadyInGroup
	}

	members, err := s.groupRepo.FindMembersByGroupID(ctx, groupID)
	if err != nil {
		return fmt.Errorf("failed to count group members: %w", err)
	}
	if len(members) >= s.maxGroupMembers {
		return ErrGroupFull
	}

	// Add to group_members
	member := &domain.GroupMember{
		GroupID: groupID,
		UserID:  userID,
		Role:    domain.GroupRoleMember,
	}
	if err := s.groupRepo.AddMember(ctx, member); err != nil {
		return fmt.Errorf("failed to add member to group: %w", err)
	}

	// Add to conversation_participants
	if err := s.convoRepo.AddParticipant(ctx, group.ConversationID, userID); err != nil {
		// TODO: Rollback AddMember
		return fmt.Errorf("failed to add member to conversation: %w", err)
	}

	return nil
}

func (s *groupService) LeaveGroup(ctx context.Context, userID, groupID int64) error {
	role, err := s.groupRepo.FindUserRoleInGroup(ctx, userID, groupID)
	if err != nil {
		return fmt.Errorf("failed to check user role: %w", err)
	}
	if role == nil {
		return ErrNotInGroup
	}
	if *role == domain.GroupRoleOwner {
		return ErrOwnerCannotLeave
	}

	if err := s.groupRepo.RemoveMember(ctx, groupID, userID); err != nil {
		return fmt.Errorf("failed to remove member: %w", err)
	}

	// TODO: Also remove from conversation participants
	return nil
}

func (s *groupService) InviteUser(ctx context.Context, inviterID, groupID int64, inviteeUsername string) error {
	// Check if inviter is in the group
	inviterRole, err := s.groupRepo.FindUserRoleInGroup(ctx, inviterID, groupID)
	if err != nil {
		return fmt.Errorf("failed to check inviter role: %w", err)
	}
	if inviterRole == nil {
		return ErrPermissionDenied // Or ErrNotInGroup
	}

	// Find invitee
	invitee, err := s.userRepo.FindByUsername(ctx, inviteeUsername)
	if err != nil {
		return fmt.Errorf("failed to find invitee: %w", err)
	}
	if invitee == nil {
		return ErrUserNotFound
	}

	// Check if invitee is already in the group
	inviteeRole, err := s.groupRepo.FindUserRoleInGroup(ctx, invitee.ID, groupID)
	if err != nil {
		return fmt.Errorf("failed to check invitee role: %w", err)
	}
	if inviteeRole != nil {
		return ErrAlreadyInGroup
	}

	// Create notification
	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil || group == nil {
		return ErrGroupNotFound
	}

	message := fmt.Sprintf("You have been invited to join the group '%s'.", group.Name)
	_, err = s.notificationService.CreateNotification(ctx, invitee.ID, domain.NotificationGroupInvite, message, &inviterID, &groupID)
	if err != nil {
		return fmt.Errorf("failed to create invitation notification: %w", err)
	}

	return nil
}

func (s *groupService) RemoveUser(ctx context.Context, actorID, groupID, targetUserID int64) error {
	if actorID == targetUserID {
		return ErrCannotRemoveSelf
	}

	actorRole, err := s.groupRepo.FindUserRoleInGroup(ctx, actorID, groupID)
	if err != nil {
		return fmt.Errorf("failed to get actor role: %w", err)
	}
	if actorRole == nil || (*actorRole != domain.GroupRoleOwner && *actorRole != domain.GroupRoleAdmin) {
		return ErrPermissionDenied
	}

	targetRole, err := s.groupRepo.FindUserRoleInGroup(ctx, targetUserID, groupID)
	if err != nil {
		return fmt.Errorf("failed to get target role: %w", err)
	}
	if targetRole == nil {
		return ErrUserNotFound // User is not in the group
	}

	// Owner can remove anyone. Admin can remove members.
	if *actorRole == domain.GroupRoleAdmin && *targetRole != domain.GroupRoleMember {
		return ErrPermissionDenied
	}

	if err := s.groupRepo.RemoveMember(ctx, groupID, targetUserID); err != nil {
		return fmt.Errorf("failed to remove member: %w", err)
	}

	return nil
}

func (s *groupService) UpdateUserRole(ctx context.Context, actorID, groupID, targetUserID int64, newRole domain.GroupRole) error {
	if actorID == targetUserID {
		return ErrCannotUpdateSelf
	}

	actorRole, err := s.groupRepo.FindUserRoleInGroup(ctx, actorID, groupID)
	if err != nil {
		return fmt.Errorf("failed to get actor role: %w", err)
	}
	if actorRole == nil || (*actorRole != domain.GroupRoleOwner && *actorRole != domain.GroupRoleAdmin) {
		return ErrPermissionDenied
	}

	targetRole, err := s.groupRepo.FindUserRoleInGroup(ctx, targetUserID, groupID)
	if err != nil {
		return fmt.Errorf("failed to get target role: %w", err)
	}
	if targetRole == nil {
		return ErrUserNotFound // User is not in the group
	}

	// Only owner can promote to admin or demote admin
	if *actorRole != domain.GroupRoleOwner && (newRole == domain.GroupRoleAdmin || *targetRole == domain.GroupRoleAdmin) {
		return ErrPermissionDenied
	}

	// Cannot change owner role
	if *targetRole == domain.GroupRoleOwner {
		return ErrInvalidRoleUpdate
	}

	// Cannot promote to owner
	if newRole == domain.GroupRoleOwner {
		return ErrInvalidRoleUpdate
	}

	if err := s.groupRepo.UpdateMemberRole(ctx, groupID, targetUserID, newRole); err != nil {
		return fmt.Errorf("failed to update member role: %w", err)
	}

	return nil
}

