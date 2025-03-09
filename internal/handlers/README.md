## Understanding `internal/handlers/`

The `internal/handlers/` package is responsible for handling HTTP requests in Taskflow. It defines route handlers that interact with the service and repository layers to manage jobs and system health.

### 1. `base_handlers.go`
- Contains generic API handlers, such as:
  - `HomeHandler`: Serves the home page using an HTML template.
  - `PingHandler`: Provides a health check endpoint that returns a JSON response with the current timestamp.

### 2. `job_handlers.go`
- Manages job-related API endpoints, including:
  - `CreateJob`: Accepts job creation requests, inserts them into the database, and enqueues them for processing.
  - Uses dependency injection to interact with the `JobRepository` and `JobQueue` components.
  - Ensures proper error handling and response formatting.

This package helps structure API request handling, making it easier to extend and maintain Taskflow's functionality.

