package cached

import (
	"context"
	"errors"
	"net/http"

	"github.com/puremike/online_auction_api/internal/config"
	"github.com/puremike/online_auction_api/internal/errs"
	"github.com/puremike/online_auction_api/internal/models"
)

type UserCached struct {
	app *config.Application
}

func (m *UserCached) reUser(ctx context.Context, userId string) (*models.User, error) {
	user, err := m.app.Store.Users.GetUserById(ctx, userId)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			m.app.Logger.Errorw("user not found", "userId", userId)
			return nil, errs.ErrUserNotFound
		}
		return nil, errs.NewHTTPError("failed to retrieve user from database", http.StatusInternalServerError)
	}

	return user, nil
}

func (m *UserCached) GetUserFromCache(ctx context.Context, userId string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	// Redis is not enabled, fetch from database
	if !m.app.AppConfig.RedisCacheConf.Enabled {
		user, err := m.reUser(ctx, userId)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	// Redis is enabled
	user, err := m.app.RedisCache.Users.Get(ctx, userId)
	if err != nil {
		m.app.Logger.Errorw("failed to get user from cache", "error", err)
		return nil, errs.NewHTTPError("failed to retrieve user from cache", http.StatusInternalServerError)
	}

	if user != nil {
		m.app.Logger.Infow("cache hit", "key", userId)
		return user, nil
	}

	// Redis cache miss
	// Fetch from database
	m.app.Logger.Infow("cache miss", "key", userId)
	user, err = m.reUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	if err := m.app.RedisCache.Users.Set(ctx, user); err != nil {
		m.app.Logger.Errorw("failed to set user in cache", "error", err)
		return nil, errs.NewHTTPError("failed to set user in cache", http.StatusInternalServerError)
	}

	return user, nil
}

func (m *UserCached) DeleteUserFromCache(ctx context.Context, userId string) error {
	ctx, cancel := context.WithTimeout(ctx, DefaultTimeout)
	defer cancel()

	// Delete user from database
	if err := m.app.Store.Users.DeleteUser(ctx, userId); err != nil {
		m.app.Logger.Errorw("failed to delete user from db", "error", err)
		return errs.ErrUserNotFound
	}

	if !m.app.AppConfig.RedisCacheConf.Enabled {
		return nil
	}

	// Redis is enabled
	// Delete user from Redis cache
	if err := m.app.RedisCache.Users.Delete(ctx, userId); err != nil {
		m.app.Logger.Errorw("failed to delete user from cache", "error", err)
		return errs.NewHTTPError("failed to delete user from cache", http.StatusInternalServerError)
	}

	return nil
}
