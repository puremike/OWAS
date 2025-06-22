package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/puremike/online_auction_api/internal/errs"
	"github.com/puremike/online_auction_api/internal/models"
	"github.com/puremike/online_auction_api/internal/store"
)

// PlaceBid is a method in your AuctionService
func (a *AuctionService) PlaceBid(ctx context.Context, req *models.PlaceBidRequest) (*models.BidResponse, error) {
	// 1. Retrieve the auction from the database
	auction, err := a.repo.GetAuctionById(ctx, req.AuctionID)
	if err != nil {
		if errors.Is(err, errs.ErrAuctionNotFound) {
			return nil, errs.ErrAuctionNotFound
		}
		return nil, errors.New("failed to retrieve auction for bidding")
	}

	// Basic Validation:
	if auction.Status != "open" {
		return nil, errs.ErrAuctionNotOpenForBids
	}
	if req.BidAmount <= auction.CurrentPrice {
		return nil, errs.ErrBidTooLow
	}
	if req.BidAmount <= auction.StartingPrice && auction.CurrentPrice == auction.StartingPrice {
		return nil, errs.ErrBidTooLow
	}
	if req.BidderID == auction.SellerID {
		return nil, errs.ErrBidBySeller
	}
	// 3. Retrieve the previous highest bid
	previousBid, err := a.bidRepo.GetHighestBid(ctx, req.AuctionID)
	if err != nil {
		if errors.Is(err, errs.ErrBidNotFound) {
			previousBid = nil
		} else {
			log.Printf("GetHighestBid failed: %v", err)
			return nil, errs.ErrFailedToGetHighestBid
		}
	}

	previousHighestBidderID := ""

	if previousBid != nil {
		previousHighestBidderID = previousBid.BidderID
	}

	// 4. Create bid
	newBid := &models.Bid{
		AuctionID: req.AuctionID,
		BidderID:  req.BidderID,
		Amount:    req.BidAmount,
	}

	savedBid, err := a.bidRepo.CreateBid(ctx, newBid)
	if err != nil {
		return nil, errs.ErrFailedToSaveBid
	}

	// 5. Update auction with new current price
	auction.CurrentPrice = req.BidAmount
	if err := a.repo.UpdateAuction(ctx, auction, req.AuctionID); err != nil {
		return nil, errs.ErrFailedToUpdateAuction
	}

	// 6. Notify WebSocket listeners
	a.auctionUpdates <- &models.AuctionUpdateEvent{
		EventType:    models.AuctionNewBid,
		ID:           req.AuctionID,
		CurrentPrice: req.BidAmount,
		SellerID:     req.BidderID,
		TimeStamp:    time.Now(),
	}

	not := &store.Notification{
		UserID:  req.BidderID,
		Message: fmt.Sprintf("You have placed a bid on auction: %s, auctionId: %s", auction.Title, req.AuctionID),
		IsRead:  false,
	}
	if err := a.notRepo.CreateNotification(ctx, not); err != nil {
		return nil, fmt.Errorf("CreateNotification failed: %v", err)
	}

	// 7. Notify previous highest bidder (if different)
	if previousHighestBidderID != "" && previousHighestBidderID != req.BidderID {
		a.notifications <- &models.NotificationEvent{
			Type:      models.NotificationOutBid,
			UserID:    previousHighestBidderID,
			Message:   fmt.Sprintf("You have been outbid on auction: %s", auction.Title),
			AuctionID: req.AuctionID,
			TimeStamp: time.Now(),
		}

		not := &store.Notification{
			UserID:  previousHighestBidderID,
			Message: fmt.Sprintf("You have been outbid on auction: %s, auctionId: %s", auction.Title, req.AuctionID),
			IsRead:  false,
		}
		if err := a.notRepo.CreateNotification(ctx, not); err != nil {
			return nil, fmt.Errorf("CreateNotification failed: %v", err)
		}
	}

	// 8. Return response
	return &models.BidResponse{
		AuctionID: req.AuctionID,
		BidderID:  req.BidderID,
		BidAmount: req.BidAmount,
		TimeStamp: savedBid.CreatedAt,
	}, nil
}

// CloseAuction method in your AuctionService
// Now accepts the authenticated userID for authorization checks.
func (a *AuctionService) CloseAuction(ctx context.Context, auctionID string, requestingUserID string) error {
	auction, err := a.repo.GetAuctionById(ctx, auctionID)
	if err != nil {
		if errors.Is(err, errs.ErrAuctionNotFound) {
			return errs.ErrAuctionNotFound
		}
		return errors.New("failed to retrieve auction for closing")
	}

	// Authorization
	if auction.SellerID != requestingUserID {
		return errs.ErrPermissionDenied
	}

	// Prevent re-closing
	if auction.Status == "closed" {
		return errs.ErrAuctionAlreadyClosed
	}

	// Update status
	auction.Status = "closed"
	if err := a.repo.CloseAuction(ctx, auction.Status, auctionID); err != nil {
		return errs.ErrFailedToUpdateAuction
	}

	// Determine the winner: highest bid
	var winnerID string
	if auction.CurrentPrice > auction.StartingPrice {
		highestBid, err := a.bidRepo.GetHighestBid(ctx, auctionID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return errors.New("failed to retrieve highest bid during auction close")
		}
		if highestBid != nil {
			winnerID = highestBid.BidderID
		}
	}

	// Notify winner
	if winnerID != "" {
		a.notifications <- &models.NotificationEvent{
			Type:      models.NotificationWon,
			UserID:    winnerID,
			Message:   fmt.Sprintf("Congratulations! You won the auction: %s", auction.Title),
			AuctionID: auctionID,
			TimeStamp: time.Now(),
		}

		not := &store.Notification{
			UserID:  winnerID,
			Message: fmt.Sprintf("Congratulations! You won the auction: %s, auctionId: %s", auction.Title, auctionID),
			IsRead:  false,
		}
		if err := a.notRepo.CreateNotification(ctx, not); err != nil {
			return fmt.Errorf("CreateNotification failed: %v", err)
		}
	}

	// Notify other bidders
	bidders, err := a.bidRepo.GetAllBidderIDsForAuction(ctx, auctionID)
	if err != nil {
		return errors.New("failed to retrieve all bidders for auction close")
	}

	uniqueBidders := make(map[string]struct{})
	for _, id := range bidders {
		if id != winnerID && id != requestingUserID {
			if _, seen := uniqueBidders[id]; !seen {
				uniqueBidders[id] = struct{}{}
				a.notifications <- &models.NotificationEvent{
					Type:      models.NotificationWon,
					UserID:    id,
					Message:   fmt.Sprintf("Auction %s has ended. You did not win.", auction.Title),
					AuctionID: auctionID,
					TimeStamp: time.Now(),
				}
			}
		}

		not := &store.Notification{
			UserID:  id,
			Message: fmt.Sprintf("Auction %s has ended. You did not win.", auction.Title),
			IsRead:  false,
		}
		if err := a.notRepo.CreateNotification(ctx, not); err != nil {
			return fmt.Errorf("CreateNotification failed: %v", err)
		}
	}

	// Notify WebSocket listeners
	a.auctionUpdates <- &models.AuctionUpdateEvent{
		EventType: models.AuctionEnded,
		ID:        auctionID,
		Status:    "closed",
		TimeStamp: time.Now(),
	}

	return nil
}
