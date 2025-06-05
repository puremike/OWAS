package models

import "time"

// Notification for user alerts and messages
type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"` // FK to User.ID
	User      User      `json:"user"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}
