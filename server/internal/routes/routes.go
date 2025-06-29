package routes

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/puremike/online_auction_api/internal/config"
	"github.com/puremike/online_auction_api/internal/handlers"
	"github.com/puremike/online_auction_api/internal/imagesuploader"
	"github.com/puremike/online_auction_api/internal/middlewares"
	"github.com/puremike/online_auction_api/internal/services"
	"github.com/puremike/online_auction_api/internal/ws"
	"github.com/puremike/online_auction_api/pkg"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routes(app *config.Application) http.Handler {

	g := gin.Default()

	g.Use(func(c *gin.Context) {
		log.Printf("DEBUG: Incoming Request - Method: %s, Path: %s, From: %s",
			c.Request.Method, c.Request.URL.Path, c.ClientIP())
		c.Next()
	})

	g.Use(cors.New(cors.Config{
		AllowOrigins:     []string{pkg.GetEnvString("FRONTEND_URL", "http://localhost:3000")},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	userService := services.NewUserService(app.Store.Users, app)
	userHandler := handlers.NewUserHandler(userService, app)

	auctionService := services.NewAuctionService(app.Store.Auctions, app.Store.Bids, app.Store.Notifications, app.WsHub.AuctionUpdates, app.WsHub.NotificationUpdates)
	auctionHandler := handlers.NewAuctionHandler(auctionService, app)

	middleware := middlewares.NewMiddleware(app)

	csService := services.NewCSService(app.Store.CS)
	csHandler := handlers.NewCSHandler(csService)

	wsHandler := ws.NewWSHandler(app.WsHub)

	paymentService := services.NewPaymentService(app.Stripe, app.Store.Payments)
	webHookHandler := handlers.NewWebHookHander(paymentService, app.Store.Auctions)

	imageService := imagesuploader.NewImageService(app.AppConfig.S3Bucket)
	imageHandler := handlers.NewImageHandler(imageService)

	api := g.Group("/api/v1")
	api.Use(middleware.RateLimiterMiddleware(app.GeneralRateLimiter))
	{
		api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
		api.GET("/health", middleware.RateLimiterMiddleware(app.HeavyOpsRateLimiter), handlers.Health)
		api.GET("/checking", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "checking",
			})
		})
		api.POST("/webhook/stripe", webHookHandler.StripeWebHookHandler)
	}

	user := api.Group("/")
	{
		user.POST("/signup", userHandler.RegisterUser)
		user.POST("/login", middleware.RateLimiterMiddleware(app.SensitiveRateLimiter), userHandler.Login)
		user.POST("/refresh", middleware.RateLimiterMiddleware(app.SensitiveRateLimiter), userHandler.RefreshToken)
	}

	authGroup := api.Group("/")
	authGroup.Use(middleware.AuthMiddleware())
	{
		authGroup.POST("/logout", userHandler.Logout)
		authGroup.GET("/me", userHandler.MeProfile)
		authGroup.GET("/:username", userHandler.UserProfile)
		authGroup.PUT("/:username/update-profile", userHandler.UpdateProfile)
		authGroup.PUT("/:username/change-password", userHandler.ChangePassword)

		authGroup.GET("/admin/users", middlewares.AuthorizeRoles(true), userHandler.AdminGetUsers)
		authGroup.GET("/auctions", auctionHandler.GetAuctions)
		authGroup.DELETE("/admin/auctions/:auctionID", middleware.AuctionMiddleware(), middlewares.AuthorizeRoles(true), auctionHandler.AdminDeleteAuction)

		authGroup.POST("/auctions", auctionHandler.CreateAuction)
		authGroup.GET("/auctions/:auctionID", middleware.AuctionMiddleware(), auctionHandler.GetAuctionById)
		authGroup.PUT("/auctions/:auctionID", middleware.AuctionMiddleware(), auctionHandler.UpdateAuction)
		authGroup.DELETE("/auctions/:auctionID", middleware.AuctionMiddleware(), auctionHandler.DeleteAuction)

		authGroup.POST("/auctions/:auctionID/bids", middleware.AuctionMiddleware(), auctionHandler.PlaceBids)
		authGroup.POST("/auctions/:auctionID/close", middleware.AuctionMiddleware(), auctionHandler.CloseAuction)

		authGroup.POST("/contact-support", csHandler.ContactSupport)

		authGroup.GET("/ws", wsHandler.ServeWs)

		authGroup.POST("/auctions/:auctionID/create-checkout-session", middleware.AuctionMiddleware(), webHookHandler.CreateCheckoutSessionHandler)

		authGroup.POST("/auctions/image_upload", imageHandler.UploadImage)
	}

	return g
}
