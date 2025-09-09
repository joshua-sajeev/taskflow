package handler

import (
	"errors"
	"net/http"
	"strconv"

	"taskflow/internal/common"
	"taskflow/internal/domain/task"
	"taskflow/internal/dto"
	"taskflow/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskHandler struct {
	service service.TaskServiceInterface
}

func NewTaskHandler(s service.TaskServiceInterface) *TaskHandler {
	return &TaskHandler{service: s}
}

var _ task.TaskHandlerInterface = (*TaskHandler)(nil)

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task with title and description
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body dto.CreateTaskRequest true "Task to create"
// @Success 201 {object} task.Task
// @Failure 400 {object} common.ErrorResponse
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: err.Error()})
		return
	}
	if err := h.service.CreateTask(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, req)
}

// GetTask godoc
// @Summary Get a task by ID
// @Description Returns details of a specific task by ID
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID (>=1)" minimum(1) example(1)
// @Success 200 {object} dto.GetTaskResponse "Task retrieved successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid ID"
// @Failure 404 {object} common.ErrorResponse "Task not found"
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Invalid ID"})
		return
	}

	resp, err := h.service.GetTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.ErrorResponse{Message: "Task not found"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// Handler Layer
// ListTasks godoc
// @Summary List all tasks
// @Description Get a list of all tasks
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {object} dto.ListTasksResponse "List of tasks retrieved successfully"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /tasks [get]
func (h *TaskHandler) ListTasks(c *gin.Context) {
	res, err := h.service.ListTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// UpdateStatus godoc
// @Summary Update task status
// @Description Update the status field of a task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID" minimum(1) example(1)
// @Param request body dto.UpdateStatusRequest true "Status payload"
// @Success 200 {object} dto.UpdateStatusResponse "Status updated successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid input"
// @Failure 404 {object} common.ErrorResponse "Task not found"
// @Router /tasks/{id}/status [patch]
func (h *TaskHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "invalid task ID"})
		return
	}

	var req dto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: err.Error()})
		return
	}

	if err := h.service.UpdateStatus(id, req.Status); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, common.ErrorResponse{Message: "Task not found"})
			return
		}
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: err.Error()})
		return
	}

	resp := dto.UpdateStatusResponse{Message: "status updated"}
	c.JSON(http.StatusOK, resp)
}

// Delete godoc
// @Summary Delete a task
// @Description Delete a task by its ID
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID" minimum(1) example(1)
// @Success 200 {object} dto.DeleteTaskResponse "Task deleted successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid ID"
// @Failure 404 {object} common.ErrorResponse "Task not found"
// @Failure 500 {object} common.ErrorResponse "Internal server error"
// @Router /tasks/{id} [delete]
func (h *TaskHandler) Delete(c *gin.Context) {
	// Parse ID from path
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: "Invalid ID"})
		return
	}

	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, common.ErrorResponse{Message: "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, common.ErrorResponse{Message: "Couldn't delete task"})
		}
		return
	}

	resp := dto.DeleteTaskResponse{Message: "Task deleted successfully"}
	c.JSON(http.StatusOK, resp)
}
