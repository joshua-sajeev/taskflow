package bootstrap

import (
	"github.com/gin-gonic/gin"
	"taskflow/internal/handlers"
	"taskflow/internal/routes"
)

func SetupRouter(jobHandler *handlers.JobHandler, dashboardHandler *handlers.DashboardHandler) *gin.Engine {
	router := routes.SetupRoutes()
	routes.RegisterJobRoutes(router, jobHandler)
	routes.RegisterDashboardRoutes(router, dashboardHandler)
	return router
}
