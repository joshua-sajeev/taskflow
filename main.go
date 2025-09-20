package main

import (
	"fmt"
	"log"
	"time"

	"taskflow/internal/domain/task"
	"taskflow/internal/domain/user"
	"taskflow/internal/handler/task"
	"taskflow/internal/handler/user"
	"taskflow/internal/repository/gorm/gorm_task"
	"taskflow/internal/repository/gorm/gorm_user"
	"taskflow/internal/service/task"
	"taskflow/internal/service/user"
	"taskflow/pkg"

	docs "taskflow/docs"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// @title           TaskFlow API
// @version         1.0
// @description     API server for managing tasks and users in the TaskFlow application.
// @host            localhost:8080
// @BasePath        /api
func main() {
	dbuser := pkg.GetEnv("MYSQL_USER", "appuser")
	pass := pkg.GetEnv("MYSQL_PASSWORD", "apppassword")
	host := pkg.GetEnv("MYSQL_HOST", "172.18.0.2")
	port := pkg.GetEnv("MYSQL_PORT", "3306")
	dbname := pkg.GetEnv("MYSQL_DATABASE", "taskdb")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbuser, pass, host, port, dbname)
	log.Printf("Connecting with user %s to %s:%s/%s", dbuser, host, port, dbname)

	var db *gorm.DB
	var err error

	for range 5 {
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

	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	if err := db.AutoMigrate(&task.Task{}, &user.User{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	// Dependency wiring
	taskRepo := gorm_task.NewTaskRepository(db)
	taskSvc := task_service.NewTaskService(taskRepo)
	taskHandler := task_handler.NewTaskHandler(taskSvc)

	userRepo := gorm_user.NewUserRepository(db)
	userSvc := user_service.NewUserService(userRepo)
	userHandler := user_handler.NewUserHandler(userSvc)

	// Router setup
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.Schemes = []string{"http"}

	api := r.Group("/api")
	{
		api.POST("/auth/register", userHandler.Register)
		api.POST("/auth/login", userHandler.Login)

		taskRoutes := api.Group("/tasks")
		taskRoutes.Use(user_handler.AuthMiddleware())
		{
			taskRoutes.POST("", taskHandler.CreateTask)
			taskRoutes.GET("/:id", taskHandler.GetTask)
			taskRoutes.GET("", taskHandler.ListTasks)
			taskRoutes.PATCH("/:id/status", taskHandler.UpdateStatus)
			taskRoutes.DELETE("/:id", taskHandler.Delete)
		}

		userRoutes := api.Group("/users")
		userRoutes.Use(user_handler.AuthMiddleware())
		{
			userRoutes.PATCH("/password", userHandler.UpdatePassword)
			userRoutes.DELETE("/account", userHandler.DeleteUser)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	log.Println("Routes registered: /auth, /tasks (protected), /users (protected)")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
