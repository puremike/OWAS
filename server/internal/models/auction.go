package models

import "time"

// Auction represents an auction listing
type Auction struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	StartingPrice float64   `json:"starting_price"`
	CurrentPrice  float64   `json:"current_price"`
	Type          string    `json:"type"`   // "english", "dutch", "sealed"
	Status        string    `json:"status"` // "open", "closed"
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	SellerID      string    `json:"seller_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateAuctionRequest struct {
	Title         string  `json:"title" binding:"required"`
	Description   string  `json:"description" binding:"required"`
	StartingPrice float64 `json:"starting_price" binding:"required,gte=1"`
	// CurrentPrice  float64   `json:"current_price" binding:"required,gte=1"`
	Type string `json:"type" binding:"required,oneof=english dutch sealed"`
	// Status        string    `json:"status" binding:"required,oneof=open closed"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
	// SellerID  string `json:"seller_id" binding:"required"`
}

type CreateAuctionResponse struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	StartingPrice float64   `json:"starting_price"`
	CurrentPrice  float64   `json:"current_price"`
	Type          string    `json:"type"`
	Status        string    `json:"status"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	SellerID      string    `json:"seller_id"`
	CreatedAt     time.Time `json:"created_at"`
}
