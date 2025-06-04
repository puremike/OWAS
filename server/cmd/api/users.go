package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/puremike/online_auction_api/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type createUserRequest struct {
	Username        string `json:"username" binding:"required,min=6,max=32"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,passwd"`
	ConfirmPassword string `json:"confirm_password" binding:"required,passwd"`
	FullName        string `json:"full_name" binding:"required,min=6,max=32"`
	Location        string `json:"location" binding:"required,min=6,max=42"`
	IsAdmin         bool   `json:"is_admin" binding:"required"`
}

type userResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	Location  string `json:"location"`
	CreatedAt string `json:"created_at"`
}

// CreateUser godoc
//
//	@Summary		Create user
//	@Description	Create a new user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		createUserRequest	true	"User payload"
//	@Success		201		{object}	userResponse
//	@Failure		400		{object}	error
//	@Failure		404		{object}	error
//	@Failure		500		{object}	error
//	@Router			/auth/signup [post]
func (app *application) registerUser(c *gin.Context) {

	var payload createUserRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if payload.Password != payload.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	user := &models.User{
		Username: payload.Username,
		Email:    payload.Email,
		Password: string(hashedPassword),
		FullName: payload.FullName,
		Location: payload.Location,
		IsAdmin:  payload.IsAdmin,
	}

	createdUser, err := app.store.Users.CreateUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "problem registering user"})
		return
	}

	c.JSON(http.StatusCreated, userResponse{
		ID:        createdUser.ID,
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		FullName:  createdUser.FullName,
		Location:  createdUser.Location,
		CreatedAt: createdUser.CreatedAt.Format(time.RFC3339),
	})
}
