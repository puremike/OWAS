package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/puremike/online_auction_api/contexts"
	"github.com/puremike/online_auction_api/internal/config"
	"github.com/puremike/online_auction_api/internal/errs"
	"github.com/puremike/online_auction_api/internal/models"
	"github.com/puremike/online_auction_api/internal/services"
)

type AuctionHandler struct {
	service *services.AuctionService
	app     *config.Application
}

func NewAuctionHandler(service *services.AuctionService, app *config.Application) *AuctionHandler {
	return &AuctionHandler{
		service: service,
		app:     app,
	}
}

// CreateAuction godoc
//
//	@Summary		Create Auction
//	@Description	Creates a new auction.
//	@Tags			Auctions
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.CreateAuctionRequest	true	"Auction create payload"
//	@Success		201			{object}	models.CreateAuctionResponse	"Created auction"
//	@Failure		400			{object}	gin.H						"Bad Request - invalid input"
//	@Failure		401			{object}	gin.H						"Unauthorized - user not authenticated"
//	@Failure		500			{object}	gin.H						"Internal Server Error - failed to create auction"
//	@Router			/auctions [post]
//
//	@Security		jwtCookieAuth
func (a *AuctionHandler) CreateAuction(c *gin.Context) {

	var payload models.CreateAuctionRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authUser, err := contexts.GetUserFromContext(c)
	if authUser == nil || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	startDate, err := time.Parse("2006-01-02", payload.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start time"})
		return
	}
	endDate, err := time.Parse("2006-01-02", payload.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end time"})
		return
	}

	auction := &models.Auction{
		SellerID:      authUser.ID,
		Title:         payload.Title,
		Description:   payload.Description,
		StartingPrice: payload.StartingPrice,
		CurrentPrice:  payload.StartingPrice,
		Type:          strings.ToLower(payload.Type),
		Status:        "open",
		StartTime:     startDate,
		EndTime:       endDate,
		ImagePath:     payload.ImagePath,
		Category:      payload.Category,
		IsPaid:        false,
	}

	createdAuction, err := a.service.CreateAuction(c.Request.Context(), auction)
	if err != nil {
		errs.MapServiceErrors(c, err)
		return
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

	c.JSON(http.StatusCreated, res)
}

// UpdateAuction godoc
//
//	@Summary		Update Auction
//	@Description	Allows a seller to update an auction they have created.
//	@Tags			Auctions
//	@Accept			json
//	@Produce		json
//	@Param			auctionID	path		string						true	"ID of the auction to update"
//	@Param			payload		body		models.CreateAuctionRequest	true	"Auction update payload"
//	@Success		201			{object}	models.Auction				"Updated auction"
//	@Failure		400			{object}	gin.H						"Bad Request - invalid input"
//	@Failure		401			{object}	gin.H						"Unauthorized - user not authenticated"
//	@Failure		404			{object}	gin.H						"NotFound - auction not found"
//	@Failure		500			{object}	gin.H						"Internal Server Error - failed to update auction"
//	@Router			/auctions/{auctionID} [put]
//
//	@Security		jwtCookieAuth
func (a *AuctionHandler) UpdateAuction(c *gin.Context) {

	var payload models.CreateAuctionRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authUser, err := contexts.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	existingAuction, err := contexts.GetAuctionFromContext(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "auction not found"})
		return
	}

	if authUser.ID != existingAuction.SellerID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	startDate, err := time.Parse("2006-01-02", payload.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start time"})
		return
	}
	endDate, err := time.Parse("2006-01-02", payload.EndTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end time"})
		return
	}

	auction := &models.Auction{
		SellerID:      authUser.ID,
		Title:         payload.Title,
		Description:   payload.Description,
		StartingPrice: payload.StartingPrice,
		CurrentPrice:  payload.StartingPrice,
		Type:          strings.ToLower(payload.Type),
		Status:        "open",
		StartTime:     startDate,
		EndTime:       endDate,
	}

	updatedAuction, err := a.service.UpdateAuction(c.Request.Context(), auction, existingAuction.ID)
	if err != nil {
		errs.MapServiceErrors(c, err)
		return
	}

	c.JSON(http.StatusCreated, updatedAuction)
}

// DeleteAuction godoc
//
//	@Summary		Delete Auction
//	@Description	Allows a seller to delete an auction they have created.
//	@Tags			Auctions
//	@Accept			json
//	@Produce		json
//	@Param			auctionID	path		string			true	"ID of the auction to delete"
//	@Success		200			{object}	models.Auction	"Deleted auction"
//	@Failure		401			{object}	gin.H			"Unauthorized - user not authenticated"
//	@Failure		404			{object}	gin.H			"NotFound - auction not found"
//	@Failure		500			{object}	gin.H			"Internal Server Error - failed to delete auction"
//	@Router			/auctions/{auctionID} [delete]
//
//	@Security		jwtCookieAuth
func (a *AuctionHandler) DeleteAuction(c *gin.Context) {

	authUser, err := contexts.GetUserFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	existingAuction, err := contexts.GetAuctionFromContext(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "auction not found"})
		return
	}

	if authUser.ID != existingAuction.SellerID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	deletedAuction, err := a.service.DeleteAuction(c.Request.Context(), existingAuction.ID)
	if err != nil {
		errs.MapServiceErrors(c, err)
		return
	}

	c.JSON(http.StatusOK, deletedAuction)
}

// DeleteAuction godoc
//
//	@Summary		Admin Delete Auction
//	@Description	Allows an admin to delete any auction.
//	@Tags			Auctions
//	@Accept			json
//	@Produce		json
//	@Param			auctionID	path		string			true	"ID of the auction to delete"
//	@Success		200			{object}	models.Auction	"Deleted auction"
//	@Failure		401			{object}	gin.H			"Unauthorized - user not authenticated"
//	@Failure		404			{object}	gin.H			"NotFound - auction not found"
//	@Failure		500			{object}	gin.H			"Internal Server Error - failed to delete auction"
//	@Router			/admin/auctions/{auctionID} [delete]
//
//	@Security		jwtCookieAuth
func (a *AuctionHandler) AdminDeleteAuction(c *gin.Context) {
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

	deletedAuction, err := a.service.DeleteAuction(c.Request.Context(), existingAuction.ID)
	if err != nil {
		errs.MapServiceErrors(c, err)
		return
	}

	c.JSON(http.StatusOK, deletedAuction)
}

// GetAuctionById godoc
//
//	@Summary		Get Auction By ID
//	@Description	Retrieves an auction by its ID.
//	@Tags			Auctions
//	@Accept			json
//	@Produce		json
//	@Param			auctionID	path		string			true	"ID of the auction to retrieve"
//	@Success		200			{object}	models.Auction	"Retrieved auction"
//	@Failure		401			{object}	gin.H			"Unauthorized - user not authenticated"
//	@Failure		404			{object}	gin.H			"NotFound - auction not found"
//	@Failure		500			{object}	gin.H			"Internal Server Error - failed to retrieve auction"
//	@Router			/auctions/{auctionID} [get]
//
//	@Security		jwtCookieAuth
func (a *AuctionHandler) GetAuctionById(c *gin.Context) {

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

	auction, err := a.service.GetAuctionById(c.Request.Context(), existingAuction.ID)
	if err != nil {
		errs.MapServiceErrors(c, err)
		return
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

	c.JSON(http.StatusOK, res)
}

// AdminGetAuctions godoc
//
//	@Summary		Retrieve All Auctions (Admin)
//	@Description	Fetches a list of all auctions with admin privileges.
//	@Tags			Auctions
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		models.CreateAuctionResponse	"List of auctions"
//	@Failure		401	{object}	gin.H							"Unauthorized - user not authenticated"
//	@Failure		500	{object}	gin.H							"Internal Server Error - failed to retrieve auctions"
//	@Router			/auctions [get]
//
//	@Security		jwtCookieAuth
func (a *AuctionHandler) GetAuctions(c *gin.Context) {

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	filter := &models.AuctionFilter{
		Type:     c.Query("type"),
		Category: c.Query("category"),
		Status:   c.Query("status"),
		StartingPrice: func() float64 {
			p, _ := strconv.ParseFloat(c.Query("starting_price"), 64)
			return p
		}(),
	}

	authUser, err := contexts.GetUserFromContext(c)
	if authUser == nil || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	auctions, err := a.service.GetAuctions(c.Request.Context(), limit, offset, filter)
	if err != nil {
		errs.MapServiceErrors(c, err)
		return
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

	c.JSON(http.StatusOK, res)
}

// GetMyWonAuctions godoc
//
//	@Summary		Get My Won Auctions
//	@Description	Retrieves a list of auctions the user has won.
//	@Tags			Auctions
//	@Accept			json
//	@Produce		json
//	@Success		200			{array}		models.CreateAuctionResponse	"List of auctions"
//	@Failure		401			{object}	gin.H							"Unauthorized - user not authenticated"
//	@Failure		500			{object}	gin.H							"Internal Server Error - failed to retrieve auctions"
//	@Router			/auctions/won [get]
//
//	@Security		jwtCookieAuth
func (a *AuctionHandler) GetMyWonAuctions(c *gin.Context) {
	authUser, err := contexts.GetUserFromContext(c)
	if authUser == nil || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	auctions, err := a.service.GetWonAuctionsByWinnerID(c.Request.Context(), authUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load won auctions"})
		return
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

	c.JSON(http.StatusOK, res)
}
