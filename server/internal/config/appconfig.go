package config

import (
	"time"

	"github.com/puremike/online_auction_api/internal/auth"
	"github.com/puremike/online_auction_api/internal/store"
	"go.uber.org/zap"
)

type Application struct {
	AppConfig AppConfig
	Logger    *zap.SugaredLogger
	JwtAUth   *auth.JWTAuthenticator
	Store     *store.Storage
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
}

const ApiVersion = "1.0.1"
