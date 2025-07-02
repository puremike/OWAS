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
	Username        string `json:"username" binding:"required,min=2,max=64"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,passwd"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
	FullName        string `json:"full_name" binding:"required,min=2,max=64"`
	Location        string `json:"location" binding:"required,min=2,max=64"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
}

type UserProfileUpdateRequest struct {
	Username string `json:"username" binding:"required,min=2,max=64"`
	Email    string `json:"email" binding:"required,email"`
	FullName string `json:"full_name" binding:"required,min=2,max=64"`
	Location string `json:"location" binding:"required,min=2,max=64"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	ID           string `json:"id"`
	Username     string `json:"username"`
}

type PasswordUpdateRequest struct {
	OldPassword     string `json:"old_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,passwd"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}
type PasswordResponse struct {
	Message  string `json:"message"`
	Password string `json:"password"`
}
