package main

import (
	_ "github.com/lib/pq"
	"github.com/puremike/online_auction_api/docs"
	"github.com/puremike/online_auction_api/internal/auth"
	"github.com/puremike/online_auction_api/internal/config"
	"github.com/puremike/online_auction_api/internal/db"
	"github.com/puremike/online_auction_api/internal/routes"
	"github.com/puremike/online_auction_api/internal/store"
	"github.com/puremike/online_auction_api/internal/ws"
	"go.uber.org/zap"
)

//	@title			Online Webbased Auction System API
//	@version		1.0.0
//	@description	This is an API for a Online Webbased Auction System
//
//	@contact.name	Puremike
//	@contact.url	http://github.com/puremike
//	@contact.email	digitalmarketfy@gmail.com
//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html
//
//	@BasePath		/api/v1

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Use a valid JWT token. Format: Bearer <token>

//	@securityDefinitions.apikey	jwtCookieAuth
//	@type						apiKey
//	@in							cookie
//	@name						jwt
//	@description				JWT (JSON Web Token) access token, sent as an HttpOnly cookie.

// @securityDefinitions.apikey	refreshTokenCookie
// @type						apiKey
// @in							cookie
// @name						refresh_token
// @description				Refresh token, sent as an HttpOnly cookie.
func main() {

	docs.SwaggerInfo.BasePath = "/api/v1"

	// configuration
	cfg := myCfg()

	logger := zap.NewExample().Sugar()
	defer logger.Sync()

	db, err := db.NewPostgresDB(cfg.DbConfig.Db_addr, cfg.DbConfig.MaxIdleConns, cfg.DbConfig.MaxOpenConns, cfg.DbConfig.ConnsMaxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Infow("Connected to database successfully")

	gLm, sLm, hLm := myRateLimiters(cfg)

	app := &config.Application{
		AppConfig: cfg,
		Logger:    logger,
		JwtAUth: auth.NewJWTAuthenticator(
			cfg.AuthConfig.Secret, cfg.AuthConfig.Iss, cfg.AuthConfig.Aud),
		Store:                store.NewStorage(db),
		WsHub:                ws.NewHub(),
		GeneralRateLimiter:   gLm,
		SensitiveRateLimiter: sLm,
		HeavyOpsRateLimiter:  hLm,
	}

	go app.WsHub.Run()

	mux := routes.Routes(app)
	logger.Fatal(routes.RunServer(mux, cfg.Port, logger))
}
