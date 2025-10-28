package domain

import "time"

type User struct {
	ID              int64      `json:"id"`
	Username        string     `json:"username"`
	PasswordHash    string     `json:"-"` // Do not expose password hash
	DisplayName     string     `json:"display_name"`
	ProfileImageURL *string    `json:"profile_image_url"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	LastSeenAt      *time.Time `json:"last_seen_at"`
}

