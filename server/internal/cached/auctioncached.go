package cached

import (
	"context"
	"errors"
	"net/http"

	"github.com/puremike/online_auction_api/internal/config"
	"github.com/puremike/online_auction_api/internal/errs"
	"github.com/puremike/online_auction_api/internal/models"
)

type AuctionCached struct {
	app *config.Application
}

func (m *AuctionCached) reAuction(ctx context.Context, auctionId string) (*models.Auction, error) {

	auction, err := m.app.Store.Auctions.GetAuctionById(ctx, auctionId)
	if err != nil {
		if errors.Is(err, errs.ErrAuctionNotFound) {
			m.app.Logger.Errorw("auction not found", "auctionId", auctionId)
			return nil, errs.ErrAuctionNotFound
		}
		return nil, errs.NewHTTPError("failed to retrieve auction from database", http.StatusInternalServerError)
	}

	return auction, nil

}
func (m *AuctionCached) GetAuctionFromCache(ctx context.Context, auctionId string) (*models.Auction, error) {

	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	if !m.app.AppConfig.RedisCacheConf.Enabled {
		auction, err := m.reAuction(ctx, auctionId)
		if err != nil {
			return nil, err
		}
		return auction, nil
	}

	// Redis is enabled
	auction, err := m.app.RedisCache.Auctions.Get(ctx, auctionId)
	if err != nil {
		m.app.Logger.Errorw("failed to get auction from cache", "error", err)
		return nil, errs.NewHTTPError("failed to retrieve auction from cache", http.StatusInternalServerError)
	}

	if auction != nil {
		m.app.Logger.Infow("cache hit", "key", auctionId)
		return auction, nil
	}

	// Cache miss
	m.app.Logger.Infow("cache miss", "key", auctionId)

	// Fetch from database if not found in cache
	auction, err = m.reAuction(ctx, auctionId)
	if err != nil {
		return nil, err
	}

	// Set auction in cache
	m.app.Logger.Infow("setting auction in cache", "auctionId", auctionId)

	if err = m.app.RedisCache.Auctions.Set(ctx, auction); err != nil {
		m.app.Logger.Errorw("failed to set auction in cache", "error", err)
		return nil, errs.NewHTTPError("failed to set auction in cache", http.StatusInternalServerError)
	}

	return auction, nil
}

func (m *AuctionCached) DeleteAuctionFromCache(ctx context.Context, auctionId string) error {
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	// Delete user from database
	if err := m.app.Store.Auctions.DeleteAuction(ctx, auctionId); err != nil {
		m.app.Logger.Errorw("failed to delete auction from db", "error", err)
		return errs.ErrAuctionNotFound
	}

	if !m.app.AppConfig.RedisCacheConf.Enabled {
		return nil
	}

	// Redis is enabled
	// Delete user from Redis cache
	if err := m.app.RedisCache.Auctions.Delete(ctx, auctionId); err != nil {
		m.app.Logger.Errorw("failed to delete auction from cache", "error", err)
		return errs.NewHTTPError("failed to delete auction from cache", http.StatusInternalServerError)
	}
	return nil
}
