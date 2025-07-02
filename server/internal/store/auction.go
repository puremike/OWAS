package store

import (
	"context"
	"database/sql"
	"log"
	"strconv"

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

	query := `SELECT id, seller_id, winner_id, title, description, starting_price, current_price, type, status, start_time, end_time, image_path, created_at FROM auctions WHERE id = $1`

	if err := a.db.QueryRowContext(ctx, query, id).Scan(&auction.ID, &auction.SellerID, &auction.WinnerID, &auction.Title, &auction.Description, &auction.StartingPrice, &auction.CurrentPrice, &auction.Type, &auction.Status, &auction.StartTime, &auction.EndTime, &auction.ImagePath, &auction.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrAuctionNotFound
		}
		return nil, err
	}

	return auction, nil
}

func (a *AuctionStore) GetAuctionBySellerId(ctx context.Context, sellerID string) (*[]models.Auction, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	auctions := []models.Auction{}

	query := `SELECT id, seller_id, winner_id, title, description, starting_price, current_price, type, status, start_time, end_time, image_path, created_at FROM auctions WHERE seller_id = $1`

	rows, err := a.db.QueryContext(ctx, query, sellerID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var a models.Auction
		err := rows.Scan(&a.ID, &a.SellerID, &a.WinnerID, &a.Title, &a.Description, &a.StartingPrice, &a.CurrentPrice, &a.Type, &a.Status, &a.StartTime, &a.EndTime, &a.ImagePath, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		auctions = append(auctions, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &auctions, nil
}

func (a *AuctionStore) GetAuctionByWinnerId(ctx context.Context, winnerID string) (*models.Auction, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	auction := &models.Auction{}

	query := `SELECT id, seller_id, winner_id, title, description, starting_price, current_price, type, status, start_time, end_time, image_path, created_at FROM auctions WHERE winner_id = $1`

	if err := a.db.QueryRowContext(ctx, query, winnerID).Scan(&auction.ID, &auction.SellerID, &auction.WinnerID, &auction.Title, &auction.Description, &auction.StartingPrice, &auction.CurrentPrice, &auction.Type, &auction.Status, &auction.StartTime, &auction.EndTime, &auction.ImagePath, &auction.CreatedAt); err != nil {
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

func (a *AuctionStore) GetAuctions(ctx context.Context, limit, offset int, filter *models.AuctionFilter) (*[]models.Auction, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	var auctions []models.Auction

	query := `SELECT id, seller_id, title, description, starting_price, current_price, type, status, start_time, end_time, image_path, category, is_paid, created_at FROM auctions WHERE 1=1`

	args := []any{}

	if filter.Type != "" {
		query += ` AND type = $` + strconv.Itoa(len(args)+1)
		args = append(args, filter.Type)
	}
	if filter.Status != "" {
		query += ` AND status = $` + strconv.Itoa(len(args)+1)
		args = append(args, filter.Status)
	}

	if filter.Category != "" {
		query += ` AND category = $` + strconv.Itoa(len(args)+1)
		args = append(args, filter.Category)
	}

	if filter.StartingPrice != 0 {
		query += ` AND starting_price = $` + strconv.Itoa(len(args)+1)
		args = append(args, filter.StartingPrice)
	}

	query += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(len(args)+1) + ` OFFSET $` + strconv.Itoa(len(args)+2)
	args = append(args, limit, offset)

	rows, err := a.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var a models.Auction

		if err := rows.Scan(&a.ID, &a.SellerID, &a.Title, &a.Description, &a.StartingPrice, &a.CurrentPrice, &a.Type, &a.Status, &a.StartTime, &a.EndTime, &a.ImagePath, &a.Category, &a.IsPaid, &a.CreatedAt); err != nil {
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

	query := `INSERT INTO auctions (seller_id, winner_id, title, description, starting_price, current_price, type, status, start_time, end_time, image_path, category, is_paid) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id, seller_id, winner_id, title, description, starting_price, current_price, type, status, start_time, end_time, image_path, category, is_paid, created_at`

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	if err = tx.QueryRowContext(ctx, query, auction.SellerID, auction.WinnerID, auction.Title, auction.Description, auction.StartingPrice, auction.CurrentPrice, auction.Type, auction.Status, auction.StartTime, auction.EndTime, auction.ImagePath, auction.Category, auction.IsPaid).Scan(&auction.ID, &auction.SellerID, &auction.WinnerID, &auction.Title, &auction.Description, &auction.StartingPrice, &auction.CurrentPrice, &auction.Type, &auction.Status, &auction.StartTime, &auction.EndTime, &auction.ImagePath, &auction.Category, &auction.IsPaid, &auction.CreatedAt); err != nil {
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

func (a *AuctionStore) GetWonAuctionsByWinnerID(ctx context.Context, winnerID string) (*[]models.Auction, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	query := `SELECT id, seller_id, winner_id, title, description, starting_price, current_price, type, status, start_time, end_time, image_path, is_paid, created_at FROM auctions WHERE winner_id = $1 AND status = 'closed';`

	rows, err := a.db.QueryContext(ctx, query, winnerID)
	if err != nil {
		log.Printf("Failed to get won auctions: %v", err)
		return nil, err
	}

	defer rows.Close()

	var auctions []models.Auction

	for rows.Next() {
		var a models.Auction
		err := rows.Scan(&a.ID, &a.SellerID, &a.WinnerID, &a.Title, &a.Description, &a.StartingPrice, &a.CurrentPrice, &a.Type, &a.Status, &a.StartTime, &a.EndTime, &a.ImagePath, &a.IsPaid, &a.CreatedAt)
		if err != nil {
			return nil, err
		}

		auctions = append(auctions, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err

	}

	return &auctions, nil
}

func (a *AuctionStore) UpdateAuctionPaymentStatus(ctx context.Context, isPaid bool, id string) error {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	query := `UPDATE auctions SET is_paid = $1 WHERE id = $2`

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if _, err = tx.ExecContext(ctx, query, isPaid, id); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (a *AuctionStore) GetBiddedAuctions(ctx context.Context, bidderID string) (*[]models.Auction, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	query := `SELECT DISTINCT 
  a.id,
  a.seller_id,
  a.winner_id,
  a.title,
  a.description,
  a.starting_price,
  a.current_price,
  a.type,
  a.status,
  a.start_time,
  a.end_time,
  a.image_path,
  a.category,
  a.is_paid,
  a.created_at
FROM auctions a
JOIN bid b ON a.id = b.auction_id
WHERE b.bidder_id = $1;`

	rows, err := a.db.QueryContext(ctx, query, bidderID)
	if err != nil {
		log.Printf("SQL query error: %v", err)
		return nil, err
	}

	defer rows.Close()

	var auctions []models.Auction

	for rows.Next() {
		var a models.Auction
		err := rows.Scan(&a.ID, &a.SellerID, &a.WinnerID, &a.Title, &a.Description, &a.StartingPrice, &a.CurrentPrice, &a.Type, &a.Status, &a.StartTime, &a.EndTime, &a.ImagePath, &a.Category, &a.IsPaid, &a.CreatedAt)
		if err != nil {
			log.Printf("SQL query error: %v", err)
			return nil, err
		}

		auctions = append(auctions, a)
	}

	if err := rows.Err(); err != nil {
		return nil, err

	}

	return &auctions, nil

}
