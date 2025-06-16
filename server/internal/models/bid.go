package models

import "time"

// Bid represents a bid placed by a user on an auction
type Bid struct {
	ID        string    `json:"id"`
	AuctionID string    `json:"auction_id"` // FK to Auction.ID
	BidderID  string    `json:"bidder_id"`  // FK to UserProfile.ID
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type PlaceBidRequest struct {
	AuctionID string  `json:"auction_id"`
	BidderID  string  `json:"bidder_id"`
	BidAmount float64 `json:"bid_amount"`
}

type BidResponse struct {
	AuctionID string    `json:"auction_id"`
	BidderID  string    `json:"bidder_id"`
	BidAmount float64   `json:"amount"`
	TimeStamp time.Time `json:"created_at"`
}
