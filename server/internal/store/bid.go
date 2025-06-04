package store

import "database/sql"

type BidStore struct {
	db *sql.DB
}
