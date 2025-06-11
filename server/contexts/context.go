package contexts

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/puremike/online_auction_api/internal/models"
)

func GetUserFromContext(c *gin.Context) (*models.User, error) {
	userContext, exist := c.Get("user")
	if !exist {
		return nil, errors.New("user not found in the context")
	}

	user, ok := userContext.(*models.User)
	if !ok {
		return nil, errors.New("user context is not of type *models.User")
	}

	return user, nil
}

func GetAuctionFromContext(c *gin.Context) (*models.Auction, error) {
	auctionContext, exist := c.Get("auction")
	if !exist {
		return nil, errors.New("auction not found in the context")
	}

	auction, ok := auctionContext.(*models.Auction)
	if !ok {
		return nil, errors.New("auction context is not of type *models.Auction")
	}

	return auction, nil
}
