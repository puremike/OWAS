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
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	StoreRefreshToken(ctx context.Context, userID, refreshToken string, expires_at time.Time) error
	ValidateRefreshToken(ctx context.Context, refreshToken string) (string, error)
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
	ErrTokenNotFound       = errors.New("token not found")
)
