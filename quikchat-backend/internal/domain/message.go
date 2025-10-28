package domain

import "time"

type MessageType string

const (
	MessageTypeText  MessageType = "text"
	MessageTypeMedia MessageType = "media"
)

type Message struct {
	ID             int64       `json:"id"`
	ConversationID int64       `json:"conversation_id"`
	SenderID       int64       `json:"sender_id"`
	Type           MessageType `json:"type"`
	Content        string      `json:"content"` // Text content or media file key/URL
	CreatedAt      time.Time   `json:"created_at"`
	Sender         *User       `json:"sender,omitempty"` // Include sender details
}

