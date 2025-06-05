package services

import (
	"context"
	"fmt"

	"github.com/puremike/online_auction_api/internal/models"
	"github.com/puremike/online_auction_api/internal/store"
	"github.com/puremike/online_auction_api/internal/utils"
)

type UserService struct {
	repo store.UserRepository
}

func NewUserService(repo store.UserRepository) *UserService {
	return &UserService{
		repo: repo,
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
