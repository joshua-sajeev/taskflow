package user_handler

import (
	"net/http"
	"strconv"
	"strings"

	"taskflow/internal/common"
	"taskflow/pkg"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware protects routes and requires a valid JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := []byte(pkg.GetEnv("JWT_SECRET", "secret-key"))

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, common.ErrorResponse{
				Message: "authorization header required",
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, common.ErrorResponse{
				Message: "invalid authorization header format",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		claims, err := pkg.ValidateToken(tokenString, secretKey)
		if err != nil {
			c.JSON(http.StatusUnauthorized, common.ErrorResponse{
				Message: "invalid or expired token",
			})
			c.Abort()
			return
		}

		// Handle numeric or string user_id from claims
		claimsMap := *claims
		var userID int
		switch v := claimsMap["user_id"].(type) {
		case float64:
			userID = int(v)
		case int:
			userID = v
		case string:
			var err error
			userID, err = strconv.Atoi(v)
			if err != nil {
				c.JSON(http.StatusUnauthorized, common.ErrorResponse{
					Message: "invalid user ID in token",
				})
				c.Abort()
				return
			}
		default:
			c.JSON(http.StatusUnauthorized, common.ErrorResponse{
				Message: "invalid token claims",
			})
			c.Abort()
			return
		}

		// Set userID in context for handlers
		c.Set("userID", userID)
		c.Next()
	}
}

// OptionalAuthMiddleware allows routes to continue even if JWT is missing/invalid
func OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := []byte(pkg.GetEnv("JWT_SECRET", "secret-key"))

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]
		claims, err := pkg.ValidateToken(tokenString, secretKey)
		if err != nil {
			c.Next()
			return
		}

		claimsMap := *claims
		switch v := claimsMap["user_id"].(type) {
		case float64:
			c.Set("userID", int(v))
		case int:
			c.Set("userID", v)
		case string:
			if userID, err := strconv.Atoi(v); err == nil {
				c.Set("userID", userID)
			}
		}

		c.Next()
	}
}
