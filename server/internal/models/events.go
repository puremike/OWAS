package models

import (
	"time"
)

type AuctionUpdateType string

const (
	AuctionNewBid       AuctionUpdateType = "AUCTION_NEW_BID"
	AuctionEnded        AuctionUpdateType = "AUCTION_ENDED"
	AuctionStatusUpdate AuctionUpdateType = "AUCTION_STATUS_UPDATE"
	// You can add more types as your application grows, e.g.:
	// AuctionCancelled  AuctionUpdateType = "AUCTION_CANCELLED"
	// AuctionExtended   AuctionUpdateType = "AUCTION_EXTENDED"
)

type AuctionUpdateEvent struct {
	EventType    AuctionUpdateType `json:"event_type"`
	ID           string            `json:"id"`
	CurrentPrice float64           `json:"current_price"`
	Type         string            `json:"type"`
	Status       string            `json:"status"`
	TimeStamp    time.Time         `json:"start_time"`
	SellerID     string            `json:"seller_id"`
}

type NotificationUpdateType string

const (
	NotificationOutBid       NotificationUpdateType = "OUTBID"
	NotificationWon          NotificationUpdateType = "AUCTION_WON"
	NotificationReminder     NotificationUpdateType = "REMINDER"
	MotificationAuctionEnded NotificationUpdateType = "AUCTION_ENDED"
)

type NotificationEvent struct {
	Type      NotificationUpdateType `json:"type"`
	ID        string                 `json:"id"`
	UserID    string                 `json:"user_id"` // FK to User.ID
	User      User                   `json:"user"`
	Message   string                 `json:"message"`
	AuctionID string                 `json:"auction_id"`
	IsRead    bool                   `json:"is_read"`
	TimeStamp time.Time              `json:"created_at"`
}
