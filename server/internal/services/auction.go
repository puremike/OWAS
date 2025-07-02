package services

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/puremike/online_auction_api/internal/errs"
	"github.com/puremike/online_auction_api/internal/models"
	"github.com/puremike/online_auction_api/internal/store"
)

type AuctionService struct {
	repo           store.AuctionRepository
	bidRepo        store.BidRepository
	notRepo        store.NotificationRepository
	auctionUpdates chan<- *models.AuctionUpdateEvent
	notifications  chan<- *models.NotificationEvent
}

func NewAuctionService(repo store.AuctionRepository, bidRepo store.BidRepository, notRepo store.NotificationRepository, auctionUpdates chan<- *models.AuctionUpdateEvent, notifications chan<- *models.NotificationEvent) *AuctionService {
	return &AuctionService{
		repo:           repo,
		bidRepo:        bidRepo,
		notRepo:        notRepo,
		auctionUpdates: auctionUpdates,
		notifications:  notifications,
	}
}

func (a *AuctionService) CreateAuction(ctx context.Context, req *models.Auction) (*models.CreateAuctionResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryDefaultContext)
	defer cancel()

	if req.Title == "" || req.Description == "" || req.StartingPrice < 1 || req.Type == "" || req.Category == "" || req.Status == "" || req.StartTime.IsZero() || req.EndTime.IsZero() || req.SellerID == "" {
		return &models.CreateAuctionResponse{}, errs.ErrInvalidAuctionDetails
	}

	auction := &models.Auction{
		Title:         req.Title,
		Description:   req.Description,
		StartingPrice: req.StartingPrice,
		CurrentPrice:  req.StartingPrice,
		Type:          strings.ToLower(req.Type),
		Status:        "open",
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		SellerID:      req.SellerID,
		WinnerID:      req.SellerID,
		ImagePath:     req.ImagePath,
		Category:      req.Category,
		IsPaid:        false,
	}

	createdAuction, err := a.repo.CreateAuction(ctx, auction)
	if err != nil {
		return &models.CreateAuctionResponse{}, errs.ErrFailedToCreateAuction
	}

	res := &models.CreateAuctionResponse{
		ID:            createdAuction.ID,
		SellerID:      createdAuction.SellerID,
		Title:         createdAuction.Title,
		Description:   createdAuction.Description,
		StartingPrice: createdAuction.StartingPrice,
		CurrentPrice:  createdAuction.CurrentPrice,
		Type:          createdAuction.Type,
		Status:        createdAuction.Status,
		StartTime:     createdAuction.StartTime,
		EndTime:       createdAuction.EndTime,
		CreatedAt:     createdAuction.CreatedAt,
		ImagePath:     createdAuction.ImagePath,
		Category:      createdAuction.Category,
	}

	return res, nil
}

func (a *AuctionService) UpdateAuction(ctx context.Context, req *models.Auction, id string) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryDefaultContext)
	defer cancel()

	if req.Title == "" || req.Description == "" || req.StartingPrice < 1 || req.Type == "" || req.Status == "" || req.StartTime.IsZero() || req.EndTime.IsZero() || req.SellerID == "" {
		return "", errs.ErrInvalidAuctionDetails
	}

	auction := &models.Auction{
		Title:         req.Title,
		Description:   req.Description,
		StartingPrice: req.StartingPrice,
		CurrentPrice:  req.StartingPrice,
		Type:          req.Type,
		Status:        "open",
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		SellerID:      req.SellerID,
		WinnerID:      req.SellerID,
	}

	if err := a.repo.UpdateAuction(ctx, auction, id); err != nil {
		return "", errs.ErrFailedToCreateAuction
	}

	return "auction updated successfully", nil
}

func (a *AuctionService) DeleteAuction(ctx context.Context, id string) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryDefaultContext)
	defer cancel()

	if err := a.bidRepo.DeleteBidsByAuction(ctx, id); err != nil {
		return "", errs.ErrFailedToDeleteBids
	}

	if err := a.repo.DeleteAuction(ctx, id); err != nil {
		return "", errs.ErrFailedToDeleteAuction
	}

	return "auction deleted successfully", nil
}

func (a *AuctionService) GetAuctionById(ctx context.Context, id string) (*models.CreateAuctionResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryDefaultContext)
	defer cancel()

	auction, err := a.repo.GetAuctionById(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrAuctionNotFound) {
			return &models.CreateAuctionResponse{}, errs.ErrAuctionNotFound
		}
		return &models.CreateAuctionResponse{}, fmt.Errorf("failed to retrieve auction: %w", err)
	}

	res := &models.CreateAuctionResponse{
		ID:            auction.ID,
		SellerID:      auction.SellerID,
		Title:         auction.Title,
		Description:   auction.Description,
		StartingPrice: auction.StartingPrice,
		CurrentPrice:  auction.CurrentPrice,
		Type:          auction.Type,
		Status:        auction.Status,
		StartTime:     auction.StartTime,
		EndTime:       auction.EndTime,
		CreatedAt:     auction.CreatedAt,
		ImagePath:     auction.ImagePath,
	}

	return res, nil
}

func (a *AuctionService) GetAuctions(ctx context.Context, limit, offset int, filter *models.AuctionFilter) (*[]models.CreateAuctionResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryDefaultContext)
	defer cancel()

	auctions, err := a.repo.GetAuctions(ctx, limit, offset, filter)
	if err != nil {
		return &[]models.CreateAuctionResponse{}, errors.New("failed to retrieve auctions")
	}

	res := &[]models.CreateAuctionResponse{}

	for _, auction := range *auctions {
		*res = append(*res, models.CreateAuctionResponse{
			ID:            auction.ID,
			SellerID:      auction.SellerID,
			Title:         auction.Title,
			Description:   auction.Description,
			StartingPrice: auction.StartingPrice,
			CurrentPrice:  auction.CurrentPrice,
			Type:          auction.Type,
			Status:        auction.Status,
			StartTime:     auction.StartTime,
			EndTime:       auction.EndTime,
			CreatedAt:     auction.CreatedAt,
			ImagePath:     auction.ImagePath,
			Category:      auction.Category,
			IsPaid:        auction.IsPaid,
		})
	}

	return res, nil
}

func (a *AuctionService) GetAuctionsBySellerID(ctx context.Context, sellerID string) (*[]models.CreateAuctionResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryDefaultContext)
	defer cancel()

	auctions, err := a.repo.GetAuctionBySellerId(ctx, sellerID)
	if err != nil {
		return &[]models.CreateAuctionResponse{}, errors.New("failed to retrieve auctions")
	}

	res := &[]models.CreateAuctionResponse{}

	for _, auction := range *auctions {
		*res = append(*res, models.CreateAuctionResponse{
			ID:            auction.ID,
			SellerID:      auction.SellerID,
			Title:         auction.Title,
			Description:   auction.Description,
			StartingPrice: auction.StartingPrice,
			CurrentPrice:  auction.CurrentPrice,
			Type:          auction.Type,
			Status:        auction.Status,
			StartTime:     auction.StartTime,
			EndTime:       auction.EndTime,
			CreatedAt:     auction.CreatedAt,
			ImagePath:     auction.ImagePath,
			IsPaid:        auction.IsPaid,
		})
	}

	return res, nil
}

func (a *AuctionService) GetWonAuctionsByWinnerID(ctx context.Context, winnerID string) (*[]models.CreateAuctionResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryDefaultContext)
	defer cancel()

	existingAuction, err := a.repo.GetAuctionByWinnerId(ctx, winnerID)
	if err != nil {
		return &[]models.CreateAuctionResponse{}, errs.ErrAuctionNotFound
	}

	if existingAuction.Status != "closed" {
		return &[]models.CreateAuctionResponse{}, errs.ErrAuctionNotFound
	}

	auctions, err := a.repo.GetWonAuctionsByWinnerID(ctx, winnerID)
	if err != nil {
		return &[]models.CreateAuctionResponse{}, errors.New("failed to retrieve auctions")
	}

	res := &[]models.CreateAuctionResponse{}

	for _, auction := range *auctions {
		*res = append(*res, models.CreateAuctionResponse{
			ID:            auction.ID,
			SellerID:      auction.SellerID,
			Title:         auction.Title,
			Description:   auction.Description,
			StartingPrice: auction.StartingPrice,
			CurrentPrice:  auction.CurrentPrice,
			Type:          auction.Type,
			Status:        auction.Status,
			StartTime:     auction.StartTime,
			EndTime:       auction.EndTime,
			CreatedAt:     auction.CreatedAt,
			ImagePath:     auction.ImagePath,
			IsPaid:        auction.IsPaid,
		})
	}

	return res, nil
}

func (a *AuctionService) GetBiddedAuctionsForUser(ctx context.Context, bidderID string) (*[]models.Auction, error) {
	return a.repo.GetBiddedAuctions(context.Background(), bidderID)
}
