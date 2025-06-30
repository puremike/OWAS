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
	auction, err := a.repo.GetAuctionById(ctx, req.AuctionID)
	if err != nil {
		if errors.Is(err, errs.ErrAuctionNotFound) {
			return nil, errs.ErrAuctionNotFound
		}
		return nil, errors.New("failed to retrieve auction for bidding")
	}

	if auction.Status != "open" {
		return nil, errs.ErrAuctionNotOpenForBids
	}
	if req.BidderID == auction.SellerID {
		return nil, errs.ErrBidBySeller
	}

	switch auction.Type {
	case models.EnglishAuction:
		if req.BidAmount <= auction.CurrentPrice {
			return nil, errs.ErrBidTooLow
		}

	case models.DutchAuction:
		// Only allow ONE bid, exactly at the current price
		if auction.CurrentPrice != req.BidAmount {
			return nil, errs.ErrDutchBidMustMatchCurrent
		}

		// Optional: prevent duplicate bids if auction is already won
		existingBid, err := a.bidRepo.GetHighestBid(ctx, req.AuctionID)
		if err == nil && existingBid != nil {
			return nil, errs.ErrDutchAuctionAlreadyWon
		}

	// case "sealed":
	// 	// ✅ Allow any amount >= starting price
	// 	if req.BidAmount < auction.StartingPrice {
	// 		return nil, errs.ErrBidTooLow
	// 	}
	// 	// Prevent duplicate bids from same user
	// 	existing, _ := a.bidRepo.GetBidByUser(ctx, req.AuctionID, req.BidderID)
	// 	if existing != nil {
	// 		return nil, errs.ErrDuplicateSealedBid
	// 	}

	default:
		return nil, errors.New("unknown auction type")
	}

	// Retrieve the previous highest bid
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

	// Same for both types:
	newBid := &models.Bid{
		AuctionID: req.AuctionID,
		BidderID:  req.BidderID,
		Amount:    req.BidAmount,
	}

	savedBid, err := a.bidRepo.CreateBid(ctx, newBid)
	if err != nil {
		return nil, errs.ErrFailedToSaveBid
	}

	auction.CurrentPrice = req.BidAmount
	auction.WinnerID = req.BidderID

	// Close the auction immediately for Dutch
	if auction.Type == models.DutchAuction {
		auction.Status = "closed"
	}

	if err := a.repo.UpdateAuction(ctx, auction, req.AuctionID); err != nil {
		return nil, errs.ErrFailedToUpdateAuction
	}

	// WebSocket broadcast
	a.auctionUpdates <- &models.AuctionUpdateEvent{
		EventType:    models.AuctionNewBid,
		ID:           req.AuctionID,
		CurrentPrice: req.BidAmount,
		SellerID:     req.BidderID,
		Status:       auction.Status,
		Type:         auction.Type,
		TimeStamp:    time.Now(),
	}

	// Notification to bidder
	not := &store.Notification{
		UserID:    req.BidderID,
		Message:   fmt.Sprintf("You have placed a bid on auction: %s, auctionId: %s", auction.Title, req.AuctionID),
		AuctionID: req.AuctionID,
		IsRead:    false,
	}
	if err := a.notRepo.CreateNotification(ctx, not); err != nil {
		return nil, fmt.Errorf("CreateNotification failed: %v", err)
	}

	// Only notify previous bidder if English type
	if auction.Type == models.EnglishAuction && previousHighestBidderID != "" && previousHighestBidderID != req.BidderID {
		a.notifications <- &models.NotificationEvent{
			Type:      models.NotificationOutBid,
			UserID:    previousHighestBidderID,
			Message:   fmt.Sprintf("You have been outbid on auction: %s", auction.Title),
			AuctionID: req.AuctionID,
			TimeStamp: time.Now(),
		}
		not := &store.Notification{
			UserID:    previousBid.BidderID,
			Message:   fmt.Sprintf("You have been outbid on auction: %s, auctionId: %s", auction.Title, req.AuctionID),
			AuctionID: req.AuctionID,
			IsRead:    false,
		}
		if err := a.notRepo.CreateNotification(ctx, not); err != nil {
			return nil, fmt.Errorf("CreateNotification failed: %v", err)
		}
	}

	return &models.BidResponse{
		AuctionID: req.AuctionID,
		BidderID:  req.BidderID,
		BidAmount: req.BidAmount,
		TimeStamp: savedBid.CreatedAt,
	}, nil
}

// CloseAuction method in your AuctionService
// Now accepts the authenticated userID for authorization checks.
func (a *AuctionService) CloseAuction(ctx context.Context, auctionID string, requestingUserID string) (*models.WinnerResponse, error) {
	auction, err := a.repo.GetAuctionById(ctx, auctionID)
	if err != nil {
		if errors.Is(err, errs.ErrAuctionNotFound) {
			return nil, errs.ErrAuctionNotFound
		}
		return nil, errors.New("failed to retrieve auction for closing")
	}

	// Authorization
	if auction.SellerID != requestingUserID {
		return nil, errs.ErrPermissionDenied
	}

	// Prevent re-closing
	if auction.Status == "closed" {
		return nil, errs.ErrAuctionAlreadyClosed
	}

	// Update status
	auction.Status = "closed"
	if err := a.repo.CloseAuction(ctx, auction.Status, auctionID); err != nil {
		return nil, errs.ErrFailedToUpdateAuction
	}

	// Determine the winner: highest bid
	var winnerID string
	if auction.CurrentPrice > auction.StartingPrice {
		highestBid, err := a.bidRepo.GetHighestBid(ctx, auctionID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("failed to retrieve highest bid during auction close")
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
			UserID:    winnerID,
			Message:   fmt.Sprintf("Congratulations! You won the auction: %s, auctionId: %s", auction.Title, auctionID),
			AuctionID: auctionID,
			IsRead:    false,
		}
		if err := a.notRepo.CreateNotification(ctx, not); err != nil {
			return nil, fmt.Errorf("CreateNotification failed: %v", err)
		}
	}

	// Notify other bidders
	bidders, err := a.bidRepo.GetAllBidderIDsForAuction(ctx, auctionID)
	if err != nil {
		return nil, errors.New("failed to retrieve all bidders for auction close")
	}

	uniqueBidders := make(map[string]struct{})
	for _, id := range bidders {
		if id != winnerID && id != requestingUserID {
			if _, seen := uniqueBidders[id]; !seen {
				uniqueBidders[id] = struct{}{}
				a.notifications <- &models.NotificationEvent{
					Type:      models.NotificationAuctionEnded,
					UserID:    id,
					Message:   fmt.Sprintf("Auction %s has ended. You did not win.", auction.Title),
					AuctionID: auctionID,
					TimeStamp: time.Now(),
				}
			}
		}

		not := &store.Notification{
			UserID:    id,
			Message:   fmt.Sprintf("Auction %s has ended. You did not win.", auction.Title),
			AuctionID: auctionID,
			IsRead:    false,
		}
		if err := a.notRepo.CreateNotification(ctx, not); err != nil {
			return nil, fmt.Errorf("CreateNotification failed: %v", err)
		}
	}

	// Notify WebSocket listeners
	a.auctionUpdates <- &models.AuctionUpdateEvent{
		EventType:    models.AuctionEnded,
		ID:           auctionID,
		Type:         auction.Type,
		Status:       auction.Status,
		SellerID:     auction.SellerID,
		CurrentPrice: auction.CurrentPrice,
		TimeStamp:    time.Now(),
	}

	// Delete bids
	if err := a.bidRepo.DeleteBidsByAuction(ctx, auctionID); err != nil {
		return nil, errs.ErrFailedToDeleteBids
	}

	// Delete notifications
	if err := a.notRepo.DeleteNotificationByAuction(ctx, auctionID); err != nil {
		return nil, errs.ErrFailedToDeleteNotifications
	}

	res := &models.WinnerResponse{
		WinnerID:   winnerID,
		WinningBid: auction.CurrentPrice,
		Status:     auction.Status,
	}

	return res, nil
}
