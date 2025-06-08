package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/puremike/online_auction_api/internal/config"
	"github.com/puremike/online_auction_api/internal/models"
	"github.com/puremike/online_auction_api/internal/store"
	"github.com/puremike/online_auction_api/internal/utils"
)

type UserService struct {
	repo store.UserRepository
	app  *config.Application
}

func NewUserService(repo store.UserRepository, app *config.Application) *UserService {
	return &UserService{
		repo: repo,
		app:  app,
	}
}

func (u *UserService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {

	// validation
	if len(user.Username) < 6 || user.Email == "" || len(user.Password) < 8 || len(user.FullName) < 8 || len(user.Location) < 8 {
		return nil, fmt.Errorf("invalid user details")
	}

	hashedPassword, err := utils.HashedPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword

	createdUser, err := u.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (u *UserService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {

	if req.Email == "" || req.Password == "" {
		return &models.LoginResponse{}, fmt.Errorf("email or password cannot be empty")
	}

	user, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			return &models.LoginResponse{}, fmt.Errorf("user not found")
		}
		return &models.LoginResponse{}, err
	}

	if err := utils.CompareHashedPassword(user.Password, req.Password); err != nil {
		return &models.LoginResponse{}, fmt.Errorf("invalid credentials")
	}

	claims := jwt.MapClaims{
		"sub":     user.ID,
		"isAdmin": user.IsAdmin,
		"iss":     u.app.AppConfig.AuthConfig.Iss,
		"aud":     u.app.AppConfig.AuthConfig.Aud,
		"iat":     time.Now().Unix(),
		"nbf":     time.Now().Unix(),
		"exp":     time.Now().Add(u.app.AppConfig.AuthConfig.TokenExp).Unix(),
	}

	token, err := u.app.JwtAUth.GenerateToken(claims)
	if err != nil {
		return &models.LoginResponse{}, fmt.Errorf("failed to generate token")
	}

	refreshToken, err := u.app.JwtAUth.GenerateRefreshToken()
	if err != nil {
		return &models.LoginResponse{}, fmt.Errorf("failed to generate refresh token")
	}

	if err := u.repo.StoreRefreshToken(ctx, user.ID, refreshToken, time.Now().Add(u.app.AppConfig.AuthConfig.RefreshTokenExp)); err != nil {
		return &models.LoginResponse{}, fmt.Errorf("failed to store refresh token")
	}

	return &models.LoginResponse{ID: user.ID, Username: user.Username, Token: token, RefreshToken: refreshToken}, nil
}

func (u *UserService) UserProfile(ctx context.Context, username string) (*models.UserResponse, error) {

	user, err := u.repo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, store.ErrUserNotFound) {
			return &models.UserResponse{}, fmt.Errorf("user not found")
		}
		return &models.UserResponse{}, err
	}
	return &models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		Location:  user.Location,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (u *UserService) Refresh(ctx context.Context, refreshToken string) (string, error) {

	userId, err := u.repo.ValidateRefreshToken(ctx, refreshToken)
	if err != nil || userId == "" {

		if errors.Is(err, store.ErrTokenNotFound) {
			return "", fmt.Errorf("refresh token not found")
		}
		return "", fmt.Errorf("invalid refresh token")
	}

	// Generate new JWT
	claims := jwt.MapClaims{
		"sub": userId,
		"iss": u.app.AppConfig.AuthConfig.Iss,
		"aud": u.app.AppConfig.AuthConfig.Aud,
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"exp": time.Now().Add(u.app.AppConfig.AuthConfig.TokenExp).Unix(),
	}

	newToken, err := u.app.JwtAUth.GenerateToken(claims)
	if err != nil {
		return "", fmt.Errorf("failed to generate new token")
	}

	return newToken, nil
}
