package handlers

import (
	"log"
	"net/http"
	"taskflow/internal/domain/entities"
	repositories "taskflow/internal/domain/repository"
	"taskflow/internal/jobqueue"

	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	repo     repositories.JobRepository
	jobQueue jobqueue.JobQueue
}

func NewJobHandler(repo repositories.JobRepository, jobQueue jobqueue.JobQueue) *JobHandler {
	return &JobHandler{
		repo:     repo,
		jobQueue: jobQueue,
	}
}

func (j *JobHandler) CreateJob(c *gin.Context) {
	var request struct {
		Task string `json:"task"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	job := entities.NewJob(request.Task)

	// Insert job into the database
	if err := j.repo.Create(&job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
		return
	}

	// Enqueue job for processing
	j.jobQueue.Enqueue(job.ID)
	log.Println(job)
	c.JSON(http.StatusCreated, job)
}
