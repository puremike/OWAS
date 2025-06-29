package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/puremike/online_auction_api/contexts"
	"github.com/puremike/online_auction_api/internal/errs"
	"github.com/puremike/online_auction_api/internal/models"
	"github.com/puremike/online_auction_api/internal/services"
	"github.com/puremike/online_auction_api/internal/store"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

type WebHookHandler struct {
	service     *services.PaymentService
	auctionRepo store.AuctionRepository
}

func NewWebHookHander(service *services.PaymentService, auctionRepo store.AuctionRepository) *WebHookHandler {
	return &WebHookHandler{
		service:     service,
		auctionRepo: auctionRepo,
	}
}

// StripeWebHookHandler handles Stripe webhook events.
//
// This handler reads and verifies incoming Stripe webhook requests, ensuring
// the request body is correctly read and the Stripe signature is valid.
// It processes different event types such as "checkout.session.completed",
// "payment_intent.succeeded", and "payment_intent.payment_failed". Each event
// type is handled by invoking appropriate service methods to update payment
// statuses or handle session completions. If the event type is unrecognized,
// it returns an error response indicating the event type is unhandled.
//
//	@Summary		Handle Stripe Webhook Events
//	@Description	Processes Stripe webhook events for payment and checkout session updates.
//	@Tags			Webhook
//	@Accept			json
//	@Produce		json
//	@Param			Stripe-Signature	header		string	true	"Stripe Signature Header"
//	@Success		200					{object}	gin.H	"success"
//	@Failure		400					{object}	gin.H	"Bad Request - invalid input or unhandled event type"
//	@Failure		500					{object}	gin.H	"Internal Server Error - failed to process event"
//	@Router			/webhook/stripe [post]
func (w *WebHookHandler) StripeWebHookHandler(c *gin.Context) {

	// read the body request
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("error reading request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "error reading request body"})
		return
	}

	// verify the signature header
	signatureHeader := c.GetHeader("Stripe-Signature")
	if signatureHeader == "" {
		log.Println("missing Stripe-Signature header")
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing Stripe-Signature header"})
		return
	}

	// verify the webhook secret
	stripeWebhookSecret, exists := os.LookupEnv("STRIPE_WEBHOOK_SECRET")
	if !exists {
		log.Println("missing STRIPE_WEBHOOK_SECRET environment variable")
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing STRIPE_WEBHOOK_SECRET environment variable"})
		return
	}

	event, err := webhook.ConstructEvent(body, signatureHeader, stripeWebhookSecret)
	if err != nil {
		log.Printf("error verifying webhook signature: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "error verifying webhook signature"})
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession

		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			log.Printf("failed to unmarshal event data: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to unmarshal event data"})
			return
		}

		if err := w.service.HandleCheckoutSessionCompleted(c, &event, &session); err != nil {
			log.Printf("failed to handle checkout.session.completed event: %v", err)
			errs.MapServiceErrors(c, err)
			return
		}

	case "payment_intent.succeeded":
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			log.Printf("failed to unmarshal event data: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to unmarshal event data"})
			return
		}

		if err := w.service.HandlePaymentIntentSucceeded(c, &event, &pi); err != nil {
			log.Printf("failed to handle payment_intent.succeeded event: %v", err)
			errs.MapServiceErrors(c, err)
			return
		}

	case "payment_intent.payment_failed":
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			log.Printf("failed to unmarshal event data: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "failed to unmarshal event data"})
			return
		}

		if err := w.service.HandlePaymentIntentFailed(c, &event, &pi); err != nil {
			log.Printf("failed to handle payment_intent.payment_failed event: %v", err)
			errs.MapServiceErrors(c, err)
			return
		}

	default:
		log.Printf("unhandled event type: %s", event.Type)
		c.JSON(http.StatusBadRequest, gin.H{"error": "unhandled event type"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

// CreateCheckoutSessionHandler godoc
//
//	@Summary		Create Stripe Checkout Session for an auction
//	@Description	Create a Stripe Checkout Session for an auction, using the current price of the auction and the authenticated user's ID.
//	@Tags			Payments
//	@Accept			json
//	@Produce		json
//	@Param			auction_id	path		string							true	"ID of the auction to create a checkout session for"
//	@Success		201			{object}	models.CreatePaymentResponse	"Stripe Checkout Session created successfully"
//	@Failure		400			{object}	gin.H							"Bad Request - invalid input"
//	@Failure		401			{object}	gin.H							"Unauthorized - user not authenticated"
//	@Failure		404			{object}	gin.H							"Not Found - auction not found"
//	@Failure		500			{object}	gin.H							"Internal Server Error - failed to create Stripe Checkout Session"
//	@Router			/auctions/{auctionID}/create-checkout-session [post]
//
//	@Security		jwtCookieAuth
func (w *WebHookHandler) CreateCheckoutSessionHandler(c *gin.Context) {

	authUser, err := contexts.GetUserFromContext(c)
	if authUser == nil || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	auction, err := contexts.GetAuctionFromContext(c)
	if err != nil {
		if errors.Is(err, errs.ErrAuctionNotFound) {
			errs.MapServiceErrors(c, err)
			return
		}
		return
	}

	if authUser.ID != auction.WinnerID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you're not allowed to proceed"})
		return
	}

	orderID := uuid.New().String()

	// Call the service layer to create the Stripe Checkout Session
	stripeSession, err := w.service.CreatePaymentCheckout(c.Request.Context(), int64(auction.CurrentPrice), orderID, authUser.ID, auction.ID)
	if err != nil {
		log.Printf("failed to create payment intent in service: %v", err)
		errs.MapServiceErrors(c, err)
		return
	}

	c.JSON(http.StatusOK, models.CreatePaymentResponse{
		CheckoutURL: stripeSession.URL,
	})

	log.Printf("successfully created Stripe Checkout Session for Order %s, Buyer %s. URL: %s", orderID, authUser.ID, stripeSession.URL)
}
