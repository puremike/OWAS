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

// Auction represents an auction listing
type Auction struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	StartingPrice float64   `son:"starting_price"`
	CurrentPrice  float64   `json:"current_price"`
	Type          string    `json:"type"`   // "english", "dutch", "sealed"
	Status        string    `json:"status"` // "open", "closed"
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	SellerID      string    `json:"seller_id"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

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

// Notification for user alerts and messages
type Notification struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"` // FK to User.ID
	User      User      `json:"user"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

type Payment struct {
	ID        string    `json:"id"`
	AuctionID string    `json:"auction_id"`
	BuyerID   string    `json:"buyer_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"` // pending, completed, failed
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
