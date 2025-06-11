package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/puremike/online_auction_api/internal/services"
)

type AuctionHandler struct {
	service *services.AuctionService
}

func NewAuctionHandler(service *services.AuctionService) *AuctionHandler {
	return &AuctionHandler{
		service: service,
	}
}

func (a *AuctionHandler) CreateAuction(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"message": "Auction created successfully"})

}
