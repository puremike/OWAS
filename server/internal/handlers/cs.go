package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/puremike/online_auction_api/contexts"
	"github.com/puremike/online_auction_api/internal/models"
	"github.com/puremike/online_auction_api/internal/services"
)

type CSHandler struct {
	csService services.CSServiceInterface
}

func NewCSHandler(csService services.CSServiceInterface) *CSHandler {
	return &CSHandler{
		csService: csService,
	}
}

// ContactSupport godoc
//
//	@Summary		Contact Support
//	@Description	Send a message to the support team
//	@Tags			Contact Support
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.ContactSupportReq	true	"Contact Support payload"
//	@Success		200		{object}	models.SupportRes
//	@Failure		400		{object}	gin.H	"Bad Request - invalid input"
//	@Failure		401		{object}	gin.H	"Unauthorized - user not authenticated"
//	@Failure		500		{object}	gin.H	"Internal Server Error - failed to contact support"
//	@Router			/contact-support [post]
//
//	@Security		jwtCookieAuth
func (cs *CSHandler) ContactSupport(c *gin.Context) {

	var payload models.ContactSupportReq
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authUser, err := contexts.GetUserFromContext(c)
	if authUser == nil || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	createdSupport := &models.ContactSupport{
		UserID:  authUser.ID,
		Subject: payload.Subject,
		Message: payload.Message,
	}

	res, err := cs.csService.ContactSupport(c, createdSupport)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
