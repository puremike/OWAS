package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/puremike/online_auction_api/internal/models"
	"github.com/puremike/online_auction_api/internal/services"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// CreateUser godoc
//
//	@Summary		Create user
//	@Description	Create a new user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.CreateUserRequest	true	"User payload"
//	@Success		201		{object}	models.UserResponse
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/signup [post]
func (u *UserHandler) RegisterUser(c *gin.Context) {

	var payload models.CreateUserRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if payload.Password != payload.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
		return
	}

	user := &models.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
		FullName: payload.FullName,
		Location: payload.Location,
		IsAdmin:  false,
	}

	createdUser, err := u.service.CreateUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "problem registering user"})
		return
	}

	c.JSON(http.StatusCreated, models.UserResponse{
		ID:        createdUser.ID,
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		FullName:  createdUser.FullName,
		Location:  createdUser.Location,
		CreatedAt: createdUser.CreatedAt.Format(time.RFC3339),
	})
}
