package main

import (
	"time"

	_ "github.com/lib/pq"
	"github.com/puremike/online_auction_api/docs"
	"github.com/puremike/online_auction_api/internal/auth"
	"github.com/puremike/online_auction_api/internal/config"
	"github.com/puremike/online_auction_api/internal/db"
	"github.com/puremike/online_auction_api/internal/routes"
	"github.com/puremike/online_auction_api/internal/store"
	"github.com/puremike/online_auction_api/pkg"
	"go.uber.org/zap"
)

// @title						Online Webbased Auction System API
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

	docs.SwaggerInfo.BasePath = "/api/v1"

	cfg := config.AppConfig{
		Port: pkg.GetEnvString("PORT", "6000"),
		Env:  pkg.GetEnvString("ENV", "development"),
		DbConfig: config.DbConfig{
			Db_addr:          pkg.GetEnvString("DB_ADDR", "postgres://user:userpassword@localhost:5432/OWAS?sslmode=disable"),
			MaxIdleConns:     pkg.GetEnvInt("DB_MAX_IDLE_CONNS", 5),
			MaxOpenConns:     pkg.GetEnvInt("DB_MAX_OPEN_CONNS", 50),
			ConnsMaxIdleTime: pkg.GetEnvTDuration("DB_CONNS_MAX_IDLE_TIME", 30*time.Minute),
		},
		AuthConfig: config.AuthConfig{
			Aud:    pkg.GetEnvString("JWT_AUD", "OWAS"),
			Iss:    pkg.GetEnvString("JWT_ISS", "OWAS"),
			Secret: pkg.GetEnvString("JWT_SECRET", "cb02ad3d42d1818c330c8a3d78f88d2b613b75cd99e56cf2182b9ad5b0c39ef20ddb7e87144d9c94dd212b2f08349c6a090eadd42736f4335a76f482e4f6762a"),
		},
	}

	logger := zap.NewExample().Sugar()
	defer logger.Sync()

	db, err := db.NewPostgresDB(cfg.DbConfig.Db_addr, cfg.DbConfig.MaxIdleConns, cfg.DbConfig.MaxOpenConns, cfg.DbConfig.ConnsMaxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Infow("Connected to database successfully")

	app := &config.Application{
		AppConfig: cfg,
		Logger:    logger,
		JwtAUth: auth.NewJWTAuthenticator(
			cfg.AuthConfig.Secret, cfg.AuthConfig.Iss, cfg.AuthConfig.Aud),
		Store: store.NewStorage(db),
	}

	mux := routes.Routes(app)
	logger.Fatal(routes.RunServer(mux, cfg.Port, logger))
}
