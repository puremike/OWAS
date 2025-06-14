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
	auctionService := services.NewAuctionService(app.Store.Auctions)
	auctionHandler := handlers.NewAuctionHandler(auctionService, app)
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

		authGroup.GET("/admin/users", middlewares.AuthorizeRoles(true), userHandler.AdminGetUsers)
		authGroup.GET("/admin/auctions", middlewares.AuthorizeRoles(true), auctionHandler.AdminGetAuctions)
		authGroup.DELETE("/admin/auctions/:auctionID", middleware.AuctionMiddleware(), middlewares.AuthorizeRoles(true), auctionHandler.AdminDeleteAuction)

		authGroup.POST("/auctions", auctionHandler.CreateAuction)
		authGroup.GET("/auctions/:auctionID", middleware.AuctionMiddleware(), auctionHandler.GetAuctionById)
		authGroup.PUT("/auctions/:auctionID", middleware.AuctionMiddleware(), auctionHandler.UpdateAuction)
		authGroup.DELETE("/auctions/:auctionID", middleware.AuctionMiddleware(), auctionHandler.DeleteAuction)
	}

	return g
}
