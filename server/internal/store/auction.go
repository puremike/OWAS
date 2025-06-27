package store

import (
	"context"
	"database/sql"

	"github.com/puremike/online_auction_api/internal/errs"
	"github.com/puremike/online_auction_api/internal/models"
)

type AuctionStore struct {
	db *sql.DB
}

func (a *AuctionStore) GetAuctionById(ctx context.Context, id string) (*models.Auction, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	auction := &models.Auction{}

	query := `SELECT id, seller_id, winner_id, title, description, starting_price, current_price, type, status, start_time, end_time, created_at FROM auctions WHERE id = $1`

	if err := a.db.QueryRowContext(ctx, query, id).Scan(&auction.ID, &auction.SellerID, &auction.WinnerID, &auction.Title, &auction.Description, &auction.StartingPrice, &auction.CurrentPrice, &auction.Type, &auction.Status, &auction.StartTime, &auction.EndTime, &auction.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrAuctionNotFound
		}
		return nil, err
	}

	return auction, nil
}

func (a *AuctionStore) CloseAuction(ctx context.Context, status, id string) error {

	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	query := `UPDATE auctions SET status = $1 WHERE id = $2`

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, query, status, id); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (a *AuctionStore) GetAuctions(ctx context.Context) (*[]models.Auction, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	var auctions []models.Auction

	query := `SELECT id, seller_id, title, description, starting_price, current_price, type, status, start_time, end_time, created_at FROM auctions`

	rows, err := a.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var a models.Auction

		if err := rows.Scan(&a.ID, &a.SellerID, &a.Title, &a.Description, &a.StartingPrice, &a.CurrentPrice, &a.Type, &a.Status, &a.StartTime, &a.EndTime, &a.CreatedAt); err != nil {
			return nil, err
		}

		auctions = append(auctions, a)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &auctions, nil
}

func (a *AuctionStore) CreateAuction(ctx context.Context, auction *models.Auction) (*models.Auction, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	query := `INSERT INTO auctions (seller_id, winner_id, title, description, starting_price, current_price, type, status, start_time, end_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, seller_id, winner_id, title, description, starting_price, current_price, type, status, start_time, end_time, created_at`

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	if err = tx.QueryRowContext(ctx, query, auction.SellerID, auction.WinnerID, auction.Title, auction.Description, auction.StartingPrice, auction.CurrentPrice, auction.Type, auction.Status, auction.StartTime, auction.EndTime).Scan(&auction.ID, &auction.SellerID, &auction.WinnerID, &auction.Title, &auction.Description, &auction.StartingPrice, &auction.CurrentPrice, &auction.Type, &auction.Status, &auction.StartTime, &auction.EndTime, &auction.CreatedAt); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return auction, nil
}

func (a *AuctionStore) UpdateAuction(ctx context.Context, auction *models.Auction, id string) error {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	query := `UPDATE auctions SET seller_id = $1, title = $2, description = $3, starting_price = $4, current_price = $5, type = $6, status = $7, start_time = $8, end_time = $9, winner_id = $10 WHERE id = $11`

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx, query, auction.SellerID, auction.Title, auction.Description, auction.StartingPrice, auction.CurrentPrice, auction.Type, auction.Status, auction.StartTime, auction.EndTime, auction.WinnerID, id); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (a *AuctionStore) DeleteAuction(ctx context.Context, id string) error {

	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	query := `DELETE FROM auctions WHERE id = $1`

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, query, id); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
