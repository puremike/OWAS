package main

import (
	"time"

	"github.com/puremike/online_auction_api/internal/auth"
	"github.com/puremike/online_auction_api/internal/db"
	"github.com/puremike/online_auction_api/internal/env"
	"go.uber.org/zap"
)

type application struct {
	config  config
	logger  *zap.SugaredLogger
	jwtAUth *auth.JWTAuthenticator
}

type config struct {
	port       string
	env        string
	dbConfig   dbConfig
	authConfig authConfig
}

type dbConfig struct {
	db_addr          string
	maxIdleConns     int
	maxOpenConns     int
	connsMaxIdleTime time.Duration
}

type authConfig struct {
	aud, iss, secret string
}

const apiVersion = "1.0.0"

// @title					Online Webbased Auction System API
// @version					1.0.0
// @description				This is an API for a Online Webbased Auction System
//
// @contact.name				Puremike
// @contact.url				http://github.com/puremike
// @contact.email				digitalmarketfy@gmail.com
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//
// @BasePath					/api/v1
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Use a valid JWT token. Format: Bearer <token>
func main() {

	cfg := &application{
		config: config{
			port: env.GetEnvString("PORT", "6000"),
			env:  env.GetEnvString("ENV", "development"),
			dbConfig: dbConfig{
				db_addr:          env.GetEnvString("DB_ADDR", "postgres://admin:adminpassword123@localhost:5432/OWAS?sslmode=disable"),
				maxIdleConns:     env.GetEnvInt("DB_MAX_IDLE_CONNS", 5),
				maxOpenConns:     env.GetEnvInt("DB_MAX_OPEN_CONNS", 50),
				connsMaxIdleTime: env.GetEnvTDuration("DB_CONNS_MAX_IDLE_TIME", 30*time.Minute),
			},
			authConfig: authConfig{
				aud:    env.GetEnvString("JWT_AUD", "OWAS"),
				iss:    env.GetEnvString("JWT_ISS", "OWAS"),
				secret: env.GetEnvString("JWT_SECRET", "cb02ad3d42d1818c330c8a3d78f88d2b613b75cd99e56cf2182b9ad5b0c39ef20ddb7e87144d9c94dd212b2f08349c6a090eadd42736f4335a76f482e4f6762a"),
			},
		},
	}

	logger := zap.NewExample().Sugar()
	defer logger.Sync()

	db, err := db.NewPostgresDB(cfg.config.dbConfig.db_addr, cfg.config.dbConfig.maxIdleConns, cfg.config.dbConfig.maxOpenConns, cfg.config.dbConfig.connsMaxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Infow("Connected to database successfully")

	app := &application{
		config: cfg.config,
		logger: logger,
		jwtAUth: auth.NewJWTAuthenticator(
			cfg.config.authConfig.secret, cfg.config.authConfig.iss, cfg.config.authConfig.aud),
	}

	mux := app.routes()
	logger.Fatal(app.server(mux))
}
