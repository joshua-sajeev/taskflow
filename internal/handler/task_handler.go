package handler

import (
	"net/http"
	"strconv"

	"taskflow/internal/common"
	"taskflow/internal/domain/task"
	"taskflow/internal/service"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(s *service.TaskService) *TaskHandler {
	return &TaskHandler{service: s}
}

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task with title and description
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body task.Task true "Task to create"
// @Success 201 {object} task.Task
// @Failure 400 {object} common.ErrorResponse
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var t task.Task
	if err := c.ShouldBindJSON(&t); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}
	if err := h.service.CreateTask(&t); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, t)
}

// GetTask godoc
// @Summary Get a task by ID
// @Description Get details of a specific task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} task.Task
// @Failure 404 {object} common.ErrorResponse
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	t, err := h.service.GetTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, common.ErrorResponse{Error: "Task not found"})
		return
	}
	c.JSON(http.StatusOK, t)
}

// ListTasks godoc
// @Summary List all tasks
// @Description Get a list of all tasks
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {array} task.Task
// @Failure 500 {object} common.ErrorResponse
// @Router /tasks [get]
func (h *TaskHandler) ListTasks(c *gin.Context) {
	tasks, err := h.service.ListTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// UpdateStatus godoc
// @Summary Update task status
// @Description Update the status field of a task
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param status body map[string]string true "Status payload"
// @Success 200 {object} map[string]string
// @Failure 400 {object} common.ErrorResponse
// @Router /tasks/{id}/status [patch]
func (h *TaskHandler) UpdateStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var body struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}
	if err := h.service.UpdateStatus(id, body.Status); err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "status updated"})
}

// Delete godoc
// @Summary Delete a task
// @Description Delete a task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /tasks/{id} [delete]
func (h *TaskHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.ErrorResponse{Error: "Invalid ID"})
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, common.ErrorResponse{Error: "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, common.ErrorResponse{Error: "Couldn't delete task"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
