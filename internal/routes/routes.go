package routes

import (
	"taskflow/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.GET("/", handlers.HomeHandler)
	r.GET("/ping", handlers.PingHandler)
	return r
}
