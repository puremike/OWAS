package store

import (
	"context"
	"database/sql"

	"github.com/puremike/online_auction_api/internal/models"
)

type UserStore struct {
	db *sql.DB
}

func (u *UserStore) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	query := `INSERT INTO users (username, email, password, full_name, location, is_admin) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, username, email, full_name, location, created_at`

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	if err = tx.QueryRowContext(ctx, query, user.Username, user.Email, user.Password, user.FullName, user.Location, user.IsAdmin).Scan(&user.ID, &user.Username, &user.Email, &user.FullName, &user.Location, &user.CreatedAt); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}
