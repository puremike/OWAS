package models

import "time"

// Bid represents a bid placed by a user on an auction
type Bid struct {
	ID        string    `json:"id"`
	AuctionID string    `json:"auction_id"` // FK to Auction.ID
	Auction   Auction   `json:"auction"`
	BidderID  string    `json:"bidder_id"` // FK to UserProfile.ID
	Bidder    User      `json:"bidder"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
