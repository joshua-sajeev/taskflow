package handlers

import (
	"net/http"
	"taskflow/internal/domain/entities"
	repositories "taskflow/internal/domain/repository"

	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	repo repositories.JobRepository
}

func NewJobHandler(repo repositories.JobRepository) *JobHandler {
	return &JobHandler{
		repo: repo,
	}
}

func (j *JobHandler) CreateJob(c *gin.Context) {
	var job entities.Job

	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	if err := j.repo.Create(&job); err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create job",
		})
		return
	}
	c.JSON(http.StatusCreated, job)
}
