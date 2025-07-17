package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/puremike/online_auction_api/contexts"
	"github.com/puremike/online_auction_api/internal/config"
	"github.com/puremike/online_auction_api/internal/errs"
	"github.com/puremike/online_auction_api/internal/ratelimiters"
)

type Middleware struct {
	app *config.Application
}

func NewMiddleware(app *config.Application) *Middleware {
	return &Middleware{
		app: app,
	}
}

func (m *Middleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		var tokenString string

		isWebSocketUpgrade := strings.ToLower(c.GetHeader("Upgrade")) == "websocket"
		if isWebSocketUpgrade {
			queryToken := c.Query("token")
			if queryToken != "" {
				tokenString = queryToken
			}
		}

		cookieToken, err := c.Cookie("jwt")
		if err == nil && cookieToken != "" {
			tokenString = cookieToken
		}

		if tokenString == "" {
			authHeader := c.GetHeader("Authorization")

			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
				c.Abort()
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is deformed"})
				c.Abort()
				return
			}

			tokenString = strings.TrimSpace(parts[1])
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication token missing"})
			c.Abort()
			return
		}

		jwtToken, err := m.app.JwtAUth.ValidateToken(tokenString)
		if jwtToken == nil || err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)
		userId, ok := claims["sub"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid sub type"})
			c.Abort()
			return
		}

		if !m.app.AppConfig.RedisCacheConf.Enabled {
			user, err := m.app.Store.Users.GetUserById(c.Request.Context(), userId)

			if err != nil {
				handleUserFetchError(c, err)
				return
			}

			c.Set("user", user)
			c.Next()
			return
		}

		// Redis is enabled
		m.app.Logger.Infow("cache hit", "key", userId)

		user, err := m.app.RedisCache.Users.Get(c.Request.Context(), userId)
		if err != nil {
			m.app.Logger.Errorw("failed to get user from cache", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user from cache"})
			c.Abort()
			return
		}

		if user != nil {
			m.app.Logger.Infow("cache hit", "key", userId)
			c.Set("user", user)
			c.Set("userId", user.ID)
			c.Next()
			return
		}

		// Fetch from database if not found in cache
		m.app.Logger.Infow("cache miss", "key", "id", userId)
		user, err = m.app.Store.Users.GetUserById(c.Request.Context(), userId)
		if err != nil {
			handleUserFetchError(c, err)
			return
		}

		// set user in cache
		m.app.Logger.Infow("setting user in cache", "userId", userId)
		if err := m.app.RedisCache.Users.Set(c.Request.Context(), user); err != nil {
			m.app.Logger.Errorw("failed to set user in cache", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set user in cache"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("userId", user.ID)
		c.Next()
	}
}

func handleUserFetchError(c *gin.Context, err error) {
	if errors.Is(err, errs.ErrUserNotFound) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user"})
	}
	c.Abort()
}

// Grant administrator access to the user
func AuthorizeRoles(allowedRole bool) gin.HandlerFunc {
	return func(c *gin.Context) {

		user, err := contexts.GetUserFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		if user.IsAdmin && allowedRole {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient role"})
	}
}

func (m *Middleware) AuctionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		auctionId := c.Param("auctionID")

		if !m.app.AppConfig.RedisCacheConf.Enabled {
			auction, err := m.app.Store.Auctions.GetAuctionById(c.Request.Context(), auctionId)
			if err != nil {
				if errors.Is(err, errs.ErrAuctionNotFound) {
					c.JSON(http.StatusNotFound, gin.H{"error": "auction not found"})
					c.Abort()
					return
				}
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve auction"})
				c.Abort()
				return
			}

			c.Set("auction", auction)
			c.Next()
			return
		}

		// Redis is enabled
		m.app.Logger.Infow("cache hit", "key", auctionId)

		auction, err := m.app.RedisCache.Auctions.Get(c.Request.Context(), auctionId)

		if err != nil {
			m.app.Logger.Errorw("failed to get auction from cache", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve auction from cache"})
			c.Abort()
			return
		}

		if auction != nil {
			m.app.Logger.Infow("cache hit", "key", auctionId)
			c.Set("auction", auction)
			c.Next()
			return
		}

		// Fetch from database if not found in cache
		auction, err = m.app.Store.Auctions.GetAuctionById(c.Request.Context(), auctionId)
		if err != nil {
			if errors.Is(err, errs.ErrAuctionNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "auction not found"})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve auction from database"})
			c.Abort()
			return
		}

		// set auction in cache
		m.app.Logger.Infow("setting auction in cache", "auctionId", auctionId)

		if err = m.app.RedisCache.Auctions.Set(c.Request.Context(), auction); err != nil {
			m.app.Logger.Errorw("failed to set auction in cache", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set auction in cache"})
			c.Abort()
			return
		}

		c.Set("auction", auction)
		c.Next()
	}
}

func (m *Middleware) RateLimiterMiddleware(limiter ratelimiters.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allowed() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}
