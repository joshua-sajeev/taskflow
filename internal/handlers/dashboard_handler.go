package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	repositories "taskflow/internal/domain/repository"
	"taskflow/internal/jobqueue"
	"time"
)

type DashboardHandler struct {
	repo       repositories.JobRepository
	workerPool *jobqueue.WorkerPool
}

func NewDashboardHandler(repo repositories.JobRepository, workerPool *jobqueue.WorkerPool) *DashboardHandler {
	return &DashboardHandler{
		repo:       repo,
		workerPool: workerPool,
	}
}

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// Render the HTML dashboard page
func (d *DashboardHandler) DisplayStats(c *gin.Context) {
	pendingJobs, completedJobs, err := d.repo.CountJobsByStatus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch job stats"})
		return
	}

	c.HTML(http.StatusOK, "dashboard.tmpl", gin.H{
		"title":         "Dashboard",
		"workers":       d.workerPool.ActiveWorkers(),
		"queueSize":     d.workerPool.JobQueue.Size(),
		"pendingJobs":   pendingJobs,
		"completedJobs": completedJobs,
	})
}

// Stream stats over WebSocket
func (d *DashboardHandler) StreamStats(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		pendingJobs, completedJobs, err := d.repo.CountJobsByStatus()
		if err != nil {
			return
		}

		stats := gin.H{
			"workers":       d.workerPool.ActiveWorkers(),
			"queueSize":     d.workerPool.JobQueue.Size(),
			"pendingJobs":   pendingJobs,
			"completedJobs": completedJobs,
		}

		if err := conn.WriteJSON(stats); err != nil {
			break
		}

		time.Sleep(1 * time.Second)
	}
}
