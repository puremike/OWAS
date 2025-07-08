package store

// import (
// 	"context"
// 	"time"

// 	"github.com/puremike/online_auction_api/internal/models"
// )

// type MockUserRepository interface {
// 	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
// 	GetUserById(ctx context.Context, id string) (*models.User, error)
// 	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
// 	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
// 	StoreRefreshToken(ctx context.Context, userID, refreshToken string, expires_at time.Time) error
// 	UpdateUser(ctx context.Context, user *models.User, id string) error
// 	ValidateRefreshToken(ctx context.Context, refreshToken string) (string, error)
// 	ChangePassword(ctx context.Context, pass, id string) error
// 	GetUsers(ctx context.Context) (*[]models.User, error)
// 	DeleteUser(ctx context.Context, id string) error
// }

// type MockAuctionRepository interface {
// 	GetAuctionById(ctx context.Context, id string) (*models.Auction, error)
// 	GetAuctions(ctx context.Context, limit, offset int, filter *models.AuctionFilter) (*[]models.Auction, error)
// 	CreateAuction(ctx context.Context, auction *models.Auction) (*models.Auction, error)
// 	CloseAuction(ctx context.Context, status, id string) error
// 	UpdateAuction(ctx context.Context, auction *models.Auction, id string) error
// 	DeleteAuction(ctx context.Context, id string) error
// 	GetWonAuctionsByWinnerID(ctx context.Context, winnerID string) (*[]models.Auction, error)
// 	UpdateAuctionPaymentStatus(ctx context.Context, isPaid bool, id string) error
// 	GetBiddedAuctions(ctx context.Context, bidderID string) (*[]models.Auction, error)
// 	GetAuctionByWinnerId(ctx context.Context, winnerID string) (*models.Auction, error)
// 	GetAuctionBySellerId(ctx context.Context, sellerID string) (*[]models.Auction, error)
// }

// type MockBidRepository interface {
// 	GetHighestBid(ctx context.Context, id string) (*models.Bid, error)
// 	GetBidById(ctx context.Context, id string) (*models.Bid, error)
// 	GetBids(ctx context.Context, userId string) (*[]models.Bid, error)
// 	CreateBid(ctx context.Context, bid *models.Bid) (*models.Bid, error)
// 	GetAllBidderIDsForAuction(ctx context.Context, auctionID string) ([]string, error)
// 	GetBidByUser(ctx context.Context, auctionID, bidderID string) (*models.Bid, error)
// 	DeleteBidsByAuction(ctx context.Context, auctionID string) error
// }

// type MockPaymentRepository interface {
// 	CreatePayment(ctx context.Context, payment *models.Payment) error
// 	GetPayment(ctx context.Context, orderID, buyerID string) (*models.Payment, error)
// 	UpdatePayment(ctx context.Context, paymentStatus, id string) error
// }

// type MockNotificationRepository interface {
// 	CreateNotification(ctx context.Context, notification *Notification) error
// 	GetNotifications(ctx context.Context, userID string) ([]*Notification, error)
// 	DeleteNotificationByAuction(ctx context.Context, auctionID string) error
// }

// type MockCSRepository interface {
// 	ContactSupport(ctx context.Context, cs *models.ContactSupport) (*models.ContactSupport, error)
// }

// type MockStorage struct {
// 	Users         MockUserRepository
// 	Auctions      MockAuctionRepository
// 	Bids          MockBidRepository
// 	Payments      MockPaymentRepository
// 	Notifications MockNotificationRepository
// 	CS            MockCSRepository
// }

// func MockNewStorage() *Storage {
// 	return &Storage{
// 		Users:         &MockUserStore{},
// 		Auctions:      &MockAuctionStore{},
// 		Bids:          &MockBidStore{},
// 		Payments:      &MockPaymentStore{},
// 		Notifications: &MockNotificationStore{},
// 		CS:            &MockCSStore{},
// 	}
// }

// type MockAuctionStore struct {
// }

// func (m *MockAuctionStore) GetAuctionById(ctx context.Context, id string) (*models.Auction, error) {
// 	return &models.Auction{ID: "id"}, nil
// }
// func (m *MockAuctionStore) GetAuctions(ctx context.Context, limit, offset int, filter *models.AuctionFilter) (*[]models.Auction, error) {
// 	return &[]models.Auction{}, nil
// }
// func (m *MockAuctionStore) CreateAuction(ctx context.Context, auction *models.Auction) (*models.Auction, error) {
// 	return &models.Auction{ID: "id"}, nil
// }
// func (m *MockAuctionStore) CloseAuction(ctx context.Context, status, id string) error {
// 	return nil
// }
// func (m *MockAuctionStore) UpdateAuction(ctx context.Context, auction *models.Auction, id string) error {
// 	return nil
// }
// func (m *MockAuctionStore) DeleteAuction(ctx context.Context, id string) error {
// 	return nil
// }
// func (m *MockAuctionStore) GetWonAuctionsByWinnerID(ctx context.Context, winnerID string) (*[]models.Auction, error) {
// 	return &[]models.Auction{}, nil
// }
// func (m *MockAuctionStore) UpdateAuctionPaymentStatus(ctx context.Context, isPaid bool, id string) error {
// 	return nil
// }
// func (m *MockAuctionStore) GetBiddedAuctions(ctx context.Context, bidderID string) (*[]models.Auction, error) {
// 	return &[]models.Auction{}, nil
// }
// func (m *MockAuctionStore) GetAuctionByWinnerId(ctx context.Context, winnerID string) (*models.Auction, error) {
// 	return &models.Auction{ID: "id"}, nil
// }
// func (m *MockAuctionStore) GetAuctionBySellerId(ctx context.Context, sellerID string) (*[]models.Auction, error) {
// 	return &[]models.Auction{}, nil
// }

// type MockBidStore struct{}

// func (m *MockBidStore) GetHighestBid(ctx context.Context, id string) (*models.Bid, error) {
// 	return &models.Bid{ID: "id"}, nil
// }
// func (m *MockBidStore) GetBidById(ctx context.Context, id string) (*models.Bid, error) {
// 	return &models.Bid{ID: "id"}, nil
// }
// func (m *MockBidStore) GetBids(ctx context.Context, userId string) (*[]models.Bid, error) {
// 	return &[]models.Bid{}, nil
// }
// func (m *MockBidStore) CreateBid(ctx context.Context, bid *models.Bid) (*models.Bid, error) {
// 	return &models.Bid{ID: "id"}, nil
// }
// func (m *MockBidStore) GetAllBidderIDsForAuction(ctx context.Context, auctionID string) ([]string, error) {
// 	return []string{"id"}, nil
// }
// func (m *MockBidStore) GetBidByUser(ctx context.Context, auctionID, bidderID string) (*models.Bid, error) {
// 	return &models.Bid{ID: "id"}, nil
// }
// func (m *MockBidStore) DeleteBidsByAuction(ctx context.Context, auctionID string) error {
// 	return nil
// }

// type MockPaymentStore struct {
// }

// func (m *MockPaymentStore) CreatePayment(ctx context.Context, payment *models.Payment) error {
// 	return nil
// }
// func (m *MockPaymentStore) GetPayment(ctx context.Context, orderID, buyerID string) (*models.Payment, error) {
// 	return &models.Payment{ID: "id"}, nil
// }
// func (m *MockPaymentStore) UpdatePayment(ctx context.Context, paymentStatus, id string) error {
// 	return nil
// }

// type MockNotificationStore struct {
// }

// func (m *MockNotificationStore) CreateNotification(ctx context.Context, notification *Notification) error {
// 	return nil
// }
// func (m *MockNotificationStore) GetNotifications(ctx context.Context, userID string) ([]*Notification, error) {
// 	return []*Notification{}, nil
// }
// func (m *MockNotificationStore) DeleteNotificationByAuction(ctx context.Context, auctionID string) error {
// 	return nil
// }

// type MockCSStore struct {
// }

// func (m *MockCSStore) ContactSupport(ctx context.Context, cs *models.ContactSupport) (*models.ContactSupport, error) {
// 	return &models.ContactSupport{ID: 12}, nil
// }
