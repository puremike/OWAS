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
	NotificationAuctionEnded NotificationUpdateType = "AUCTION_ENDED"
)

type NotificationEvent struct {
	Type      NotificationUpdateType `json:"type"`
	UserID    string                 `json:"user_id"` // FK to User.ID
	Message   string                 `json:"message"`
	AuctionID string                 `json:"auction_id"`
	TimeStamp time.Time              `json:"created_at"`
}

const (
	EnglishAuction = "english"
	DutchAuction   = "dutch"
	SealedAuction  = "sealed"
)

type WinnerResponse struct {
	WinnerID   string  `json:"winner_id"`
	WinningBid float64 `json:"winning_bid"`
	Status     string  `json:"status"`
}
