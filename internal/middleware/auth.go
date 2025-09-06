package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/immunesh/automated-social-media-scheduler/internal/auth"
)

// AuthRequired is middleware that validates JWT tokens
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		
		// Check if Authorization header exists
		if authHeader == "" {
			// For web pages, check for cookie
			cookie, err := c.Cookie("auth_token")
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
				c.Abort()
				return
			}
			authHeader = "Bearer " + cookie
		}

		// Extract token from header
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Validate token
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)

		c.Next()
	}
}

// OptionalAuth is middleware that optionally validates JWT tokens
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			cookie, err := c.Cookie("auth_token")
			if err != nil {
				c.Next()
				return
			}
			authHeader = "Bearer " + cookie
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := auth.ValidateToken(tokenString)
		if err == nil {
			c.Set("user_id", claims.UserID)
			c.Set("user_email", claims.Email)
		}

		c.Next()
	}
}