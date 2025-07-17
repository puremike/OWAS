package config

import (
	"time"

	"github.com/puremike/online_auction_api/internal/auth"
	"github.com/puremike/online_auction_api/internal/payments"
	"github.com/puremike/online_auction_api/internal/ratelimiters"
	"github.com/puremike/online_auction_api/internal/store"
	"github.com/puremike/online_auction_api/internal/store/cache"

	"github.com/puremike/online_auction_api/internal/ws"
	"github.com/puremike/online_auction_api/pkg"
	"go.uber.org/zap"
)

type Application struct {
	AppConfig *AppConfig
	Logger    *zap.SugaredLogger
	JwtAUth   *auth.JWTAuthenticator
	Store     *store.Storage
	WsHub     *ws.Hub
	// RateLimiter *ratelimiters.HybridLimiter
	GeneralRateLimiter   ratelimiters.Limiter
	SensitiveRateLimiter ratelimiters.Limiter
	HeavyOpsRateLimiter  ratelimiters.Limiter
	Stripe               *payments.StripePayment
	RedisCache           *cache.Storage
}

type AppConfig struct {
	Port           string
	Env            string
	DbConfig       DbConfig
	AuthConfig     AuthConfig
	GeneralRL      RateLimiterConf
	SensitiveRL    RateLimiterConf
	HeavyOpsRL     RateLimiterConf
	StripeConf     StripeConf
	S3Bucket       string
	RedisCacheConf RedisCacheConf
}

type RedisCacheConf struct {
	Addr     string
	Password string
	DB       int
	Enabled  bool
}

type StripeConf struct {
	StripeSecretKey string
	CancelURL       string
	SuccessURL      string
}

type RateLimiterConf struct {
	Window   time.Duration
	Limit    int
	Rate     float64
	Capacity float64
	Enabled  bool
}

type DbConfig struct {
	Db_addr          string
	MaxIdleConns     int
	MaxOpenConns     int
	ConnsMaxIdleTime time.Duration
}

type AuthConfig struct {
	Aud, Iss, Secret string
	TokenExp         time.Duration
	RefreshTokenExp  time.Duration
}

const ApiVersion = "1.0.1"

func MyCfg() *AppConfig {
	return &AppConfig{
		RedisCacheConf: RedisCacheConf{
			Addr:     pkg.GetEnvString("REDIS_ADDR", "localhost:6379"),
			Password: pkg.GetEnvString("REDIS_PASSWORD", ""),
			DB:       pkg.GetEnvInt("REDIS_DB", 0),
			Enabled:  pkg.GetEnvBool("REDIS_ENABLED", false),
		},
		Port:     pkg.GetEnvString("PORT", "8080"),
		Env:      pkg.GetEnvString("ENV", "development"),
		S3Bucket: pkg.GetEnvString("S3_BUCKET", ""),
		DbConfig: DbConfig{
			Db_addr:          pkg.GetEnvString("DB_ADDR", "postgres://user:userpassword@localhost:5432/OWAS?sslmode=disable"),
			MaxIdleConns:     pkg.GetEnvInt("DB_MAX_IDLE_CONNS", 5),
			MaxOpenConns:     pkg.GetEnvInt("DB_MAX_OPEN_CONNS", 50),
			ConnsMaxIdleTime: pkg.GetEnvTDuration("DB_CONNS_MAX_IDLE_TIME", 30*time.Minute),
		},
		AuthConfig: AuthConfig{
			Aud:             pkg.GetEnvString("JWT_AUD", "OWAS"),
			Iss:             pkg.GetEnvString("JWT_ISS", "OWAS"),
			Secret:          pkg.GetEnvString("JWT_SECRET", "cb02ad3d42d1818c330c8a3d78f88d2b613b75cd99e56cf2182b9ad5b0c39ef20ddb7e87144d9c94dd212b2f08349c6a090eadd42736f4335a76f482e4f6762a"),
			TokenExp:        pkg.GetEnvTDuration("JWT_TOKEN_EXP", 30*time.Minute),
			RefreshTokenExp: pkg.GetEnvTDuration("JWT_REFRESH_TOKEN_EXP", 7*24*time.Hour),
		},
		GeneralRL: RateLimiterConf{
			Enabled:  pkg.GetEnvBool("GEN_RATE_LIMITER_ENABLED", false),
			Limit:    pkg.GetEnvInt("GEN_RATE_LIMITER_LIMIT", 10),
			Window:   pkg.GetEnvTDuration("GEN_RATE_LIMITER_WINDOW", 1*time.Minute),
			Rate:     pkg.GetEnvFloat("GEN_RATE_LIMITER_RATE", 1.666),
			Capacity: pkg.GetEnvFloat("GEN_RATE_LIMITER_CAPACITY", 100),
		},
		SensitiveRL: RateLimiterConf{
			Enabled:  pkg.GetEnvBool("SEN_RATE_LIMITER_ENABLED", false),
			Limit:    pkg.GetEnvInt("SEN_RATE_LIMITER_LIMIT", 5),
			Window:   pkg.GetEnvTDuration("SEN_RATE_LIMITER_WINDOW", 1*time.Minute),
			Rate:     pkg.GetEnvFloat("SEN_RATE_LIMITER_RATE", 0.083),
			Capacity: pkg.GetEnvFloat("SEN_RATE_LIMITER_CAPACITY", 5),
		},
		HeavyOpsRL: RateLimiterConf{
			Enabled:  pkg.GetEnvBool("HEA_RATE_LIMITER_ENABLED", false),
			Limit:    pkg.GetEnvInt("HEA_RATE_LIMITER_LIMIT", 2),
			Window:   pkg.GetEnvTDuration("HEA_RATE_LIMITER_WINDOW", 5*time.Minute),
			Rate:     pkg.GetEnvFloat("HEA_RATE_LIMITER_RATE", 0.00666),
			Capacity: pkg.GetEnvFloat("HEA_RATE_LIMITER_CAPACITY", 2),
		},

		StripeConf: StripeConf{
			StripeSecretKey: pkg.GetEnvString("STRIPE_SECRET_KEY", ""),
			CancelURL:       pkg.GetEnvString("STRIPE_CANCEL_URL", ""),
			SuccessURL:      pkg.GetEnvString("STRIPE_SUCCESS_URL", ""),
		},
	}
}

func MyRateLimiters(cfg *AppConfig) (gLm, sLm, hLm ratelimiters.Limiter) {

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
