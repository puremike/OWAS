package main

import (
	"time"

	"github.com/puremike/online_auction_api/internal/db"
	"github.com/puremike/online_auction_api/internal/env"
	"go.uber.org/zap"
)

type application struct {
	config config
	logger *zap.SugaredLogger
}

type config struct {
	port     string
	env      string
	dbconfig dbconfig
}

type dbconfig struct {
	db_addr          string
	maxIdleConns     int
	maxOpenConns     int
	connsMaxIdleTime time.Duration
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
			dbconfig: dbconfig{
				db_addr:          env.GetEnvString("DB_ADDR", "postgres://admin:adminpassword123@localhost:5432/OWAS?sslmode=disable"),
				maxIdleConns:     env.GetEnvInt("DB_MAX_IDLE_CONNS", 5),
				maxOpenConns:     env.GetEnvInt("DB_MAX_OPEN_CONNS", 50),
				connsMaxIdleTime: env.GetEnvTDuration("DB_CONNS_MAX_IDLE_TIME", 30*time.Minute),
			},
		},
	}

	logger := zap.NewExample().Sugar()
	defer logger.Sync()

	db, err := db.NewPostgresDB(cfg.config.dbconfig.db_addr, cfg.config.dbconfig.maxIdleConns, cfg.config.dbconfig.maxOpenConns, cfg.config.dbconfig.connsMaxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Infow("Connected to database successfully")

	app := &application{
		config: cfg.config,
		logger: logger,
	}

	mux := app.routes()
	logger.Fatal(app.server(mux))
}
