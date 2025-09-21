package auth

import "github.com/gin-gonic/gin"

type UserAuthInterface interface {
	AuthMiddleware() gin.HandlerFunc
	OptionalAuthMiddleware() gin.HandlerFunc
}
