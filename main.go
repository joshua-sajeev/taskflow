package main

import (
	"fmt"
	"log"
	"os"
	"taskflow/internal/domain/task"
	"taskflow/internal/handler"
	gg "taskflow/internal/repository/gorm"
	"taskflow/internal/service"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func main() {
	user := getEnv("MYSQL_USER", "appuser")
	pass := getEnv("MYSQL_PASSWORD", "apppassword")
	host := getEnv("MYSQL_HOST", "mysql")
	port := getEnv("MYSQL_PORT", "3306")
	dbname := getEnv("MYSQL_DATABASE", "taskdb")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, dbname)
	log.Printf("Connecting with user %s to %s:%s/%s", user, host, port, dbname)

	var db *gorm.DB
	var err error

	// Retry logic
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Println("DB open error:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Println("Failed to get sql.DB from GORM:", err)
			time.Sleep(5 * time.Second)
			continue
		}

		err = sqlDB.Ping()
		if err == nil {
			log.Println("Database connected successfully")
			break
		}

		log.Println("Waiting for DB:", err)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("Database connection failed after retries: %v", err)
	}

	// Optional cleanup on exit
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// Auto-migrate your Task schema
	if err := db.AutoMigrate(&task.Task{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	// Dependency wiring
	repo := gg.NewTaskRepository(db)
	svc := service.NewTaskService(repo)
	h := handler.NewTaskHandler(svc)

	// Router setup
	r := gin.Default()
	r.POST("/tasks", h.CreateTask)
	r.GET("/tasks/:id", h.GetTask)
	r.GET("/tasks", h.ListTasks)
	r.PATCH("/tasks/:id/status", h.UpdateStatus)

	// Start the HTTP server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
