package domain

import "time"

type FriendRequestStatus string

const (
	FriendRequestStatusPending  FriendRequestStatus = "pending"
	FriendRequestStatusAccepted FriendRequestStatus = "accepted"
	FriendRequestStatusDeclined FriendRequestStatus = "declined"
)

// FriendRequest represents a request from one user to another to become friends.
type FriendRequest struct {
	ID         int64               `json:"id"`
	SenderID   int64               `json:"sender_id"`
	ReceiverID int64               `json:"receiver_id"`
	Status     FriendRequestStatus `json:"status"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
}

// Friendship represents an established friendship between two users.
type Friendship struct {
	UserID1   int64     `json:"user_id1"`
	UserID2   int64     `json:"user_id2"`
	CreatedAt time.Time `json:"created_at"`
}

