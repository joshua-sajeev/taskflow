# Taskflow - Job Scheduler

## Understanding `internal/db/`
The `internal/db/` package is responsible for handling database connectivity and migrations for Taskflow.

### 1. `db.go`
Manages the database connection using GORM and ensures a singleton pattern for database access.

**Key Functions:**
```go
func InitDB() (*gorm.DB, error)
```
- Loads environment variables from `.env`.
- Establishes a connection to PostgreSQL.
- Returns a singleton database instance.

```go
func CloseDB()
```
- Closes the database connection safely.
- Ensures proper resource cleanup.

### 2. `migrations/`
Contains SQL migration files for managing database schema.

- **`000001_create_job_table.up.sql`**: Defines the schema for the `jobs` table.
- **`000001_create_job_table.down.sql`**: Rolls back the `jobs` table creation.

