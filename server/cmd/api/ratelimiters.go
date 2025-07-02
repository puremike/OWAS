package main

import (
	"github.com/puremike/online_auction_api/internal/config"
	"github.com/puremike/online_auction_api/internal/ratelimiters"
)

func MyRateLimiters(cfg *config.AppConfig) (gLm, sLm, hLm ratelimiters.Limiter) {

	var generalRL ratelimiters.Limiter
	if cfg.GeneralRL.Enabled {
		generalRL = ratelimiters.NewHybridLimiters(ratelimiters.NewSlidingWindow(cfg.GeneralRL.Limit, cfg.GeneralRL.Window), ratelimiters.NewTokenBucket(cfg.GeneralRL.Rate, cfg.GeneralRL.Capacity))
	}

	var sensitiveRL ratelimiters.Limiter
	if cfg.SensitiveRL.Enabled {
		sensitiveRL = ratelimiters.NewHybridLimiters(ratelimiters.NewSlidingWindow(cfg.SensitiveRL.Limit, cfg.SensitiveRL.Window), ratelimiters.NewTokenBucket(cfg.SensitiveRL.Rate, cfg.SensitiveRL.Capacity))
	}

	var heavyOpsRL ratelimiters.Limiter
	if cfg.HeavyOpsRL.Enabled {
		heavyOpsRL = ratelimiters.NewHybridLimiters(ratelimiters.NewSlidingWindow(cfg.HeavyOpsRL.Limit, cfg.HeavyOpsRL.Window), ratelimiters.NewTokenBucket(cfg.HeavyOpsRL.Rate, cfg.HeavyOpsRL.Capacity))
	}

	return generalRL, sensitiveRL, heavyOpsRL
}
