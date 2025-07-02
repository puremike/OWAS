package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/puremike/online_auction_api/contexts"
	"github.com/puremike/online_auction_api/internal/config"
	"github.com/puremike/online_auction_api/internal/errs"
	"github.com/puremike/online_auction_api/internal/models"
	"github.com/puremike/online_auction_api/internal/services"
)

type UserHandler struct {
	service *services.UserService
	app     *config.Application
}

func NewUserHandler(service *services.UserService, app *config.Application) *UserHandler {
	return &UserHandler{
		service: service,
		app:     app,
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
		errs.MapServiceErrors(c, err)
		return
	}

	c.JSON(http.StatusCreated, models.UserResponse{
		ID:        createdUser.ID,
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		FullName:  createdUser.FullName,
		Location:  createdUser.Location,
		CreatedAt: createdUser.CreatedAt,
	})
}

// Login godoc
//
//	@Summary		Login User
//	@Description	Authenticates a user using email and password.
//	@Description	Upon successful authentication, a short-lived **JWT (access token)** is set as an `HttpOnly` cookie named `jwt`.
//	@Description	A long-lived **refresh token** is also set as an `HttpOnly` cookie named `refresh_token`.
//	@Description	Both cookies are crucial for maintaining user session and subsequent authenticated requests.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.LoginRequest	true	"Login credentials"
//	@Success		200		{object}	string				"login successful"
//	@Success		200		{header}	string				Set-Cookie	"Two HttpOnly cookies are set: 'jwt' (access token) and 'refresh_token' (refresh token)."
//	@Failure		400		{object}	gin.H				"Bad Request - invalid input"
//	@Failure		401		{object}	gin.H				"Unauthorized - invalid credentials"
//	@Failure		500		{object}	gin.H				"Internal Server Error"
//	@Router			/login [post]
func (u *UserHandler) Login(c *gin.Context) {

	var payload models.LoginRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := u.service.Login(c.Request.Context(), &payload)
	if err != nil {
		errs.MapServiceErrors(c, err)
		return
	}

	// c.SetCookie("jwt", user.Token, int(u.app.AppConfig.AuthConfig.TokenExp.Seconds()), "/", "", false, true)
	// c.SetCookie("refresh_token", user.RefreshToken, int(u.app.AppConfig.AuthConfig.RefreshTokenExp.Seconds()), "/", "", false, true)
	// c.SetSameSite(http.SameSiteStrictMode)

	u.setJwtCookie(c, user)

	c.JSON(http.StatusOK, "login successful")
}

// AdminLogin godoc
//
//	@Summary		Login Admin
//	@Description	Authenticates an admin user using email and password.
//	@Description	Upon successful authentication, a short-lived **JWT (access token)** is set as an `HttpOnly` cookie named `jwt`.
//	@Description	A long-lived **refresh token** is also set as an `HttpOnly` cookie named `refresh_token`.
//	@Description	Both cookies are crucial for maintaining user session and subsequent authenticated requests.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.LoginRequest	true	"Login credentials"
//	@Success		200		{object}	string				"login successful"
//	@Success		200		{header}	string				Set-Cookie	"Two HttpOnly cookies are set: 'jwt' (access token) and 'refresh_token' (refresh token)."
//	@Failure		400		{object}	gin.H				"Bad Request - invalid input"
//	@Failure		401		{object}	gin.H				"Unauthorized - invalid credentials"
//	@Failure		500		{object}	gin.H				"Internal Server Error"
//	@Router			/admin/login [post]
func (u *UserHandler) AdminLogin(c *gin.Context) {

	var payload models.LoginRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := u.app.Store.Users.GetUserByEmail(c.Request.Context(), payload.Email)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
			return
		}
		return
	}

	if !existingUser.IsAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user is not an admin"})
		return
	}
	user, err := u.service.Login(c.Request.Context(), &payload)
	if err != nil {
		errs.MapServiceErrors(c, err)
		return
	}

	u.setJwtCookie(c, user)

	c.JSON(http.StatusOK, "login successful")
}

func (u *UserHandler) setJwtCookie(c *gin.Context, user *models.LoginResponse) {

	isTrue := true
	sameSite := http.SameSiteNoneMode

	if u.app.AppConfig.Env == "development" {
		isTrue = false
		sameSite = http.SameSiteLaxMode
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "jwt",
		Value:    user.Token,
		Path:     "/",
		MaxAge:   int(u.app.AppConfig.AuthConfig.TokenExp.Seconds()),
		HttpOnly: true,
		Secure:   isTrue, // change to true in production
		SameSite: sameSite,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    user.RefreshToken,
		Path:     "/",
		MaxAge:   int(u.app.AppConfig.AuthConfig.RefreshTokenExp.Seconds()),
		HttpOnly: true,
		Secure:   isTrue, // change to true in production
		SameSite: sameSite,
	})
}

// Logout godoc
//
//	@Summary		Logout User
//	@Description	Clears the user's authentication cookies, effectively logging them out.
//	@Tags			Users
//	@Success		200	{object}	gin.H	"Logout successful"
//	@Router			/logout [post]
//
//	@Security		jwtCookieAuth
func (u *UserHandler) Logout(c *gin.Context) {
	c.SetCookie("jwt", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)
	c.SetSameSite(http.SameSiteStrictMode)

	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}

// RefreshToken godoc
//
//	@Summary		Refresh JWT Token
//	@Description	Refreshes the JWT access token using a valid refresh token.
//	@Description	If the refresh token is valid, a new JWT is generated and set as an `HttpOnly` cookie.
//	@Description	A valid refresh token must be provided as an `HttpOnly` cookie named `refresh_token`.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	gin.H	"Token refreshed successfully"
//	@Failure		401	{object}	gin.H	"Unauthorized - Refresh token not found or invalid"
//	@Failure		500	{object}	gin.H	"Internal Server Error - Failed to generate new token"
//	@Router			/refresh [post]
func (u *UserHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")

	if refreshToken == "" || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token not found"})
		return
	}

	newToken, err := u.service.Refresh(c.Request.Context(), refreshToken)
	if err != nil {
		errs.MapServiceErrors(c, err)
		return
	}

	c.SetCookie("jwt", newToken, int(u.app.AppConfig.AuthConfig.TokenExp.Seconds()), "/", "", false, true)
	c.SetSameSite(http.SameSiteStrictMode)

	c.JSON(http.StatusOK, gin.H{"message": "token refreshed successfully"})
}

// UserProfile godoc
//
//	@Summary		Get User Profile
//	@Description	Retrieves the profile of the user associated with the access token.
//	@Description	Access token must be provided as an `HttpOnly` cookie named `jwt`.
//	@Tags			Users
//	@Accept			json
//	@Param			username	path	string	true	"Username of the user to retrieve profile for"
//	@Produce		json
//	@Success		200	{object}	models.UserResponse
//	@Failure		401	{object}	gin.H	"Unauthorized - invalid or expired token"
//	@Failure		404	{object}	gin.H	"User not found"
//	@Failure		500	{object}	gin.H	"Internal Server Error"
//	@Router			/{username} [get]
//
//	@Security		jwtCookieAuth
func (u *UserHandler) UserProfile(c *gin.Context) {

	authUser, err := contexts.GetUserFromContext(c)
	if authUser == nil || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	user, err := u.service.UserProfile(c.Request.Context(), username)
	if err != nil {
		errs.MapServiceErrors(c, err)
		return
	}

	res := models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		Location:  user.Location,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusOK, res)
}

// MeProfile godoc
//
//	@Summary		Get authenticated user's profile
//	@Description	Retrieves the profile details of the authenticated user.
//	@Tags			Users
//	@Produce		json
//	@Success		200	{object}	models.UserResponse
//	@Failure		401	{object}	gin.H	"Unauthorized - user not authenticated"
//	@Failure		500	{object}	gin.H	"Internal Server Error - failed to retrieve profile"
//	@Router			/me [get]
//
//	@Security		jwtCookieAuth

func (u *UserHandler) MeProfile(c *gin.Context) {
	authUser, err := contexts.GetUserFromContext(c)
	if authUser == nil || err != nil {
		log.Printf("ERROR: Could not get authenticated user from context: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := u.service.MeProfile(c.Request.Context(), authUser.ID)
	if err != nil {
		log.Printf("failed to retrieve user profile for ID %s: %v", authUser.ID, err)
		errs.MapServiceErrors(c, err)
		return
	}

	// Respond with the UserResponse model
	res := models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		Location:  user.Location,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusOK, res)
}

// UpdateProfile godoc
//
//	@Summary		Update user profile
//	@Description	Allows an authenticated user to update their profile details such as username, email, full name, and location.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			payload		body		models.UserProfileUpdateRequest	true	"Profile update payload"
//	@Param			username	path		string							true	"Username of the user to update profile for"
//	@Success		201			{object}	string							"Profile updated successfully"
//	@Failure		400			{object}	gin.H							"Bad Request - invalid input"
//	@Failure		401			{object}	gin.H							"Unauthorized - user not authenticated"
//	@Failure		500			{object}	gin.H							"Internal Server Error - failed to update profile"
//	@Router			/{username}/update-profile [put]
//
//	@Security		jwtCookieAuth
func (u *UserHandler) UpdateProfile(c *gin.Context) {

	var payload models.UserProfileUpdateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authUser, err := contexts.GetUserFromContext(c)
	if authUser == nil || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
		return
	}

	currentUser, err := u.app.Store.Users.GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user"})
		return
	}

	if authUser.ID != currentUser.ID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized to update this profile"})
		return
	}

	user := &models.User{
		Username: payload.Username,
		Email:    payload.Email,
		FullName: payload.FullName,
		Location: payload.Location,
	}

	msg, err := u.service.UpdateProfile(c.Request.Context(), user, authUser.ID)
	if err != nil {
		errs.MapServiceErrors(c, err)
		return
	}

	c.JSON(http.StatusCreated, msg)
}

// ChangePassword godoc
//
//	@Summary		Change User Password
//	@Description	Allows an authenticated user to change their password.
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Param			payload		body		models.PasswordUpdateRequest	true	"Password update payload"
//	@Param			username	path		string							true	"Username of the user to change password for"
//	@Success		201			{object}	string							"Password changed successfully"
//	@Failure		400			{object}	gin.H							"Bad Request - invalid input"
//	@Failure		401			{object}	gin.H							"Unauthorized - user not authenticated"
//	@Failure		500			{object}	gin.H							"Internal Server Error - failed to change password"
//	@Router			/change-password [put]
//
//	@Security		jwtCookieAuth
func (u *UserHandler) ChangePassword(c *gin.Context) {

	var payload models.PasswordUpdateRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authUser, err := contexts.GetUserFromContext(c)
	if authUser == nil || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized to update profile"})
		return
	}

	// username := c.Param("username")
	// if username == "" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "username is required"})
	// 	return
	// }

	// currentUser, err := u.app.Store.Users.GetUserByUsername(c.Request.Context(), username)
	// if err != nil {
	// 	if errors.Is(err, errs.ErrUserNotFound) {
	// 		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
	// 		return
	// 	}
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user"})
	// 	return
	// }

	// if authUser.ID != currentUser.ID {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized to update this profile"})
	// 	return
	// }

	msg, err := u.service.ChangePassword(c.Request.Context(), &payload, authUser.ID)
	if err != nil {
		errs.MapServiceErrors(c, err)
		return
	}

	c.JSON(http.StatusCreated, msg)
}

// AdminGetUsers godoc
//
//	@Summary		Get all users
//	@Description	Retrieves all users
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]models.UserResponse
//	@Failure		401	{object}	gin.H	"Unauthorized - user not authenticated"
//	@Failure		500	{object}	gin.H	"Internal Server Error - failed to retrieve users"
//	@Router			/admin/users [get]
//
//	@Security		jwtCookieAuth
func (u *UserHandler) AdminGetUsers(c *gin.Context) {

	authUser, err := contexts.GetUserFromContext(c)
	if authUser == nil || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	users, err := u.service.GetUsers(c.Request.Context())
	if err != nil {
		errs.MapServiceErrors(c, err)
		return
	}

	res := &[]models.UserResponse{}
	for _, user := range *users {
		*res = append(*res, models.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			FullName:  user.FullName,
			Location:  user.Location,
			CreatedAt: user.CreatedAt,
		})
	}

	c.JSON(http.StatusOK, res)
}

func (u *UserHandler) DeleteUser(c *gin.Context) {

	authUser, err := contexts.GetUserFromContext(c)
	if authUser == nil || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	msg, err := u.service.DeleteUser(c.Request.Context(), authUser.ID)
	if err != nil {
		errs.MapServiceErrors(c, err)
		return
	}

	c.JSON(http.StatusOK, msg)
}

func (u *UserHandler) AdminDeleteUser(c *gin.Context) {

	id := c.Param("userID")

	authUser, err := contexts.GetUserFromContext(c)
	if authUser == nil || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	msg, err := u.service.DeleteUser(c.Request.Context(), id)
	if err != nil {
		errs.MapServiceErrors(c, err)
		return
	}

	c.JSON(http.StatusOK, msg)
}
