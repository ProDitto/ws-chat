package domain

import "time"

type NotificationType string

const (
	NotificationFriendRequestReceived NotificationType = "friend_request_received"
	NotificationFriendRequestAccepted NotificationType = "friend_request_accepted"
	NotificationGroupInvite           NotificationType = "group_invite"
	NotificationNewMessage            NotificationType = "new_message"
)

type Notification struct {
	ID        int64            `json:"id"`
	UserID    int64            `json:"user_id"`
	Type      NotificationType `json:"type"`
	Message   string           `json:"message"`
	Read      bool             `json:"read"`
	ActorID   *int64           `json:"actor_id,omitempty"`
	ObjectID  *int64           `json:"object_id,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	Actor     *User            `json:"actor,omitempty"` // For rich notifications
}

