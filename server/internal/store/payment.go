package store

import (
	"context"
	"database/sql"

	"github.com/puremike/online_auction_api/internal/models"
)

type PaymentStore struct {
	db *sql.DB
}

func (p *PaymentStore) CreatePayment(ctx context.Context, payment *models.Payment) error {

	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	query := `INSERT INTO payment (auction_id, buyer_id, order_id, amount, status) VALUES ($1, $2, $3, $4, $5)`

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil
	}

	defer tx.Rollback()

	if err = tx.QueryRowContext(ctx, query, payment.AuctionID, payment.BuyerID, payment.OrderID, payment.Amount, payment.Status).Scan(&payment.ID); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (p *PaymentStore) GetPayment(ctx context.Context, orderID, buyerID string) (*models.Payment, error) {
	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	var payment models.Payment

	query := `SELECT id, auction_id, buyer_id, order_id, amount, status FROM payment WHERE order_id = $1 AND buyer_id = $2`

	if err := p.db.QueryRowContext(ctx, query, orderID, buyerID).Scan(&payment.ID, &payment.AuctionID, &payment.BuyerID, &payment.OrderID, &payment.Amount, &payment.Status, &payment.CreatedAt); err != nil {
		return nil, err
	}

	return &payment, nil
}

func (p *PaymentStore) UpdatePayment(ctx context.Context, paymentStatus, id string) error {

	ctx, cancel := context.WithTimeout(ctx, QueryBackgroundTimeout)
	defer cancel()

	query := `UPDATE payment SET status = $1 WHERE id = $2`

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, query, paymentStatus, id); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
