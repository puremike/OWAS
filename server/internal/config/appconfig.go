package config

import (
	"time"

	"github.com/puremike/online_auction_api/internal/auth"
	"github.com/puremike/online_auction_api/internal/payments"
	"github.com/puremike/online_auction_api/internal/ratelimiters"
	"github.com/puremike/online_auction_api/internal/store"
	"github.com/puremike/online_auction_api/internal/ws"
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
}

type AppConfig struct {
	Port        string
	Env         string
	DbConfig    DbConfig
	AuthConfig  AuthConfig
	GeneralRL   RateLimiterConf
	SensitiveRL RateLimiterConf
	HeavyOpsRL  RateLimiterConf
	StripeConf  StripeConf
	S3Bucket    string
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
