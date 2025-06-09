package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/puremike/online_auction_api/internal/config"
	"github.com/puremike/online_auction_api/internal/handlers"
	"github.com/puremike/online_auction_api/internal/middlewares"
	"github.com/puremike/online_auction_api/internal/services"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routes(app *config.Application) http.Handler {

	g := gin.Default()

	userService := services.NewUserService(app.Store.Users, app)
	userHandler := handlers.NewUserHandler(userService, app)
	middleware := middlewares.NewMiddleware(app)

	api := g.Group("/api/v1")
	{
		api.GET("swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		api.GET("/health", handlers.Health)
	}

	user := api.Group("/")
	{
		user.POST("/signup", userHandler.RegisterUser)
		user.POST("/login", userHandler.Login)
		user.POST("/refresh", userHandler.RefreshToken)
	}

	authGroup := api.Group("/")
	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.POST("/logout", userHandler.Logout)
		authGroup.GET("/:username", userHandler.UserProfile)
		authGroup.PUT("/:username/update-profile", userHandler.UpdateProfile)
		authGroup.PUT("/:username/change-password", userHandler.ChangePassword)
	}

	return g
}
