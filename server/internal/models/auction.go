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
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	StartingPrice float64 `json:"starting_price"`
	// CurrentPrice  float64   `json:"current_price"`
	Type string `json:"type"`
	// Status        string    `json:"status"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	SellerID  string    `json:"seller_id"`
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
