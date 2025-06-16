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

	query := `INSERT INTO bids (id, auction_id, bidder_id, amount) VALUES ($1, $2, $3, $4, $5) RETURNING id, auction_id, bidder_id, amount, created_at`

	tx, err := b.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	if err = tx.QueryRowContext(ctx, query, bid.ID, bid.AuctionID, bid.BidderID, bid.Amount).Scan(&bid.ID, &bid.AuctionID, &bid.BidderID, &bid.Amount, &bid.CreatedAt); err != nil {
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

	query := `SELECT id, auction_id, user_id, amount, created_at FROM bids WHERE user_id = $1`

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
	query := `SELECT id, auction_id, bidder_id, amount, created_at FROM bids WHERE auction_id = $1 ORDER BY amount DESC LIMIT 1`
	var bid models.Bid
	err := b.db.QueryRowContext(ctx, query, id).Scan(&bid.ID, &bid.AuctionID, &bid.BidderID, &bid.Amount, &bid.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &bid, nil
}

func (b *BidStore) GetBidById(ctx context.Context, id string) (*models.Bid, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	bid := &models.Bid{}

	query := `SELECT id, auction_id, user_id, amount, created_at FROM bids WHERE id = $1`

	if err := b.db.QueryRowContext(ctx, query, id).Scan(&bid.ID, &bid.AuctionID, &bid.BidderID, &bid.Amount, &bid.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrBidNotFound
		}
		return nil, err
	}

	return bid, nil
}

func (b *BidStore) GetAllBidderIDsForAuction(ctx context.Context, id string) ([]string, error) {
	query := `SELECT DISTINCT bidder_id FROM bids WHERE auction_id = $1`
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
