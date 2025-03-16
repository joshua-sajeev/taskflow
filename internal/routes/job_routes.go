package routes

import (
	"taskflow/internal/handlers"

	"github.com/gin-gonic/gin"
)

// RegisterJobRoutes registers job-related routes
func RegisterJobRoutes(router *gin.Engine, jobHandler *handlers.JobHandler) {

	jobRoutes := router.Group("/jobs")
	{
		jobRoutes.POST("/", jobHandler.CreateJob)
	}

}
