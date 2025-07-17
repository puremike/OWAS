package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/puremike/online_auction_api/internal/models"
)

type UserCacheStore struct {
	rdb *redis.Client
}

const timeExp = time.Minute * 2

func (u *UserCacheStore) Get(ctx context.Context, id string) (*models.User, error) {

	cacheKey := "user:" + id
	data, err := u.rdb.Get(ctx, cacheKey).Result()

	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	user := &models.User{}
	if err := json.Unmarshal([]byte(data), user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserCacheStore) Set(ctx context.Context, user *models.User) error {

	cacheKey := "user:" + user.ID
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return u.rdb.SetEX(ctx, cacheKey, data, timeExp).Err()
}

func (u *UserCacheStore) Delete(ctx context.Context, id string) error {
	cacheKey := "user:" + id
	return u.rdb.Del(ctx, cacheKey).Err()
}
