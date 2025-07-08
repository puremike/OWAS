package mock_services

import (
	"context"

	"github.com/puremike/online_auction_api/internal/models"
	"github.com/puremike/online_auction_api/internal/services"
	"github.com/stretchr/testify/mock"
)

var _ services.UserServiceInterface = (*MockUserService)(nil)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, user *models.User) (*models.UserResponse, error) {
	ret := m.Called(ctx, user)
	return ret.Get(0).(*models.UserResponse), ret.Error(1)
}

func (m *MockUserService) Login(ctx context.Context, req *models.LoginRequest) (*models.LoginResponse, error) {
	ret := m.Called(ctx, req)
	return ret.Get(0).(*models.LoginResponse), ret.Error(1)
}
func (m *MockUserService) UserProfile(ctx context.Context, username string) (*models.UserResponse, error) {
	ret := m.Called(ctx, username)
	return ret.Get(0).(*models.UserResponse), ret.Error(1)
}
func (m *MockUserService) MeProfile(ctx context.Context, userID string) (*models.UserResponse, error) {
	ret := m.Called(ctx, userID)
	return ret.Get(0).(*models.UserResponse), ret.Error(1)
}
func (m *MockUserService) Refresh(ctx context.Context, refreshToken string) (string, error) {
	ret := m.Called(ctx, refreshToken)
	return ret.String(0), ret.Error(1)
}
func (m *MockUserService) UpdateProfile(ctx context.Context, req *models.User, id string) (string, error) {
	ret := m.Called(ctx, req, id)
	return ret.String(0), ret.Error(1)
}
func (m *MockUserService) ChangePassword(ctx context.Context, req *models.PasswordUpdateRequest, id string) (string, error) {
	ret := m.Called(ctx, req, id)
	return ret.String(0), ret.Error(1)
}
func (m *MockUserService) GetUsers(ctx context.Context) (*[]models.UserResponse, error) {
	ret := m.Called(ctx)
	return ret.Get(0).(*[]models.UserResponse), ret.Error(1)
}
func (m *MockUserService) DeleteUser(ctx context.Context, id string) (string, error) {
	ret := m.Called(ctx, id)
	return ret.String(0), ret.Error(1)
}
