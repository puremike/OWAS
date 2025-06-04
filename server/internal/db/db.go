package db

import (
	"context"
	"database/sql"
	"time"
)

func NewPostgresDB(db_addr string, maxIdleConns, maxOpenConns int, connsMaxIdleTime time.Duration) (*sql.DB, error) {
	db, err := sql.Open("postgres", db_addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetConnMaxIdleTime(connsMaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}
