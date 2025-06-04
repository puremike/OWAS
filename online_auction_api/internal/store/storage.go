package store

import "database/sql"

type UserRepository interface {
}

type AuctionRepository interface {
}

type BidRepository interface {
}

type PaymentRepository interface {
}

type NotificationRepository interface {
}

type Storage struct {
	User         UserRepository
	Auction      AuctionRepository
	Bid          BidRepository
	Payment      PaymentRepository
	Notification NotificationRepository
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		User:         &UserStore{db},
		Auction:      &AuctionStore{db},
		Bid:          &BidStore{db},
		Payment:      &PaymentStore{db},
		Notification: &NotificationStore{db},
	}
}
