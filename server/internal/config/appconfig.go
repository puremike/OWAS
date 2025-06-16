package config

import (
	"time"

	"github.com/puremike/online_auction_api/internal/auth"
	"github.com/puremike/online_auction_api/internal/store"
	"github.com/puremike/online_auction_api/internal/ws"
	"go.uber.org/zap"
)

type Application struct {
	AppConfig AppConfig
	Logger    *zap.SugaredLogger
	JwtAUth   *auth.JWTAuthenticator
	Store     *store.Storage
	WsHub     *ws.Hub
}

type AppConfig struct {
	Port       string
	Env        string
	DbConfig   DbConfig
	AuthConfig AuthConfig
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
