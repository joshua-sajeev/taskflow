# Taskflow - Job Scheduler

Taskflow is a **job scheduler** built using Golang. It allows users to enqueue jobs, process them concurrently using worker pools, and manage job execution efficiently. The project is designed to handle background task execution using a queue system.

## Understanding `main.go`
The `main.go` file is responsible for initializing and running the Taskflow application. Below is a breakdown of its key components:

1. **Loading Environment Variables:**
   ```go
   _ = godotenv.Load()
   ```
   Loads environment variables from a `.env` file.

2. **Database Initialization:**
   ```go
   database, err := bootstrap.InitDatabase()
   if err != nil {
       logrus.Fatal("Failed to connect to the database")
   }
   repo := repositories.NewGormJobRepository(database)
   ```
   Establishes a connection to PostgreSQL and initializes the job repository.

3. **Worker Pool Initialization:**
   ```go
   jobQueue, workerPool := bootstrap.InitWorkerPool(repo)
   defer workerPool.Stop()
   ```
   Creates a job queue and starts worker processes to handle tasks concurrently.

4. **Router Setup:**
   ```go
   jobHandler := handlers.NewJobHandler(repo, *jobQueue)
   router := bootstrap.SetupRouter(jobHandler)
   ```
   Initializes the API router and assigns request handlers.

5. **Starting the HTTP Server:**
   ```go
   port := os.Getenv("PORT")
   if port == "" {
       port = "8080"
   }
   srv := &http.Server{
       Addr:    ":" + port,
       Handler: router,
   }
   ```
   Starts the server on the specified port.

6. **Graceful Shutdown Handling:**
   ```go
   quit := make(chan os.Signal, 1)
   signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
   
   <-quit
   logrus.Info("Server shutting down")
   
   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
   defer cancel()
   
   if err := srv.Shutdown(ctx); err != nil {
       log.Fatal("Server Shutdown:", err)
   }
   ```
   Ensures that when the server is stopped, it shuts down cleanly, allowing workers to complete tasks before exiting.

