package store

import "database/sql"

type NotificationStore struct {
	db *sql.DB
}
