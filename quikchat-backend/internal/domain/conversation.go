package domain

import "time"

type ConversationType string

const (
	ConversationTypePrivate ConversationType = "private"
	ConversationTypeGroup   ConversationType = "group"
)

type Conversation struct {
	ID        int64            `json:"id"`
	Type      ConversationType `json:"type"`
	CreatedAt time.Time        `json:"created_at"`
}

type ConversationParticipant struct {
	ConversationID int64     `json:"conversation_id"`
	UserID         int64     `json:"user_id"`
	JoinedAt       time.Time `json:"joined_at"`
}

