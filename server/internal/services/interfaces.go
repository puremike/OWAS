package services

import (
	"context"

	"github.com/puremike/online_auction_api/internal/models"
	"github.com/stripe/stripe-go/v82"
)

type UserServiceInterface interface {
	CreateUser(ctx context.Context, user *models.User) (*models.UserResponse, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error)
	UserProfile(ctx context.Context, username string) (*models.UserResponse, error)
	MeProfile(ctx context.Context, userID string) (*models.UserResponse, error)
	Refresh(ctx context.Context, refreshToken string) (string, error)
	UpdateProfile(ctx context.Context, req *models.User, id string) (string, error)
	ChangePassword(ctx context.Context, req *models.PasswordUpdateRequest, id string) (string, error)
	GetUsers(ctx context.Context) (*[]models.UserResponse, error)
	DeleteUser(ctx context.Context, id string) (string, error)
}

type AuctionServiceInterface interface {
	CreateAuction(ctx context.Context, req *models.Auction) (*models.CreateAuctionResponse, error)
	UpdateAuction(ctx context.Context, req *models.Auction, id string) (string, error)
	DeleteAuction(ctx context.Context, id string) (string, error)
	GetAuctionById(ctx context.Context, id string) (*models.CreateAuctionResponse, error)
	GetAuctions(ctx context.Context, limit, offset int, filter *models.AuctionFilter) (*[]models.CreateAuctionResponse, error)
	GetAuctionsBySellerID(ctx context.Context, sellerID string) (*[]models.CreateAuctionResponse, error)
	GetWonAuctionsByWinnerID(ctx context.Context, winnerID string) (*[]models.CreateAuctionResponse, error)
	GetBiddedAuctionsForUser(ctx context.Context, bidderID string) (*[]models.Auction, error)
	PlaceBid(ctx context.Context, req *models.PlaceBidRequest) (*models.BidResponse, error)
	CloseAuction(ctx context.Context, auctionID string, requestingUserID string) (*models.WinnerResponse, error)
}

type CSServiceInterface interface {
	ContactSupport(ctx context.Context, req *models.ContactSupport) (*models.SupportRes, error)
}

type PaymentServiceInterface interface {
	CreatePaymentCheckout(ctx context.Context, amount int64, orderID, buyerID, auctionID string) (*stripe.CheckoutSession, error)
	HandleCheckoutSessionCompleted(ctx context.Context, event *stripe.Event, session *stripe.CheckoutSession) error
	HandlePaymentIntentSucceeded(ctx context.Context, event *stripe.Event, pi *stripe.PaymentIntent) error
	HandlePaymentIntentFailed(ctx context.Context, event *stripe.Event, pi *stripe.PaymentIntent) error
	GetPayment(ctx context.Context, orderID, buyerID string) (*models.Payment, error)
	UpdateAuctionPayment(ctx context.Context, isPaid bool, id string) error
}
