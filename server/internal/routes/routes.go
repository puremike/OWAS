package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/puremike/online_auction_api/internal/config"
	"github.com/puremike/online_auction_api/internal/handlers"
	"github.com/puremike/online_auction_api/internal/middlewares"
	"github.com/puremike/online_auction_api/internal/services"
	"github.com/puremike/online_auction_api/internal/ws"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routes(app *config.Application) http.Handler {

	g := gin.Default()

	userService := services.NewUserService(app.Store.Users, app)
	userHandler := handlers.NewUserHandler(userService, app)

	auctionService := services.NewAuctionService(app.Store.Auctions, app.Store.Bids, app.Store.Notifications, app.WsHub.AuctionUpdates, app.WsHub.NotificationUpdates)
	auctionHandler := handlers.NewAuctionHandler(auctionService, app)

	middleware := middlewares.NewMiddleware(app)

	csService := services.NewCSService(app.Store.CS)
	csHandler := handlers.NewCSHandler(csService)

	wsHandler := ws.NewWSHandler(app.WsHub)
	g.POST("/contact-support", csHandler.ContactSupport)

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
		authGroup.GET("/auctions", auctionHandler.AdminGetAuctions)
		authGroup.DELETE("/admin/auctions/:auctionID", middleware.AuctionMiddleware(), middlewares.AuthorizeRoles(true), auctionHandler.AdminDeleteAuction)

		authGroup.POST("/auctions", auctionHandler.CreateAuction)
		authGroup.GET("/auctions/:auctionID", middleware.AuctionMiddleware(), auctionHandler.GetAuctionById)
		authGroup.PUT("/auctions/:auctionID", middleware.AuctionMiddleware(), auctionHandler.UpdateAuction)
		authGroup.DELETE("/auctions/:auctionID", middleware.AuctionMiddleware(), auctionHandler.DeleteAuction)

		authGroup.POST("/auctions/:auctionID/bids", middleware.AuctionMiddleware(), auctionHandler.PlaceBids)
		authGroup.POST("/auctions/:auctionID/close", middleware.AuctionMiddleware(), auctionHandler.CloseAuction)

		authGroup.POST("/contact-support", csHandler.ContactSupport)

		authGroup.GET("/ws", wsHandler.ServeWs)
	}

	return g
}
