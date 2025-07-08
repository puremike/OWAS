package mock_services

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/puremike/online_auction_api/internal/config"
	"github.com/puremike/online_auction_api/internal/errs"
	"github.com/puremike/online_auction_api/internal/models"
	"github.com/puremike/online_auction_api/internal/services"
	"github.com/puremike/online_auction_api/internal/store/mock_store"
	"github.com/puremike/online_auction_api/internal/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUser_Success(t *testing.T) {
	// t.Helper()
	assert := assert.New(t)
	require := require.New(t)

	mockRepo := new(mock_store.MockUserStore)
	app := &config.Application{}
	hashedPassword, err := utils.HashedPassword("@SecurePassword123")
	require.NoError(err, "Expected no error when hashing password")

	userService := services.NewUserService(mockRepo, app)

	expectedUser := &models.User{
		ID:        "test-id-1",
		Username:  "Testuser",
		Email:     "test@example.com",
		FullName:  "Full Name Of User",
		Location:  "Some Location",
		CreatedAt: time.Now().Truncate(time.Second),
	}

	mockRepo.
		On("CreateUser", mock.Anything, mock.MatchedBy(func(user *models.User) bool {
			return user.Username == "testuser" &&
				user.Email == "test@example.com" &&
				strings.HasPrefix(user.Password, "$2a$")
		})).
		// mock.AnythingOfType()
		Return(expectedUser, nil).Once()

	inputUser := &models.User{
		Username: strings.ToLower("TestUser"),
		Email:    strings.ToLower("test@Example.com"),
		Password: hashedPassword,
		FullName: "Full Name Of User",
		Location: "Some Location",
	}

	res, err := userService.CreateUser(context.Background(), inputUser)
	assert.NoError(err, "Expected no error when creating user")
	assert.NotNil(res, "Expected a user response")
	assert.Equal(expectedUser.ID, res.ID, "Expected user ID to match")
	assert.Equal(expectedUser.Username, res.Username, "Expected username to be formatted correctly")
	assert.Equal(expectedUser.Email, res.Email, "Expected email to match")
	assert.Equal(expectedUser.FullName, res.FullName, "Expected full name to")
	assert.WithinDuration(expectedUser.CreatedAt, res.CreatedAt, time.Second, "Expected created at time to match")

	mockRepo.AssertExpectations(t)
}

func TestCreateUser_InvalidDetails(t *testing.T) {
	mockRepo := new(mock_store.MockUserStore)
	app := &config.Application{}
	userService := services.NewUserService(mockRepo, app)

	tests := []struct {
		name          string
		user          *models.User
		expectedError error
	}{
		{
			name: "unHashed password",
			user: &models.User{
				Username: "usr", Password: "unhashedPasword", Email: "usr@email.com",
			}, expectedError: errs.ErrInvalidUserDetails,
		},
		{
			name: "empty email",
			user: &models.User{
				Username: "usr", Password: "@SecurePassword123", Email: "",
			},
			expectedError: errs.ErrInvalidUserDetails,
		},
		{
			name: "short password",
			user: &models.User{
				Username: "usr", Password: "pass", Email: "usr@email.com",
			},
			expectedError: errs.ErrInvalidUserDetails,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			res, err := userService.CreateUser(context.Background(), tt.user)
			assert.ErrorIs(err, tt.expectedError)
			assert.Nil(res, "Expected nil response for invalid input")

			mockRepo.AssertNotCalled(t, "CreateUser", mock.Anything, mock.Anything)
		})
	}
}
