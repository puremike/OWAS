package models

import "time"

type ContactSupport struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

type ContactSupportReq struct {
	Subject string `json:"subject" binding:"required,min=1,max=100"`
	Message string `json:"message" binding:"required,min=1,max=1000"`
}

type SupportRes struct {
	ID      int    `json:"id"`
	UserID  string `json:"user_id"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}
