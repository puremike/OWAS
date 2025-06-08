package contexts

import (
	"github.com/gin-gonic/gin"
	"github.com/puremike/online_auction_api/internal/models"
)

func GetUserFromContext(c *gin.Context) *models.User {
	userContext, exist := c.Get("user")
	if !exist {
		return &models.User{}
	}

	user, ok := userContext.(*models.User)
	if !ok {
		return &models.User{}
	}
	return user
}
