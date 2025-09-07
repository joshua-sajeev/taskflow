package task

import "github.com/gin-gonic/gin"

type TaskHandlerInterface interface {
	// CreateTask handles POST /api/tasks
	CreateTask(c *gin.Context)

	// GetTask handles GET /api/tasks/:id
	GetTask(c *gin.Context)

	// ListTasks handles GET /api/tasks
	ListTasks(c *gin.Context)

	// UpdateStatus handles PATCH /api/tasks/:id/status
	UpdateStatus(c *gin.Context)

	// Delete handles DELETE /api/tasks/:id
	Delete(c *gin.Context)
}
