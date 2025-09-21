package user_handler

import "github.com/gin-gonic/gin"

type UserHandlerInterface interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	UpdatePassword(c *gin.Context)
	DeleteUser(c *gin.Context)
}
