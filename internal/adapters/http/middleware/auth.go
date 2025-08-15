package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"hexagonal-go/internal/utils"
)

// AuthMiddleware validates JWT tokens from the Authorization header.
// It expects the header to be in the format: "Bearer <token>".
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header missing"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		userID, err := utils.ValidateJWT(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// store userID in context for downstream handlers
		c.Set("userID", userID)
		c.Next()
	}
}
