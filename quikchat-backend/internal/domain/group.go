package domain

import "time"

type GroupRole string

const (
	GroupRoleOwner  GroupRole = "owner"
	GroupRoleAdmin  GroupRole = "admin"
	GroupRoleMember GroupRole = "member"
)

type Group struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Tag             string    `json:"tag"`
	Description     *string   `json:"description,omitempty"`
	ProfileImageURL *string   `json:"profile_image_url,omitempty"`
	CreatedBy       int64     `json:"created_by"`
	ConversationID  int64     `json:"conversation_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type GroupMember struct {
	GroupID  int64     `json:"group_id"`
	UserID   int64     `json:"user_id"`
	Role     GroupRole `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
	User     *User     `json:"user,omitempty"` // For returning member lists with user details
}

