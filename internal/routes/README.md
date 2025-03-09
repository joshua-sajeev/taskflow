## Understanding `internal/routes/`

The `internal/routes/` package is responsible for defining API routes in Taskflow.

### 1. `job_routes.go`
- Defines job-related routes under the `/jobs` endpoint.
- Calls `RegisterJobRoutes`, which:
  - Groups job-related routes.
  - Maps HTTP `POST /jobs/` to `CreateJob` for job creation.

### 2. `routes.go`
- Defines the `SetupRouter` function, which:
  - Initializes the Gin router.
  - Loads HTML templates from the `templates/` directory.
  - Registers base routes (`/` for home, `/ping` for health check).

This package centralizes route management, ensuring clean and organized API definitions.

