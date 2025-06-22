package store

import (
	"context"
	"database/sql"

	"github.com/puremike/online_auction_api/internal/models"
)

type CSStore struct {
	db *sql.DB
}

func (c *CSStore) ContactSupport(ctx context.Context, cs *models.ContactSupport) (*models.ContactSupport, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	query := `INSERT INTO contact_support (user_id, subject, message) VALUES ($1, $2, $3)`

	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, query, cs.UserID, cs.Subject, cs.Message); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return cs, nil
}
