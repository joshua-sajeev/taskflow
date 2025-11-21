package auth

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"taskflow/internal/common"
	"taskflow/internal/repository/gorm/gorm_user"
	"taskflow/pkg/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserAuth struct {
	secretKey []byte
	userRepo  gorm_user.UserRepositoryInterface
}

func NewUserAuth(secret string, userRepo gorm_user.UserRepositoryInterface) *UserAuth {
	return &UserAuth{secretKey: []byte(secret), userRepo: userRepo}
}

var _ UserAuthInterface = (*UserAuth)(nil)

func (ua *UserAuth) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := ua.secretKey

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

		claims, err := jwt.ValidateToken(tokenString, secretKey)
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

		if ua.userRepo != nil {
			_, err := ua.userRepo.GetByID(userID)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					c.JSON(http.StatusUnauthorized, common.ErrorResponse{
						Message: "user account not found",
					})
				} else {
					c.JSON(http.StatusInternalServerError, common.ErrorResponse{
						Message: "authentication failed",
					})
				}
				c.Abort()
				return
			}
		}

		c.Set("userID", userID)
		c.Next()
	}
}

// OptionalAuthMiddleware allows routes to continue even if JWT is missing/invalid
func (ua *UserAuth) OptionalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := ua.secretKey

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
		claims, err := jwt.ValidateToken(tokenString, secretKey)
		if err != nil {
			c.Next()
			return
		}

		claimsMap := *claims
		var userID int
		var userIDValid bool

		switch v := claimsMap["user_id"].(type) {
		case float64:
			userID = int(v)
			userIDValid = true
		case int:
			userID = v
			userIDValid = true
		case string:
			if id, err := strconv.Atoi(v); err == nil {
				userID = id
				userIDValid = true
			}
		}

		if userIDValid && ua.userRepo != nil {
			_, err := ua.userRepo.GetByID(userID)
			if err != nil {
				c.Next()
				return
			}
			c.Set("userID", userID)
		}

		c.Next()
	}
}

func SetUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, "userID", userID)
}
