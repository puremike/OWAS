package store

import "database/sql"

type AuctionStore struct {
	db *sql.DB
}
