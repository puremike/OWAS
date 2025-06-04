package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/puremike/online_auction_api/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
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
	Users         UserRepository
	Auctions      AuctionRepository
	Bids          BidRepository
	Payments      PaymentRepository
	Notifications NotificationRepository
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Users:         &UserStore{db},
		Auctions:      &AuctionStore{db},
		Bids:          &BidStore{db},
		Payments:      &PaymentStore{db},
		Notifications: &NotificationStore{db},
	}
}

var (
	QueryBackgroundTimeout = 5 * time.Second
	ErrUserNotFound        = errors.New("user not found")
)
