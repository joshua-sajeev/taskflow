package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Server index.tmpl file
func HomeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "Home Page",
	})
}

// Checks Health
func PingHandler(c *gin.Context) {
	now := time.Now()
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"time":    now,
	})
}
