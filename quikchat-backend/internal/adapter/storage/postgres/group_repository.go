package postgres

import (
	"context"
	"errors"
	"quikchat/internal/domain"
	"quikchat/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GroupRepository struct {
	db *pgxpool.Pool
}

func NewGroupRepository(db *pgxpool.Pool) repository.GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) Create(ctx context.Context, group *domain.Group) error {
	query := `INSERT INTO groups (name, tag, description, profile_image_url, created_by, conversation_id)
              VALUES ($1, $2, $3, $4, $5, $6)
              RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query, group.Name, group.Tag, group.Description, group.ProfileImageURL, group.CreatedBy, group.ConversationID).
		Scan(&group.ID, &group.CreatedAt, &group.UpdatedAt)
	return err
}

func (r *GroupRepository) FindByID(ctx context.Context, groupID int64) (*domain.Group, error) {
	query := `SELECT id, name, tag, description, profile_image_url, created_by, conversation_id, created_at, updated_at
              FROM groups WHERE id = $1`
	group := &domain.Group{}
	err := r.db.QueryRow(ctx, query, groupID).
		Scan(&group.ID, &group.Name, &group.Tag, &group.Description, &group.ProfileImageURL, &group.CreatedBy, &group.ConversationID, &group.CreatedAt, &group.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return group, nil
}

func (r *GroupRepository) SearchByTag(ctx context.Context, tag string, limit int) ([]*domain.Group, error) {
	query := `SELECT id, name, tag, description, profile_image_url, created_by, conversation_id, created_at, updated_at
              FROM groups WHERE tag ILIKE $1 ORDER BY similarity(tag, $2) DESC LIMIT $3`
	rows, err := r.db.Query(ctx, query, "%"+tag+"%", tag, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*domain.Group
	for rows.Next() {
		group := &domain.Group{}
		err := rows.Scan(&group.ID, &group.Name, &group.Tag, &group.Description, &group.ProfileImageURL, &group.CreatedBy, &group.ConversationID, &group.CreatedAt, &group.UpdatedAt)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (r *GroupRepository) AddMember(ctx context.Context, member *domain.GroupMember) error {
	query := `INSERT INTO group_members (group_id, user_id, role) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	_, err := r.db.Exec(ctx, query, member.GroupID, member.UserID, member.Role)
	return err
}

func (r *GroupRepository) UpdateMemberRole(ctx context.Context, groupID, userID int64, role domain.GroupRole) error {
	query := `UPDATE group_members SET role = $3 WHERE group_id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, groupID, userID, role)
	return err
}

func (r *GroupRepository) RemoveMember(ctx context.Context, groupID, userID int64) error {
	query := `DELETE FROM group_members WHERE group_id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, groupID, userID)
	return err
}

func (r *GroupRepository) FindMembersByGroupID(ctx context.Context, groupID int64) ([]*domain.GroupMember, error) {
	query := `SELECT gm.group_id, gm.user_id, gm.role, gm.joined_at, u.username, u.display_name, u.profile_image_url
              FROM group_members gm
              JOIN users u ON gm.user_id = u.id
              WHERE gm.group_id = $1
              ORDER BY gm.role, u.display_name`
	rows, err := r.db.Query(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*domain.GroupMember
	for rows.Next() {
		member := &domain.GroupMember{User: &domain.User{}}
		err := rows.Scan(&member.GroupID, &member.UserID, &member.Role, &member.JoinedAt, &member.User.Username, &member.User.DisplayName, &member.User.ProfileImageURL)
		if err != nil {
			return nil, err
		}
		member.User.ID = member.UserID
		members = append(members, member)
	}
	return members, nil
}

func (r *GroupRepository) FindUserRoleInGroup(ctx context.Context, userID, groupID int64) (*domain.GroupRole, error) {
	query := `SELECT role FROM group_members WHERE user_id = $1 AND group_id = $2`
	var role domain.GroupRole
	err := r.db.QueryRow(ctx, query, userID, groupID).Scan(&role)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

