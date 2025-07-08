package mock_store

import (
	"context"
	"time"

	"github.com/puremike/online_auction_api/internal/models"
	"github.com/puremike/online_auction_api/internal/store"
	"github.com/stretchr/testify/mock"
)

var _ store.UserRepository = (*MockUserStore)(nil)

type MockUserStore struct {
	mock.Mock
}

func (u *MockUserStore) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	ret := u.Called(ctx, user)

	// VERBOSE METHOD
	// var r0 *models.User
	// if rf, ok := ret.Get(0).((func(user *models.User) *models.User)); ok {
	// 	r0 = rf(ctx, user)
	// } else {
	// 	if ret.Get(0) != nil {
	// 		r0 = ret.Get(0).(*models.User)
	// 	}

	// }

	// var r1 error
	// if rf, ok := ret.Get(1).((func(user *models.User) *models.User)); ok {
	// 	r1 = rf(ctx, user)
	// } else {
	// 	r1 = ret.Error(1)

	// }

	// return r0, r1

	return ret.Get(0).(*models.User), ret.Error(1)
}

// return ret.Get(0).(*models.User), ret.Error(1)

func (u *MockUserStore) GetUserById(ctx context.Context, id string) (*models.User, error) {
	ret := u.Called(ctx, id)
	return ret.Get(0).(*models.User), ret.Error(1)
}
func (u *MockUserStore) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	ret := u.Called(ctx, email)
	return ret.Get(0).(*models.User), ret.Error(1)
}
func (u *MockUserStore) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	ret := u.Called(ctx, username)
	return ret.Get(0).(*models.User), ret.Error(1)
}
func (u *MockUserStore) StoreRefreshToken(ctx context.Context, userID, refreshToken string, expires_at time.Time) error {
	ret := u.Called(ctx, userID, refreshToken, expires_at)
	return ret.Error(0)
}
func (u *MockUserStore) UpdateUser(ctx context.Context, user *models.User, id string) error {
	ret := u.Called(ctx, user, id)
	return ret.Error(0)
}
func (u *MockUserStore) ValidateRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	ret := u.Called(ctx, refreshToken)
	return ret.String(0), ret.Error(1)
}
func (u *MockUserStore) ChangePassword(ctx context.Context, pass, id string) error {
	ret := u.Called(ctx, pass, id)
	return ret.Error(0)
}
func (u *MockUserStore) GetUsers(ctx context.Context) (*[]models.User, error) {
	ret := u.Called(ctx)
	return ret.Get(0).(*[]models.User), ret.Error(1)
}
func (u *MockUserStore) DeleteUser(ctx context.Context, id string) error {
	ret := u.Called(ctx, id)
	return ret.Error(0)
}
