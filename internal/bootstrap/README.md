# Taskflow - Job Scheduler
## Understanding `internal/bootstrap/`
The `internal/bootstrap/` package is responsible for initializing core components of the Taskflow application, including database setup, router configuration, and worker pool management.

### 1. `database.go`
Handles the initialization of the PostgreSQL database connection using GORM.

**Key Function:**
```go
func InitDatabase() (*gorm.DB, error)
```
- Calls `db.InitDB()` to establish the database connection.
- Logs a warning if the connection fails.

### 2. `router.go`
Configures the API routes and initializes the Gin router.

**Key Function:**
```go
func SetupRouter(jobHandler *handlers.JobHandler) *gin.Engine
```
- Calls `routes.SetupRouter()` to create a new Gin router.
- Registers job-related API endpoints.

### 3. `workers.go`
Initializes and manages the worker pool for processing jobs.

**Key Function:**
```go
func InitWorkerPool(repo repositories.JobRepository) (*jobqueue.JobQueue, *jobqueue.WorkerPool)
```
- Creates a `JobQueue` with a fixed size.
- Initializes a `WorkerPool` with a set number of workers.
- Starts the worker pool to begin processing jobs.

