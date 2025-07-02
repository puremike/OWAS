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
	DeleteUser(ctx context.Context, id string) error
}

type AuctionRepository interface {
	GetAuctionById(ctx context.Context, id string) (*models.Auction, error)
	GetAuctions(ctx context.Context, limit, offset int, filter *models.AuctionFilter) (*[]models.Auction, error)
	CreateAuction(ctx context.Context, auction *models.Auction) (*models.Auction, error)
	CloseAuction(ctx context.Context, status, id string) error
	UpdateAuction(ctx context.Context, auction *models.Auction, id string) error
	DeleteAuction(ctx context.Context, id string) error
	GetWonAuctionsByWinnerID(ctx context.Context, winnerID string) (*[]models.Auction, error)
	UpdateAuctionPaymentStatus(ctx context.Context, isPaid bool, id string) error
	GetBiddedAuctions(ctx context.Context, bidderID string) (*[]models.Auction, error)
	GetAuctionByWinnerId(ctx context.Context, winnerID string) (*models.Auction, error)
	GetAuctionBySellerId(ctx context.Context, sellerID string) (*[]models.Auction, error)
}

type BidRepository interface {
	GetHighestBid(ctx context.Context, id string) (*models.Bid, error)
	GetBidById(ctx context.Context, id string) (*models.Bid, error)
	GetBids(ctx context.Context, userId string) (*[]models.Bid, error)
	CreateBid(ctx context.Context, bid *models.Bid) (*models.Bid, error)
	GetAllBidderIDsForAuction(ctx context.Context, auctionID string) ([]string, error)
	GetBidByUser(ctx context.Context, auctionID, bidderID string) (*models.Bid, error)
	DeleteBidsByAuction(ctx context.Context, auctionID string) error
}

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *models.Payment) error
	GetPayment(ctx context.Context, orderID, buyerID string) (*models.Payment, error)
	UpdatePayment(ctx context.Context, paymentStatus, id string) error
}

type NotificationRepository interface {
	CreateNotification(ctx context.Context, notification *Notification) error
	GetNotifications(ctx context.Context, userID string) ([]*Notification, error)
	DeleteNotificationByAuction(ctx context.Context, auctionID string) error
}

type CSRepository interface {
	ContactSupport(ctx context.Context, cs *models.ContactSupport) (*models.ContactSupport, error)
}

type Storage struct {
	Users         UserRepository
	Auctions      AuctionRepository
	Bids          BidRepository
	Payments      PaymentRepository
	Notifications NotificationRepository
	CS            CSRepository
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		Users:         &UserStore{db},
		Auctions:      &AuctionStore{db},
		Bids:          &BidStore{db},
		Payments:      &PaymentStore{db},
		Notifications: &NotificationStore{db},
		CS:            &CSStore{db},
	}
}

var (
	QueryBackgroundTimeout = 5 * time.Second
)
