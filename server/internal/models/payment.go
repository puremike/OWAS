package models

import "time"

type Payment struct {
	ID        string    `json:"id"`
	AuctionID string    `json:"auction_id"`
	BuyerID   string    `json:"buyer_id"`
	OrderID   string    `json:"order_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"` // pending, completed, failed
	SessionID string    `json:"session_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePaymentRequest struct {
	Amount int64 `json:"amount" binding:"required"`
}

type CreatePaymentResponse struct {
	CheckoutURL string `json:"checkout_url"`
}
