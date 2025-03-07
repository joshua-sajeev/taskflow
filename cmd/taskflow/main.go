package main

import (
	"taskflow/internal/db"
	"taskflow/internal/domain/entities"
	"taskflow/internal/domain/repositories"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize DB
	database, err := db.InitDB()
	if err != nil {
		logrus.Fatal("Error initializing database:", err)
	}

	// Initialize repositories
	jobRepo := repositories.NewJobRepo(database)

	job := &entities.Job{Task: "milk"} // Use pointer here
	if err := jobRepo.Create(job); err != nil {
		logrus.Error("Failed to create job:", err)
	} else {
		logrus.Info("Job created successfully:", job)
	}

	gin.SetMode(gin.ReleaseMode)
	// Setup router
	r := gin.New()
	// Set trusted proxies
	r.SetTrustedProxies(nil)

	// Add middleware manually
	r.Use(gin.Logger(), gin.Recovery())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Start server
	r.Run(":8080")
}
