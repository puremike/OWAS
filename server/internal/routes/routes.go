package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/puremike/online_auction_api/internal/config"
	"github.com/puremike/online_auction_api/internal/handlers"
	"github.com/puremike/online_auction_api/internal/services"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routes(app *config.Application) http.Handler {

	g := gin.Default()

	userService := services.NewUserService(app.Store.Users)
	userHandler := handlers.NewUserHandler(userService)

	api := g.Group("/api/v1")
	{
		api.GET("swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		api.GET("/health", handlers.Health)
	}

	user := api.Group("/")
	{
		user.POST("/signup", userHandler.RegisterUser)
	}

	return g
}
