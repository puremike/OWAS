package cache

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/puremike/online_auction_api/internal/models"
)

type AuctionCacheStore struct {
	rdb *redis.Client
}

func (u *AuctionCacheStore) Get(ctx context.Context, id string) (*models.Auction, error) {

	cacheKey := "auction:" + id
	data, err := u.rdb.Get(ctx, cacheKey).Result()

	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	auction := &models.Auction{}
	if err := json.Unmarshal([]byte(data), auction); err != nil {
		return nil, err
	}

	return auction, nil
}

func (u *AuctionCacheStore) Set(ctx context.Context, auction *models.Auction) error {

	cacheKey := "auction:" + auction.ID
	data, err := json.Marshal(auction)
	if err != nil {
		return err
	}

	return u.rdb.SetEX(ctx, cacheKey, data, timeExp).Err()
}

func (u *AuctionCacheStore) Delete(ctx context.Context, id string) error {
	cacheKey := "auction:" + id
	return u.rdb.Del(ctx, cacheKey).Err()
}
