package store

import (
	"context"
	"database/sql"

	"github.com/puremike/online_auction_api/internal/errs"
	"github.com/puremike/online_auction_api/internal/models"
)

type BidStore struct {
	db *sql.DB
}

func (b *BidStore) CreateBid(ctx context.Context, bid *models.Bid) (*models.Bid, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	query := `INSERT INTO bid (auction_id, bidder_id, amount) VALUES ($1, $2, $3) RETURNING id, auction_id, bidder_id, amount, created_at`

	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	if err = tx.QueryRowContext(ctx, query, bid.AuctionID, bid.BidderID, bid.Amount).Scan(&bid.ID, &bid.AuctionID, &bid.BidderID, &bid.Amount, &bid.CreatedAt); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return bid, nil
}

func (b *BidStore) GetBids(ctx context.Context, userId string) (*[]models.Bid, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	var bids []models.Bid

	query := `SELECT id, auction_id, user_id, amount, created_at FROM bid WHERE user_id = $1`

	rows, err := b.db.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var b models.Bid

		if err := rows.Scan(&b.ID, &b.AuctionID, &b.BidderID, &b.Amount, &b.CreatedAt); err != nil {
			return nil, err
		}

		bids = append(bids, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &bids, nil

}

func (b *BidStore) GetHighestBid(ctx context.Context, id string) (*models.Bid, error) {
	query := `SELECT id, auction_id, bidder_id, amount, created_at FROM bid WHERE auction_id = $1 ORDER BY amount DESC LIMIT 1`
	var bid models.Bid
	err := b.db.QueryRowContext(ctx, query, id).Scan(&bid.ID, &bid.AuctionID, &bid.BidderID, &bid.Amount, &bid.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrBidNotFound
		}
		return nil, err
	}
	return &bid, nil
}

func (b *BidStore) GetBidById(ctx context.Context, id string) (*models.Bid, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	bid := &models.Bid{}

	query := `SELECT id, auction_id, user_id, amount, created_at FROM bid WHERE id = $1`

	if err := b.db.QueryRowContext(ctx, query, id).Scan(&bid.ID, &bid.AuctionID, &bid.BidderID, &bid.Amount, &bid.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrBidNotFound
		}
		return nil, err
	}

	return bid, nil
}

func (b *BidStore) GetBidByUser(ctx context.Context, auctionID, bidderID string) (*models.Bid, error) {

	query := `SELECT * FROM bid WHERE auction_id = $1 AND bidder_id = $2 LIMIT 1`

	bid := &models.Bid{}
	if err := b.db.QueryRowContext(ctx, query, auctionID, bidderID).Scan(&bid.ID, &bid.AuctionID, &bid.BidderID, &bid.Amount, &bid.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrBidNotFound
		}
		return nil, err
	}
	return bid, nil
}

func (b *BidStore) GetAllBidderIDsForAuction(ctx context.Context, id string) ([]string, error) {
	query := `SELECT DISTINCT bidder_id FROM bid WHERE auction_id = $1`
	rows, err := b.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bidders []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		bidders = append(bidders, id)
	}
	return bidders, nil
}

func (b *BidStore) DeleteBidsByAuction(ctx context.Context, auctionID string) error {
	query := `DELETE FROM bid WHERE auction_id = $1`

	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if _, err = b.db.ExecContext(ctx, query, auctionID); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
