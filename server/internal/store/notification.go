package store

import (
	"context"
	"database/sql"
)

type Notification struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	Message string `json:"message"`
	IsRead  bool   `json:"is_read"`
}

type NotificationStore struct {
	db *sql.DB
}

func (n *NotificationStore) CreateNotification(ctx context.Context, notification *Notification) error {

	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	query := `INSERT INTO notification (user_id, message, is_read) VALUES ($1, $2, $3)`

	tx, err := n.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, query, notification.UserID, notification.Message, notification.IsRead); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (n *NotificationStore) GetNotifications(ctx context.Context, userID string) ([]*Notification, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	var notifications []*Notification

	query := `SELECT id, user_id, message FROM notification WHERE user_id = $1`

	rows, err := n.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		not := &Notification{}
		if err := rows.Scan(&not.ID, &not.UserID, &not.Message); err != nil {
			return nil, err
		}
		notifications = append(notifications, not)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}
