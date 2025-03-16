package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"taskflow/internal/bootstrap"
	repositories "taskflow/internal/domain/repository"
	"taskflow/internal/handlers"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	_ = godotenv.Load() // Load environment variables

	// Initialize Database
	database, err := bootstrap.InitDatabase()
	if err != nil {
		logrus.Fatal("Failed to connect to the database")
	}
	repo := repositories.NewGormJobRepository(database)

	// Initialize Workers
	jobQueue, workerPool := bootstrap.InitWorkerPool(repo)
	defer workerPool.Stop()

	// Initialize Router
	jobHandler := handlers.NewJobHandler(repo, *jobQueue)

	dashboardHanlder := handlers.NewDashboardHandler(repo, workerPool)

	router := bootstrap.SetupRouter(jobHandler, dashboardHanlder)

	// Start HTTP Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Graceful Shutdown Handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		logrus.Info("Server running on :", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("listen: %s\n", err)
		}
	}()

	<-quit
	logrus.Info("Server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	log.Println("Server exited gracefully")
}
