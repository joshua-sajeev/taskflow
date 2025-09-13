# TaskFlow

A REST API for task management built with Go, following clean architecture principles.

## What I've Built

### 📋 Core Features
- **Task Management**: Create, read, update, delete tasks
- **Status Updates**: Change task status (pending → in-progress → completed)
- **Clean Architecture**: Separated layers (handler → service → repository → database)
- **JWT Authentication**: Complete JWT utility package with token creation/validation

### 🏗️ Architecture Implementation
```
internal/
├── domain/task/        # Task entity and repository interface
├── dto/               # Request/response data structures  
├── handler/           # HTTP handlers with comprehensive tests
├── service/           # Business logic layer
├── repository/gorm/   # Database layer with GORM
└── common/           # Shared error responses

pkg/
└── jwt.go            # JWT utilities (create, validate, extract username)
```

### 🧪 Testing Coverage
- **Handler Tests**: All HTTP endpoints with success/failure scenarios
- **Service Tests**: Business logic validation with mocks
- **Repository Tests**: Database operations with in-memory SQLite
- **JWT Tests**: Token creation, validation, and error handling
- **Race Detection**: `go test -race` ready

### 🐳 Docker Setup
- **Development**: Hot reload with Air, volume mounting for live coding
- **Production**: Multi-stage Alpine build for optimized images
- **Database**: MySQL 8.0 with health checks and proper networking
- **Docker Compose**: Separate configs for dev/prod environments

### 📊 API Endpoints Built
| Method | Endpoint | Status |
|--------|----------|--------|
| `POST` | `/api/tasks` | ✅ Complete |
| `GET` | `/api/tasks` | ✅ Complete |
| `GET` | `/api/tasks/:id` | ✅ Complete |
| `PATCH` | `/api/tasks/:id/status` | ✅ Complete |
| `DELETE` | `/api/tasks/:id` | ✅ Complete |

### 🔐 JWT Implementation
- Token creation with configurable expiration
- Token validation with proper error handling
- Username extraction from tokens
- Comprehensive test coverage
- Ready for authentication middleware integration

## Tech Stack
- **Go** with Gin framework
- **MySQL** with GORM ORM
- **Docker** & Docker Compose
- **Testify** for testing
- **Swagger** documentation (integrated)

## Quick Start

1. **Clone and setup**:
```bash
git clone <repo>
cd taskflow
```

2. **Create `.env` file**:
```env
MYSQL_ROOT_PASSWORD=root
MYSQL_DATABASE=taskdb
MYSQL_USER=appuser
MYSQL_PASSWORD=apppassword
```

3. **Run with Docker**:
```bash
# Development with hot reload
docker-compose -f docker-compose.yml -f docker-compose.override.yml up

# Production
docker-compose up
```

4. **Access**: 
   - API: http://localhost:8080/api
   - Swagger: http://localhost:8080/swagger/

## Testing

```bash
# All tests
go test ./...

# With race detection  
go test -race ./...

# With coverage
go test -cover ./...
```

## What's Next
Planning to transform this into a **concurrent, high-performance system** with:
- Redis integration for caching/sessions
- Worker pools for parallel processing  
- WebSocket real-time updates
- Event-driven architecture
- User authentication system

---

*Built with Go clean architecture principles and comprehensive testing*
