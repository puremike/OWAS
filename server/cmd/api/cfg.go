package main

import (
	"time"

	"github.com/puremike/online_auction_api/internal/config"
	"github.com/puremike/online_auction_api/pkg"
)

func myCfg() *config.AppConfig {
	return &config.AppConfig{
		Port: pkg.GetEnvString("PORT", "8080"),
		Env:  pkg.GetEnvString("ENV", "development"),
		DbConfig: config.DbConfig{
			Db_addr:          pkg.GetEnvString("DB_ADDR", "postgres://user:userpassword@localhost:5432/OWAS?sslmode=disable"),
			MaxIdleConns:     pkg.GetEnvInt("DB_MAX_IDLE_CONNS", 5),
			MaxOpenConns:     pkg.GetEnvInt("DB_MAX_OPEN_CONNS", 50),
			ConnsMaxIdleTime: pkg.GetEnvTDuration("DB_CONNS_MAX_IDLE_TIME", 30*time.Minute),
		},
		AuthConfig: config.AuthConfig{
			Aud:             pkg.GetEnvString("JWT_AUD", "OWAS"),
			Iss:             pkg.GetEnvString("JWT_ISS", "OWAS"),
			Secret:          pkg.GetEnvString("JWT_SECRET", "cb02ad3d42d1818c330c8a3d78f88d2b613b75cd99e56cf2182b9ad5b0c39ef20ddb7e87144d9c94dd212b2f08349c6a090eadd42736f4335a76f482e4f6762a"),
			TokenExp:        pkg.GetEnvTDuration("JWT_TOKEN_EXP", 30*time.Minute),
			RefreshTokenExp: pkg.GetEnvTDuration("JWT_REFRESH_TOKEN_EXP", 7*24*time.Hour),
		},
		GeneralRL: config.RateLimiterConf{
			Enabled:  pkg.GetEnvBool("GEN_RATE_LIMITER_ENABLED", false),
			Limit:    pkg.GetEnvInt("GEN_RATE_LIMITER_LIMIT", 10),
			Window:   pkg.GetEnvTDuration("GEN_RATE_LIMITER_WINDOW", 1*time.Minute),
			Rate:     pkg.GetEnvFloat("GEN_RATE_LIMITER_RATE", 1.666),
			Capacity: pkg.GetEnvFloat("GEN_RATE_LIMITER_CAPACITY", 100),
		},
		SensitiveRL: config.RateLimiterConf{
			Enabled:  pkg.GetEnvBool("SEN_RATE_LIMITER_ENABLED", false),
			Limit:    pkg.GetEnvInt("SEN_RATE_LIMITER_LIMIT", 5),
			Window:   pkg.GetEnvTDuration("SEN_RATE_LIMITER_WINDOW", 1*time.Minute),
			Rate:     pkg.GetEnvFloat("SEN_RATE_LIMITER_RATE", 0.083),
			Capacity: pkg.GetEnvFloat("SEN_RATE_LIMITER_CAPACITY", 5),
		},
		HeavyOpsRL: config.RateLimiterConf{
			Enabled:  pkg.GetEnvBool("HEA_RATE_LIMITER_ENABLED", false),
			Limit:    pkg.GetEnvInt("HEA_RATE_LIMITER_LIMIT", 2),
			Window:   pkg.GetEnvTDuration("HEA_RATE_LIMITER_WINDOW", 5*time.Minute),
			Rate:     pkg.GetEnvFloat("HEA_RATE_LIMITER_RATE", 0.00666),
			Capacity: pkg.GetEnvFloat("HEA_RATE_LIMITER_CAPACITY", 2),
		},

		StripeConf: config.StripeConf{
			StripeSecretKey: pkg.GetEnvString("STRIPE_SECRET_KEY", ""),
			CancelURL:       pkg.GetEnvString("STRIPE_CANCEL_URL", ""),
			SuccessURL:      pkg.GetEnvString("STRIPE_SUCCESS_URL", ""),
		},
	}
}
