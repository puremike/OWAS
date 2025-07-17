package cached

import (
	"context"
	"time"

	"github.com/puremike/online_auction_api/internal/config"
	"github.com/puremike/online_auction_api/internal/models"
)

type CachedUserInterface interface {
	GetUserFromCache(ctx context.Context, userId string) (*models.User, error)
	DeleteUserFromCache(ctx context.Context, userId string) error
}

type CachedAuctionInterface interface {
	GetAuctionFromCache(ctx context.Context, auctionId string) (*models.Auction, error)
	DeleteAuctionFromCache(ctx context.Context, auctionId string) error
}

type Cached struct {
	User    CachedUserInterface
	Auction CachedAuctionInterface
}

func NewCached(app *config.Application) *Cached {
	return &Cached{
		User:    &UserCached{app},
		Auction: &AuctionCached{app},
	}
}

var (
	DefaultTimeout = 5 * time.Second
)
