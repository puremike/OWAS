package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/puremike/online_auction_api/internal/models"
)

type Users interface {
	Get(ctx context.Context, id string) (*models.User, error)
	Set(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id string) error
}

type Auctions interface {
	Get(ctx context.Context, id string) (*models.Auction, error)
	Set(ctx context.Context, user *models.Auction) error
	Delete(ctx context.Context, id string) error
}

type Bids interface {
	Get(ctx context.Context, id string) (*models.Bid, error)
	Set(ctx context.Context, user *models.Bid) error
	Delete(ctx context.Context, auctionID string) error
}

type Storage struct {
	Users    Users
	Auctions Auctions
	Bids     Bids
}

func NewRDBCacheStorage(rdb *redis.Client) *Storage {
	return &Storage{
		Users:    &UserCacheStore{rdb},
		Auctions: &AuctionCacheStore{rdb},
		// Bids:     &BidCacheStore{rdb},
	}
}
