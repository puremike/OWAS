package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/puremike/online_auction_api/internal/errs"
	"github.com/puremike/online_auction_api/internal/models"
	"github.com/puremike/online_auction_api/internal/payments"
	"github.com/puremike/online_auction_api/internal/store"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
)

type PaymentService struct {
	stripe *payments.StripePayment
	repo   store.PaymentRepository
}

func NewPaymentService(stripe *payments.StripePayment, repo store.PaymentRepository) *PaymentService {
	return &PaymentService{
		stripe: stripe,
		repo:   repo,
	}
}

const (
	PaymentStatusPending   = "pending"
	PaymentStatusCompleted = "completed"
	PaymentStatusFailed    = "failed"
)

func (p *PaymentService) CreatePaymentCheckout(ctx context.Context, amount int64, orderID, buyerID, auctionID string) (*stripe.CheckoutSession, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryDefaultContext)
	defer cancel()

	if amount < 0 {
		log.Printf("amount cannot be negative: %v", amount)
		return nil, errs.ErrAmountCannotBeNegative
	}

	amountInSmallestUnit := amount * 100 // convert to cent
	params := &stripe.CheckoutSessionParams{
		LineItems: []*stripe.CheckoutSessionLineItemParams{{
			PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
				Currency: stripe.String(string(stripe.CurrencyUSD)),
				ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
					Name: stripe.String("Order Payment"),
				},
				UnitAmount: stripe.Int64(amountInSmallestUnit),
			},
			Quantity: stripe.Int64(1),
		}},

		Mode:       stripe.String(stripe.CheckoutSessionModePayment),
		SuccessURL: stripe.String(p.stripe.SuccessURL),
		CancelURL:  stripe.String(p.stripe.CancelURL),

		PaymentMethodTypes: stripe.StringSlice([]string{
			string(stripe.PaymentMethodTypeCard),
		}),
	}

	params.AddMetadata("order_id", orderID)
	params.AddMetadata("buyer_id", buyerID)

	// --- ADD THESE LOGS ---
	log.Printf("DEBUG: CreatePaymentCheckout - About to add metadata. Incoming orderID: '%s', Incoming buyerID: '%s'", orderID, buyerID)

	params.AddMetadata("order_id", orderID)
	params.AddMetadata("buyer_id", buyerID)

	// --- ADD THIS LOG to inspect params.Metadata *before* API call ---
	if params.Metadata != nil {
		log.Printf("DEBUG: CreatePaymentCheckout - Params metadata before API call: %+v", params.Metadata)
	} else {
		log.Println("DEBUG: CreatePaymentCheckout - Params metadata is nil before API call.")
	}

	session, err := session.New(params)
	if err != nil {
		log.Printf("failed to create Stripe checkout session: %v", err)
		return nil, errs.ErrFailedToCreateStripeCheckout
	}

	// --- ADD THIS LOG to inspect session.Metadata *after* API call ---
	log.Printf("DEBUG: Stripe Session Created - SessionID: %s, Metadata from Stripe response: %+v", session.ID, session.Metadata)

	req := &models.Payment{
		Amount:    float64(amount),
		OrderID:   orderID,
		BuyerID:   buyerID,
		Status:    PaymentStatusPending,
		AuctionID: auctionID,
		SessionID: session.ID,
	}

	// create payment and save to DB
	if err := p.repo.CreatePayment(ctx, req); err != nil {
		log.Printf("failed to create payment: %v", err)
		return nil, errs.ErrFailedToCreatePayment
	}

	return session, nil
}

// func (p *PaymentService) GetPaymentStatus(sessionID string) (*stripe.CheckoutSession, error) {
// 	session, err := session.Get(sessionID, nil)
// 	if err != nil {
// 		return nil, errs.ErrFailedToGetPaymentSession
// 	}

// 	return session, nil
// }

// func (p *PaymentService) GetPayment(ctx context.Context, orderID, buyerID string) (*models.Payment, error) {
// 	return p.repo.GetPayment(ctx, orderID, buyerID)
// }

func (p *PaymentService) HandleCheckoutSessionCompleted(ctx context.Context, event *stripe.Event, session *stripe.CheckoutSession) error {

	ctx, cancel := context.WithTimeout(ctx, QueryDefaultContext)
	defer cancel()

	err := json.Unmarshal(event.Data.Raw, &session)
	if err != nil {
		log.Printf("failed to unmarshal event data: %v", err)
		return errs.ErrFailedToUnmarshalEvent
	}

	orderID := session.Metadata["order_id"]
	buyerID := session.Metadata["buyer_id"]
	stripeSessionID := session.ID

	if orderID == "" || buyerID == "" {
		log.Printf("Missing order_id or user_id in session metadata for session %s", stripeSessionID)
		return errs.ErrMissingRequiredSessionMetadata
	}

	log.Printf("handling checkout session completed for session: %s, order: %s and user: %s", stripeSessionID, orderID, buyerID)

	payment, err := p.repo.GetPayment(ctx, orderID)
	if err != nil {
		log.Printf("failed to get payment for order %s and user %s:", orderID, buyerID)
		return errs.ErrFailedToGetPayment
	}

	if payment.Status == PaymentStatusCompleted || payment.Status == PaymentStatusFailed {
		log.Printf("payment %s (Order: %s) already in terminal status '%s', skipping checkout.session.completed update.", payment.ID, orderID, payment.Status)
		return nil
	}

	newStatus := ""

	switch session.PaymentStatus {
	case stripe.CheckoutSessionPaymentStatusPaid:
		newStatus = PaymentStatusCompleted
		log.Printf("Checkout session %s is paid. Setting internal payment status t %s", stripeSessionID, newStatus)
	case stripe.CheckoutSessionPaymentStatusUnpaid:
		newStatus = PaymentStatusPending
		log.Printf("Checkout session %s is unpaid. Setting internal payment status t %s", stripeSessionID, newStatus)
	default:
		log.Printf("Checkout session %s has unknown payment status '%s'. Setting internal payment status t %s", stripeSessionID, session.PaymentStatus, newStatus)
	}

	if newStatus != "" && payment.Status != newStatus {
		if err := p.repo.UpdatePayment(ctx, newStatus, payment.ID); err != nil {
			log.Printf("failed to update payment: %v", err)
			return errs.ErrFailedToUpdatePayment
		}
	}

	return nil
}

func (p *PaymentService) HandlePaymentIntentSucceeded(ctx context.Context, event *stripe.Event, pi *stripe.PaymentIntent) error {

	ctx, cancel := context.WithTimeout(ctx, QueryDefaultContext)
	defer cancel()

	err := json.Unmarshal(event.Data.Raw, &pi)
	if err != nil {
		log.Printf("failed to unmarshal event data: %v", err)
		return errs.ErrFailedToUnmarshalEvent
	}

	// buyerID := pi.Metadata["buyer_id"]
	orderID := pi.Metadata["order_id"]

	payment, err := p.repo.GetPayment(ctx, orderID)
	if err != nil {
		log.Printf("failed to get payment: %v", err)
		return errs.ErrFailedToGetPayment
	}

	if payment.Status == PaymentStatusCompleted {
		log.Printf("payment %s (Order: %s) already in terminal status '%s', skipping payment_intent.succeeded update.", payment.ID, orderID, payment.Status)
		return nil
	}

	if payment.Status == PaymentStatusFailed {
		log.Printf("Payment %s (Order: %s) was previously '%s', but received payment_intent.succeeded. Transitioning to 'completed'. PI: %s",
			payment.ID, orderID, PaymentStatusFailed, pi.ID)
	}

	newStatus := PaymentStatusCompleted
	if err := p.repo.UpdatePayment(ctx, newStatus, payment.ID); err != nil {
		log.Printf("failed to update payment: %v", err)
		return errs.ErrFailedToUpdatePayment
	}

	return nil
}

func (p *PaymentService) HandlePaymentIntentFailed(ctx context.Context, event *stripe.Event, pi *stripe.PaymentIntent) error {

	ctx, cancel := context.WithTimeout(ctx, QueryDefaultContext)
	defer cancel()

	err := json.Unmarshal(event.Data.Raw, &pi)
	if err != nil {
		log.Printf("failed to unmarshal event data: %v", err)
		return errs.ErrFailedToUnmarshalEvent
	}

	//buyerID := pi.Metadata["buyer_id"]
	orderID := pi.Metadata["order_id"]

	payment, err := p.repo.GetPayment(ctx, orderID)
	if err != nil {
		log.Printf("failed to get payment: %v", err)
		return errs.ErrFailedToGetPayment
	}

	if payment.Status == PaymentStatusCompleted || payment.Status == PaymentStatusFailed {
		log.Printf("payment %s (Order: %s) already in terminal status '%s', skipping update from payment_intent.payment_failed.", payment.ID, orderID, payment.Status)
		return nil
	}

	newStatus := PaymentStatusFailed
	if err := p.repo.UpdatePayment(ctx, newStatus, payment.ID); err != nil {
		log.Printf("failed to update payment: %v", err)
		return errs.ErrFailedToUpdatePayment
	}

	return nil
}
