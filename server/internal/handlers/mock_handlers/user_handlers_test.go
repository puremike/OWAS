package mock_handlers

import (
	"log"
	"testing"

	"github.com/puremike/online_auction_api/internal/config"
	"github.com/puremike/online_auction_api/internal/handlers"
	"github.com/puremike/online_auction_api/internal/services/mock_services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUserSuccess(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	app := new(config.Application)
	mockUserService := new(mock_services.MockUserService)
	userHandler := handlers.NewUserHandler(mockUserService, app)

	log.Print(assert, require, userHandler)
}
