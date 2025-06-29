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

type Middeware struct {
	app *config.Application
}

func NewMiddleware(app *config.Application) *Middeware {
	return &Middeware{
		app: app,
	}
}

func (m *Middeware) AuthMiddleware() gin.HandlerFunc {
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

		user, err := m.app.Store.Users.GetUserById(c.Request.Context(), userId)
		if err != nil {
			if errors.Is(err, errs.ErrUserNotFound) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve user"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("userId", user.ID)
		c.Next()
	}
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

func (m *Middeware) AuctionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		auctionId := c.Param("auctionID")

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
	}
}

func (m *Middeware) RateLimiterMiddleware(limiter ratelimiters.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allowed() {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}
