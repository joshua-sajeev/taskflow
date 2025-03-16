package routes

import (
	"taskflow/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// RegisterJobRoutes registers job-related routes
func RegisterDashboardRoutes(router *gin.Engine, dashboardHandler *handlers.DashboardHandler) {
	router.GET("/dashboard", func(c *gin.Context) {
		// Check if it's an upgrade request (for WebSocket)
		if websocket.IsWebSocketUpgrade(c.Request) {
			dashboardHandler.StreamStats(c) // WebSocket connection
		} else {
			dashboardHandler.DisplayStats(c) // Render HTML page
		}
	})
}
