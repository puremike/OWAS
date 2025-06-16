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
	GetUsers(ctx context.Context) (*[]models.User, error)
}

type AuctionRepository interface {
	GetAuctionById(ctx context.Context, id string) (*models.Auction, error)
	GetAuctions(ctx context.Context) (*[]models.Auction, error)
	CreateAuction(ctx context.Context, auction *models.Auction) (*models.Auction, error)
	CloseAuction(ctx context.Context, status, id string) error
	UpdateAuction(ctx context.Context, auction *models.Auction, id string) error
	DeleteAuction(ctx context.Context, id string) error
}

type BidRepository interface {
	GetHighestBid(ctx context.Context, id string) (*models.Bid, error)
	GetBidById(ctx context.Context, id string) (*models.Bid, error)
	GetBids(ctx context.Context, userId string) (*[]models.Bid, error)
	CreateBid(ctx context.Context, bid *models.Bid) (*models.Bid, error)
	GetAllBidderIDsForAuction(ctx context.Context, auctionID string) ([]string, error)
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
