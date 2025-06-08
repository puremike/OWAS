package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/puremike/online_auction_api/contexts"
	"github.com/puremike/online_auction_api/internal/config"
	"github.com/puremike/online_auction_api/internal/store"
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

		tokenString, err := c.Cookie("jwt")

		if err != nil || tokenString == "" {

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
			if errors.Is(err, store.ErrUserNotFound) {
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
		user := contexts.GetUserFromContext(c)

		if user.IsAdmin && allowedRole {
			c.Next()
			return
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden: insufficient role"})
	}
}
