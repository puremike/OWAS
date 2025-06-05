package models

import "time"

type Payment struct {
	ID        string    `json:"id"`
	AuctionID string    `json:"auction_id"`
	BuyerID   string    `json:"buyer_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"` // pending, completed, failed
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
