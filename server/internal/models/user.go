package models

import "time"

// UserProfile represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	FullName  string    `json:"full_name"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsAdmin   bool      `json:"is_admin"`
}

type CreateUserRequest struct {
	Username        string `json:"username" binding:"required,min=6,max=32"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,passwd"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	FullName        string `json:"full_name" binding:"required,min=6,max=32"`
	Location        string `json:"location" binding:"required,min=6,max=42"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	Location  string `json:"location"`
	CreatedAt string `json:"created_at"`
}
