package main

import (
	"log"
	"time"

	"taskflow/internal/auth"
	"taskflow/internal/domain/task"
	"taskflow/internal/domain/user"
	task_handler "taskflow/internal/handler/task"
	user_handler "taskflow/internal/handler/user"
	"taskflow/internal/middleware/ratelimiter"
	"taskflow/internal/repository/gorm/gorm_task"
	"taskflow/internal/repository/gorm/gorm_user"
	task_service "taskflow/internal/service/task"
	user_service "taskflow/internal/service/user"
	"taskflow/pkg"
	"taskflow/pkg/database"

	docs "taskflow/docs"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
)

// @title           TaskFlow API
// @version         1.0
// @description     API server for managing tasks and users in the TaskFlow application.
// @host            localhost:8080
// @BasePath        /api
func main() {

	cfg := database.LoadConfigFromEnv()
	db, err := database.ConnectDB(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := database.MigrateModels(db, &user.User{}, &task.Task{}); err != nil {
		log.Fatal(err)
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	secretKey := []byte(pkg.GetEnv("JWT_SECRET", ""))
	// Dependency wiring
	taskRepo := gorm_task.NewTaskRepository(db)
	taskSvc := task_service.NewTaskService(taskRepo)
	userRepo := gorm_user.NewUserRepository(db)
	userSvc := user_service.NewUserService(userRepo, string(secretKey))

	userAuth := auth.NewUserAuth(string(secretKey), userRepo)

	taskHandler := task_handler.NewTaskHandler(taskSvc, userAuth)
	userHandler := user_handler.NewUserHandler(userSvc, userAuth)

	// Rate limiter setup for auth endpoints
	// Allows 5 requests per second with a burst of 10 requests
	authRateLimiter := ratelimiter.NewIPRateLimiter(rate.Limit(5), 10)
	// Clean up old IP entries every hour to prevent memory leaks
	authRateLimiter.StartCleanupRoutine(1 * time.Hour)

	// Router setup
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.Schemes = []string{"http"}

	api := r.Group("/api")
	{
		// Auth routes with rate limiting
		authRoutes := api.Group("/auth")
		authRoutes.Use(authRateLimiter.Middleware())
		{
			authRoutes.POST("/register", userHandler.Register)
			authRoutes.POST("/login", userHandler.Login)
		}

		taskRoutes := api.Group("/tasks")
		taskRoutes.Use(userAuth.AuthMiddleware())
		{
			taskRoutes.POST("", taskHandler.CreateTask)
			taskRoutes.GET("/:id", taskHandler.GetTask)
			taskRoutes.GET("", taskHandler.ListTasks)
			taskRoutes.PATCH("/:id/status", taskHandler.UpdateStatus)
			taskRoutes.DELETE("/:id", taskHandler.Delete)
		}

		userRoutes := api.Group("/users")
		userRoutes.Use(userAuth.AuthMiddleware())
		{
			userRoutes.PATCH("/password", userHandler.UpdatePassword)
			userRoutes.DELETE("/account", userHandler.DeleteUser)
		}
	}

	public := r.Group("/")
	{
		// Auth routes with rate limiting
		authRoutes := public.Group("/auth")
		authRoutes.Use(authRateLimiter.Middleware())
		{
			authRoutes.POST("/register", userHandler.Register)
			authRoutes.POST("/login", userHandler.Login)
		}

		taskRoutes := public.Group("/tasks")
		taskRoutes.Use(userAuth.AuthMiddleware())
		{
			taskRoutes.POST("", taskHandler.CreateTask)
			taskRoutes.GET("/:id", taskHandler.GetTask)
			taskRoutes.GET("", taskHandler.ListTasks)
			taskRoutes.PATCH("/:id/status", taskHandler.UpdateStatus)
			taskRoutes.DELETE("/:id", taskHandler.Delete)
		}

		userRoutes := public.Group("/users")
		userRoutes.Use(userAuth.AuthMiddleware())
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
