# Development Guide

This guide covers how to set up your development environment, run the application locally, and contribute to TaskFlow.

---

## Prerequisites

- **Go** 1.25+ ([Download](https://golang.org/dl/))
- **Docker & Docker Compose** ([Download](https://www.docker.com/products/docker-desktop))
- **Git** ([Download](https://git-scm.com/))
- **Make** (optional, for running commands)

### Verify Installation

```bash
go version          # Go 1.25+
docker --version    # Docker 28.5+
docker-compose --version  # Docker Compose 2.40+
git --version
```

---

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/yourusername/taskflow.git
cd taskflow
```

### 2. Set Up Environment Variables

Create a `.env` file in the root directory:

```bash
# JWT Configuration
JWT_SECRET=ADD-YOUR-SECRET

# Database Configuration
MYSQL_USER=NEW-USER
MYSQL_PASSWORD=CHANGE-PASSWORD
MYSQL_HOST=mysql
MYSQL_PORT=3306
MYSQL_DATABASE=CHANGE-DB-NAME
```

### 3. Start Services with Docker Compose

```bash
# Development mode with hot-reload
docker-compose -f docker-compose.override.yml up

# Or production mode
docker-compose up
```

This starts:
- **MySQL 8.0** on port 3306
- **Go API** on port 8080 (with hot-reload in dev)
- **MySQL** initialization with `.env` settings

### 4. Verify the Application

```bash
# Check API is running
curl http://localhost:8080/swagger/index.html

# Or test an endpoint
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

---

## Development Workflow

### Without Docker (Local Development)

If you prefer running Go locally:

```bash
# 1. Start only MySQL
docker-compose up mysql

# 2. Download dependencies
go mod download

# 3. Install development tools
go install github.com/air-verse/air@latest
go install github.com/swaggo/swag/cmd/swag@latest

# 4. Run with Air (hot-reload)
air

# 5. In another terminal, generate Swagger docs
swag init
```

**Air** watches for file changes and automatically rebuilds the app.

---

## Running Tests

### Unit Tests

```bash
# Run all unit tests
go test ./...

# Run tests for specific package
go test ./internal/service/task/...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Integration Tests

```bash
# Run integration tests (requires Docker)
cd test/integration
go test -v .

# Or from root
go test -v ./test/integration/...
```

---

## Code Structure & Best Practices

### Adding a New Endpoint

1. **Create Handler** in `internal/handler/{entity}/`
   ```go
   func (h *TaskHandler) MyNewAction(c *gin.Context) {
       // Implementation
   }
   ```

2. **Add to Interface** in `internal/handler/{entity}/{entity}_handler_interface.go`
   ```go
   type TaskHandlerInterface interface {
       MyNewAction(c *gin.Context)
   }
   ```

3. **Add Router** in `main.go`
   ```go
   taskRoutes.GET("/new-action", taskHandler.MyNewAction)
   ```

4. **Add Tests** in `internal/handler/{entity}/{entity}_handler_test.go`
   ```go
   func TestTaskHandler_MyNewAction(t *testing.T) {
       // Test cases
   }
   ```

### Adding Business Logic

1. **Add Service Method** in `internal/service/{entity}/{entity}_service.go`
   ```go
   func (s *{Entity}Service) NewBusinessLogic(/* params */) error {
       // Validation & logic
   }
   ```

2. **Add to Interface** in `internal/service/{entity}/{entity}_service_interface.go`

3. **Call from Handler**
   ```go
   if err := h.service.NewBusinessLogic(/* args */); err != nil {
       c.JSON(http.StatusBadRequest, common.ErrorResponse{Message: err.Error()})
       return
   }
   ```

4. **Test the Service** in `internal/service/{entity}/{entity}_service_test.go`
   ```go
   func TestService_NewBusinessLogic(t *testing.T) {
       mockRepo := new(gorm_entity.MockRepository)
       // Mock setup
       service := NewService(mockRepo)
       // Assertions
   }
   ```

### Adding Database Operations

1. **Add Repository Method** in `internal/repository/gorm/gorm_{entity}/{entity}_repository.go`
   ```go
   func (r *{Entity}Repository) NewOperation(/* params */) error {
       return r.db.Where("...").Do(...)
   }
   ```

2. **Add to Interface** in `internal/repository/gorm/gorm_{entity}/{entity}_repository_interface.go`

3. **Add Mock** in `internal/repository/gorm/gorm_{entity}/{entity}_repository_mock.go`

4. **Test with SQLite** in `internal/repository/gorm/gorm_{entity}/{entity}_repository_test.go`

---

## Database Operations

### View Database

```bash
# Connect to MySQL container
docker-compose exec mysql mysql -u <username> -p <db-name>

# Inside MySQL shell
SHOW TABLES;
SELECT * FROM users;
SELECT * FROM tasks;
DESCRIBE users;
```

### Reset Database

```bash
# Drop all tables (careful!)
docker-compose exec mysql mysql -u <username> -p <db-name> -e "DROP TABLE tasks; DROP TABLE users;"

# Restart services
docker-compose down
docker-compose up
```

### Migrations

Migrations run automatically on startup via `database.MigrateModels()` in `main.go`:

```go
if err := database.MigrateModels(db, &user.User{}, &task.Task{}); err != nil {
    log.Fatal(err)
}
```

To add migrations:
1. Modify entity struct in `internal/domain/{entity}/entity.go`
2. Restart the application (GORM auto-migrates)

---

## API Documentation

### Swagger UI

Swagger documentation is auto-generated from code comments:

```bash
# View docs (running locally)
http://localhost:8080/swagger/index.html

# Generate new docs after changes
swag init
```

### Adding Swagger Comments

```go
// GetTask godoc
// @Summary Get a task by ID
// @Description Returns details of a specific task by ID
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID" minimum(1) example(1)
// @Success 200 {object} dto.GetTaskResponse
// @Failure 404 {object} common.ErrorResponse
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
    // Implementation
}
```

---

## Common Commands

```bash
# Build the application
go build -o taskflow .

# Run linter (requires golangci-lint)
golangci-lint run

# Format code
go fmt ./...

# Run go vet (catch common mistakes)
go vet ./...

# Download/verify dependencies
go mod download
go mod verify

# Tidy dependencies
go mod tidy

# View dependency tree
go mod graph
```

---
### Using Logs

The application uses Go's standard `log` package:

```go
log.Println("Debug message")
log.Printf("Formatted message: %v", value)
log.Fatal("Fatal error")
```

View logs in Docker:

```bash
docker-compose logs -f goapp  # Follow logs
docker-compose logs goapp     # View history
```

---

## Troubleshooting

### MySQL Connection Issues

**Error**: `Error 1045 (28000): Access denied for user`

```bash
# Check credentials in .env
# Verify MySQL is healthy
docker-compose ps

# Restart MySQL
docker-compose restart mysql

# View MySQL logs
docker-compose logs mysql
```

### Port Already in Use

```bash
# Kill process using port 8080
lsof -i :8080
kill -9 <PID>

# Or use different port in docker-compose.override.yml
ports:
  - "8081:8080"
```

### Tests Failing

```bash
# Clear test cache
go clean -testcache

# Run with verbose output
go test -v ./...

# Run specific test
go test -run TestName ./package/...
```

### Hot-Reload Not Working

```bash
# Check Air is running
docker-compose logs goapp

# Restart service
docker-compose restart goapp

# Check .air.toml config exists
cat .air.toml
```

---

## Resources

- [Go Documentation](https://golang.org/doc/)
- [Gin Web Framework](https://gin-gonic.com/)
- [GORM Documentation](https://gorm.io/)
- [JWT-Go](https://github.com/golang-jwt/jwt)
- [Docker Documentation](https://docs.docker.com/)
- [Air Hot Reload](https://github.com/cosmtrek/air)