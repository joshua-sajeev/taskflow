package bootstrap

import (
	"github.com/gin-gonic/gin"
	"taskflow/internal/handlers"
	"taskflow/internal/routes"
)

func SetupRouter(jobHandler *handlers.JobHandler) *gin.Engine {
	router := routes.SetupRouter()
	routes.RegisterJobRoutes(router, jobHandler)
	return router
}
