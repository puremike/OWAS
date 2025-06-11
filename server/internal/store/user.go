package store

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/puremike/online_auction_api/internal/errs"
	"github.com/puremike/online_auction_api/internal/models"
	"github.com/puremike/online_auction_api/internal/utils"
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

func (u *UserStore) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	user := &models.User{}

	query := `SELECT id, username, email, password, full_name, location, created_at, is_admin FROM users WHERE email = $1`

	if err := u.db.QueryRowContext(ctx, query, strings.ToLower(email)).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.FullName, &user.Location, &user.CreatedAt, &user.IsAdmin); err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (u *UserStore) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	user := &models.User{}

	query := `SELECT id, username, email, password, full_name, location, created_at, is_admin FROM users WHERE username = $1`

	if err := u.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.FullName, &user.Location, &user.CreatedAt, &user.IsAdmin); err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (u *UserStore) GetUserById(ctx context.Context, id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	user := &models.User{}

	query := `SELECT id, username, email, password, full_name, location, created_at, is_admin FROM users WHERE id = $1`

	if err := u.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.FullName, &user.Location, &user.CreatedAt, &user.IsAdmin); err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

func (u *UserStore) UpdateUser(ctx context.Context, user *models.User, id string) error {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	query := `UPDATE users SET username = $1, email = $2, full_name = $3, location = $4 WHERE id = $5`

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx, query, user.Username, user.Email, user.FullName, user.Location, id); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (u *UserStore) ChangePassword(ctx context.Context, pass, id string) error {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	hashedPassword, err := utils.HashedPassword(pass)
	if err != nil {
		return err
	}

	query := `UPDATE users SET password = $1 WHERE id = $2`

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx, query, hashedPassword, id); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (u *UserStore) StoreRefreshToken(ctx context.Context, userID, refreshToken string, expires_at time.Time) error {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	// delete existing refresh tokens for the user
	deleteQuery := `DELETE FROM refresh_tokens WHERE user_id = $1`

	if _, err = tx.ExecContext(ctx, deleteQuery, userID); err != nil {
		return err
	}

	// Insert the new refresh token
	insertQuery := `INSERT INTO refresh_tokens (user_id, token, expires_at) VALUES ($1, $2, $3) `

	if _, err = tx.ExecContext(ctx, insertQuery, userID, refreshToken, expires_at); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (u *UserStore) ValidateRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	var userID string
	var expiresAt time.Time

	query := `SELECT user_id, expires_at FROM refresh_tokens WHERE token = $1`

	if err := u.db.QueryRowContext(ctx, query, refreshToken).Scan(&userID, &expiresAt); err != nil {
		if err == sql.ErrNoRows {
			return "", errs.ErrTokenNotFound
		}
		return "", err
	}
	if time.Now().After(expiresAt) {
		return "", errors.New("refresh token expired")
	}

	return userID, nil
}

func (u *UserStore) GetUsers(ctx context.Context) (*[]models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	var users []models.User

	query := `SELECT id, username, email, password, full_name, location, created_at, is_admin FROM users`

	rows, err := u.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var u models.User

		if err := rows.Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.FullName, &u.Location, &u.CreatedAt, &u.IsAdmin); err != nil {
			return nil, err
		}

		users = append(users, u)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &users, nil
}
