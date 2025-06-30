package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/puremike/online_auction_api/contexts"
	"github.com/puremike/online_auction_api/internal/errs"
	"github.com/puremike/online_auction_api/internal/models"
)

type PlaceBidRequest struct {
	BidAmount float64 `json:"bidAmount" binding:"required"`
}

// PlaceBids godoc
//
//	@Summary		Place a Bid
//	@Description	Allows a user to place a bid on an existing auction.
//	@Tags			Bids
//	@Accept			json
//	@Produce		json
//	@Param			auctionID	path		string				true	"ID of the auction to bid on"
//	@Param			payload		body		PlaceBidRequest		true	"Bid request payload"
//	@Success		200			{object}	models.BidResponse	"Bid placed successfully"
//	@Failure		400			{object}	gin.H				"Bad Request - invalid input"
//	@Failure		401			{object}	gin.H				"Unauthorized - user not authenticated or authorized"
//	@Failure		404			{object}	gin.H				"NotFound - auction not found"
//	@Failure		500			{object}	gin.H				"Internal Server Error - failed to place bid"
//	@Router			/auctions/{auctionID}/bids [post]
//
//	@Security		jwtCookieAuth
func (a *AuctionHandler) PlaceBids(c *gin.Context) {
	var payload PlaceBidRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authUser, err := contexts.GetUserFromContext(c)
	if authUser == nil || err != nil || authUser.IsAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	existingAuction, err := contexts.GetAuctionFromContext(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "auction not found"})
		return
	}

	if existingAuction.SellerID == authUser.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	auction := &models.PlaceBidRequest{
		AuctionID: c.Param("auctionID"),
		BidderID:  authUser.ID,
		BidAmount: payload.BidAmount,
	}

	bid, err := a.service.PlaceBid(c.Request.Context(), auction)
	if err != nil {
		errs.MapServiceErrors(c, err)
		return
	}

	c.JSON(http.StatusOK, bid)
}

// CloseAuction godoc
//
//	@Summary		Close Auction
//	@Description	Allows the seller of an auction to close the auction. The auction must be open.
//	@Tags			Auctions
//	@Accept			json
//	@Produce		json
//	@Param			auctionID	path		string	true	"ID of the auction to close"
//	@Success		200			{object}	gin.H	"Closed auction message"
//	@Failure		401			{object}	gin.H	"Unauthorized - user not authenticated or not the seller"
//	@Failure		404			{object}	gin.H	"Auction not found"
//	@Failure		500			{object}	gin.H	"Internal Server Error - failed to close auction"
//	@Router			/auctions/{auctionID}/close [post]
//
//	@Security		jwtCookieAuth
func (a *AuctionHandler) CloseAuction(c *gin.Context) {

	authUser, err := contexts.GetUserFromContext(c)
	if authUser == nil || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	existingAuction, err := contexts.GetAuctionFromContext(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "auction not found"})
		return
	}

	if existingAuction.SellerID != authUser.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	response, err := a.service.CloseAuction(c.Request.Context(), existingAuction.ID, existingAuction.SellerID)
	if err != nil {
		errs.MapServiceErrors(c, err)
		return
	}

	res := &models.WinnerResponse{
		WinnerID:   response.WinnerID,
		WinningBid: response.WinningBid,
		Status:     response.Status,
	}

	c.JSON(http.StatusOK, res)
}
