## Understanding `internal/domain/entities/`
The `internal/entities/` package defines the core data structures for Taskflow.

### 1. `job.go`
Defines the `Job` struct, which represents a job in the system.

**Key Struct:**
```go
type Job struct {
    ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    Task      string    `gorm:"not null"`
    Status    string    `gorm:"default:pending"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
```
- `ID`: Unique identifier for the job.
- `Task`: The job description.
- `Status`: Tracks the job status (`pending`, `processing`, `completed`, or `failed`).
- `CreatedAt` and `UpdatedAt`: Timestamps for tracking job lifecycle.

**Key Functions:**
```go
func NewJob(task string) Job
```
- Creates a new `Job` instance with the default status `pending`.

```go
func (j *Job) UpdateStatus(status string) error
```
- Updates the job status while validating the input.

```go
func (j *Job) Execute()
```
- Simulates job execution by updating the status and introducing a delay.

```go
func Migrate(db *gorm.DB) error
```
- Automates database schema migrations for the `Job` table.

