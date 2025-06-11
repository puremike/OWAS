package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/puremike/online_auction_api/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	StoreRefreshToken(ctx context.Context, userID, refreshToken string, expires_at time.Time) error
	UpdateUser(ctx context.Context, user *models.User, id string) error
	ValidateRefreshToken(ctx context.Context, refreshToken string) (string, error)
	ChangePassword(ctx context.Context, pass, id string) error
}

type AuctionRepository interface {
	GetAuctionById(ctx context.Context, id string) (*models.Auction, error)
	GetAuctions(ctx context.Context) (*[]models.Auction, error)
	CreateAuction(ctx context.Context, auction *models.Auction) (*models.Auction, error)
	UpdateAuction(ctx context.Context, auction *models.Auction, id string) error
	DeleteAuction(ctx context.Context, id string) error
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
)
