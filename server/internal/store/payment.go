package store

import "database/sql"

type PaymentStore struct {
	db *sql.DB
}
