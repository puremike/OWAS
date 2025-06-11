package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/puremike/online_auction_api/internal/config"
	"github.com/puremike/online_auction_api/internal/errs"
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

const QueryDefaultContext = 5 * time.Second

func (u *UserService) CreateUser(ctx context.Context, user *models.User) (*models.UserResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryDefaultContext)
	defer cancel()

	// validation
	if len(user.Username) < 6 || user.Email == "" || len(user.Password) < 8 || len(user.FullName) < 8 || len(user.Location) < 8 {
		return nil, errs.ErrInvalidUserDetails
	}

	hashedPassword, err := utils.HashedPassword(user.Password)
	if err != nil {
		return nil, errs.ErrFailedToHashPassword
	}

	us := &models.User{
		Username: user.Username,
		Email:    user.Email,
		Password: hashedPassword,
		FullName: user.FullName,
		Location: user.Location,
	}

	createdUser, err := u.repo.CreateUser(ctx, us)
	if err != nil {
		return nil, errs.ErrFailedToCreateUser
	}

	res := &models.UserResponse{
		ID:        createdUser.ID,
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		FullName:  createdUser.FullName,
		Location:  createdUser.Location,
		CreatedAt: createdUser.CreatedAt.Format(time.RFC3339),
	}

	return res, nil
}

func (u *UserService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryDefaultContext)
	defer cancel()

	if req.Email == "" || req.Password == "" {
		return &models.LoginResponse{}, errs.ErrEmailPasswordRequired
	}

	user, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return &models.LoginResponse{}, errs.ErrUserNotFound
		}
		return &models.LoginResponse{}, fmt.Errorf("failed to retrieve user: %w", err)
	}

	if err := utils.CompareHashedPassword(user.Password, req.Password); err != nil {
		return &models.LoginResponse{}, errs.ErrInvalidCredentials
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
		if errors.Is(err, errs.ErrFailedToGenToken) {
			return &models.LoginResponse{}, errs.ErrFailedToGenToken
		}
		return &models.LoginResponse{}, fmt.Errorf("problem generating token: %w", err)
	}

	refreshToken, err := u.app.JwtAUth.GenerateRefreshToken()
	if err != nil {
		if errors.Is(err, errs.ErrFailedToGenRefreshToken) {
			return &models.LoginResponse{}, errs.ErrFailedToGenRefreshToken
		}
		return &models.LoginResponse{}, fmt.Errorf("problem generating refresh token: %w", err)
	}

	if err := u.repo.StoreRefreshToken(ctx, user.ID, refreshToken, time.Now().Add(u.app.AppConfig.AuthConfig.RefreshTokenExp)); err != nil {
		if errors.Is(err, errs.ErrFailedToStoreToken) {
			return &models.LoginResponse{}, errs.ErrFailedToStoreToken
		}
		return &models.LoginResponse{}, fmt.Errorf("problem storing refresh token: %w", err)
	}

	return &models.LoginResponse{ID: user.ID, Username: user.Username, Token: token, RefreshToken: refreshToken}, nil
}

func (u *UserService) UserProfile(ctx context.Context, username string) (*models.UserResponse, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryDefaultContext)
	defer cancel()

	user, err := u.repo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return &models.UserResponse{}, errs.ErrUserNotFound
		}
		return &models.UserResponse{}, fmt.Errorf("failed to retrieve user: %w", err)
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

	ctx, cancel := context.WithTimeout(ctx, QueryDefaultContext)
	defer cancel()

	userId, err := u.repo.ValidateRefreshToken(ctx, refreshToken)
	if err != nil || userId == "" {
		if errors.Is(err, errs.ErrUserNotFound) {
			return "", fmt.Errorf("user not found")
		}
		return "", fmt.Errorf("invalid user details: %w", err)
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
		return "", errs.ErrFailedToGenToken
	}

	return newToken, nil
}

func (u *UserService) UpdateProfile(ctx context.Context, req *models.User, id string) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryDefaultContext)
	defer cancel()

	if req.Username == "" || len(req.Username) < 4 || req.Email == "" || req.FullName == "" || len(req.FullName) < 6 || req.Location == "" || len(req.Location) < 6 {
		return "", errs.ErrInvalidUserDetails
	}

	us := &models.User{
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		Location: req.Location,
	}

	if err := u.repo.UpdateUser(ctx, us, id); err != nil {
		return "", errs.ErrFailedToUpdateUser
	}

	return "user updated successfully", nil
}

func (u *UserService) ChangePassword(ctx context.Context, req *models.PasswordUpdateRequest, id string) (string, error) {

	ctx, cancel := context.WithTimeout(ctx, QueryDefaultContext)
	defer cancel()

	if req.OldPassword == "" || len(req.OldPassword) < 8 || req.ConfirmPassword == "" || len(req.ConfirmPassword) < 8 || req.NewPassword == "" || len(req.NewPassword) < 8 {
		return "", errs.ErrInvalidPassword
	}

	existingUser, err := u.repo.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			return "", errs.ErrUserNotFound
		}
		return "", fmt.Errorf("failed to retrieve user: %w", err)
	}

	if err := utils.CompareHashedPassword(existingUser.Password, req.OldPassword); err != nil {
		return "", errs.ErrInvalidPassword
	}

	if req.NewPassword != req.ConfirmPassword {
		return "", errs.ErrPasswordsDoNotMatch
	}

	if req.OldPassword == req.NewPassword {
		return "", errs.ErrPasswordCannotBeSame
	}

	hashedPassword, err := utils.HashedPassword(req.NewPassword)
	if err != nil {
		return "", errs.ErrFailedToHashPassword
	}

	if err := u.repo.ChangePassword(ctx, hashedPassword, id); err != nil {
		return "", errs.ErrFailedToChangePassword
	}

	return "password changed successfully", nil
}
