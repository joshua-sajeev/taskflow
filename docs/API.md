# API Documentation

Complete reference for TaskFlow API endpoints, request/response formats, and status codes.

---

## Base URL

```
http://localhost:8080/api
```

---

## Authentication

TaskFlow uses **JWT (JSON Web Tokens)** for authentication.

### Getting a Token
1. [[API#Register User|Register]] or [[#Login User|Login]] to receive a token
2. Include token in all protected requests using:
   ```
   Authorization: Bearer <token>
   ```

### Token Details

- **Type**: JWT
- **Expiration**: 24 hours
- **Algorithm**: HS256
- **Storage**: Client-side (in header)

### Example Request with Token

```bash
curl -X GET http://localhost:8080/api/tasks \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

---

## Response Format

All responses follow a standard format:

### Success Response
```json
{
  "id": 1,
  "email": "user@example.com",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Error Response
```json
{
  "error": "invalid credentials"
}
```

### HTTP Status Codes

| Code | Meaning | Example |
|------|---------|---------|
| 200 | OK - Request succeeded | GET /tasks |
| 201 | Created - Resource created | POST /auth/register |
| 400 | Bad Request - Invalid input | Missing required field |
| 401 | Unauthorized - Invalid/missing token | Missing Authorization header |
| 404 | Not Found - Resource doesn't exist | GET /tasks/999 |
| 409 | Conflict - Duplicate resource | Email already exists |
| 500 | Server Error - Internal error | Database connection failed |

---

## Authentication Endpoints

### Register User

Create a new user account.

**Endpoint**: `POST /auth/register`

**Authentication**: None required

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Query Parameters**: None

**Validation**:
- `email`: Required, valid email format
- `password`: Required, minimum 6 characters

**Response** (201 Created):
```json
{
  "id": 1,
  "email": "user@example.com"
}
```

**Error Examples**:
```json
// Invalid email format
{
  "error": "invalid email format"
}

// Password too short
{
  "error": "password must be at least 6 characters"
}

// Email already exists
{
  "error": "email already exists"
}
```

**Example Request**:
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "myPassword123"
  }'
```

---

### Login User

Authenticate with email and password, receive JWT token.

**Endpoint**: `POST /auth/login`

**Authentication**: None required

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

**Response** (200 OK):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJpYXQiOjE2OTUzMTg1NjMsImV4cCI6MTY5NTQwNDk2M30.abc123...",
  "id": 1,
  "email": "user@example.com"
}
```

**Error Examples**:
```json
// Invalid credentials
{
  "error": "invalid credentials"
}

// User not found
{
  "error": "invalid credentials"
}
```

**Example Request**:
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token-here>" \
  -d '{
    "email": "john@example.com",
    "password": "myPassword123"
  }'
```

**Usage**: Store the returned `token` and include in `Authorization: Bearer <token>` header for all protected endpoints.

---

## Task Endpoints

All task endpoints require authentication.

### Create Task

Create a new task for the authenticated user.

**Endpoint**: `POST /tasks`

**Authentication**: Required ✓

**Request Body**:
```json
{
  "task": "Buy groceries"
}
```

**Validation**:
- `task`: Required, max 20 characters

**Response** (201 Created):
```json
{
  "task": "Buy groceries"
}
```

**Error Examples**:
```json
// Missing task field
{
  "error": "Key: 'CreateTaskRequest.Task' Error:Field validation for 'Task' failed on the 'required' tag"
}

// Task too long
{
  "error": "Key: 'CreateTaskRequest.Task' Error:Field validation for 'Task' failed on the 'max' tag"
}
```

**Example Request**:
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "task": "Complete project"
  }'
```

---

### Get Task

Retrieve a single task by ID.

**Endpoint**: `GET /tasks/{id}`

**Authentication**: Required ✓

**Path Parameters**:
- `id`: Task ID (integer, minimum 1)

**Query Parameters**: None

**Response** (200 OK):
```json
{
  "id": 1,
  "task": "Buy groceries",
  "status": "pending"
}
```

**Error Examples**:
```json
// Invalid ID format
{
  "error": "Invalid ID"
}

// Task not found or belongs to different user
{
  "error": "Task not found"
}
```

**Example Request**:
```bash
curl -X GET http://localhost:8080/api/tasks/1 \
  -H "Authorization: Bearer <token>"
```

---

### List Tasks

Retrieve all tasks for the authenticated user.

**Endpoint**: `GET /tasks`

**Authentication**: Required ✓

**Query Parameters**: None

**Response** (200 OK):
```json
{
  "tasks": [
    {
      "id": 1,
      "task": "Buy groceries",
      "status": "pending"
    },
    {
      "id": 2,
      "task": "Write report",
      "status": "completed"
    }
  ]
}
```

**Example Request**:
```bash
curl -X GET http://localhost:8080/api/tasks \
  -H "Authorization: Bearer <token>"
```

**Example Response with No Tasks**:
```json
{
  "tasks": []
}
```

---

### Update Task Status

Update the status of a task.

**Endpoint**: `PATCH /tasks/{id}/status`

**Authentication**: Required ✓

**Path Parameters**:
- `id`: Task ID (integer, minimum 1)

**Request Body**:
```json
{
  "status": "completed"
}
```

**Valid Status Values**:
- `pending` - Task not started
- `completed` - Task finished

**Response** (200 OK):
```json
{
  "message": "status updated"
}
```

**Error Examples**:
```json
// Invalid status value
{
  "error": "Key: 'UpdateStatusRequest.Status' Error:Field validation for 'Status' failed on the 'oneof' tag"
}

// Task not found
{
  "error": "Task not found"
}

// Invalid ID format
{
  "error": "invalid task ID"
}
```

**Example Request**:
```bash
curl -X PATCH http://localhost:8080/api/tasks/1/status \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "completed"
  }'
```

---

### Delete Task

Delete a task.

**Endpoint**: `DELETE /tasks/{id}`

**Authentication**: Required ✓

**Path Parameters**:
- `id`: Task ID (integer, minimum 1)

**Query Parameters**: None

**Response** (200 OK):
```json
{
  "message": "Task deleted successfully"
}
```

**Error Examples**:
```json
// Task not found
{
  "error": "Task not found"
}

// Invalid ID
{
  "error": "Invalid ID"
}
```

**Example Request**:
```bash
curl -X DELETE http://localhost:8080/api/tasks/1 \
  -H "Authorization: Bearer <token>"
```

---

## User Endpoints

### Update Password

Change the password for the authenticated user.

**Endpoint**: `PATCH /users/password`

**Authentication**: Required ✓

**Request Body**:
```json
{
  "old_password": "currentPassword123",
  "new_password": "newPassword456"
}
```

**Validation**:
- `old_password`: Required
- `new_password`: Required, minimum 6 characters

**Response** (200 OK):
```json
{
  "message": "Password updated successfully"
}
```

**Error Examples**:
```json
// Incorrect old password
{
  "error": "invalid old password"
}

// New password too short
{
  "error": "new password must be at least 6 characters"
}

// User not found
{
  "error": "user not found"
}
```

**Example Request**:
```bash
curl -X PATCH http://localhost:8080/api/users/password \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "currentPassword123",
    "new_password": "newPassword456"
  }'
```

---

### Delete Account

Permanently delete the authenticated user account.

**Endpoint**: `DELETE /users/account`

**Authentication**: Required ✓

**Request Body**: None

**Response** (200 OK):
```json
{
  "message": "User account deleted successfully"
}
```

**Note**: This performs a soft delete. The user can be recovered from backups but will not appear in normal queries.

**Error Examples**:
```json
// User not found (already deleted)
{
  "error": "user not found"
}
```

**Example Request**:
```bash
curl -X DELETE http://localhost:8080/api/users/account \
  -H "Authorization: Bearer <token>"
```

---

## Common Workflows

### Complete Flow: Register, Create Task, Update Status

```bash
# 1. Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
# Returns: { "id": 1, "email": "user@example.com" }

# 2. Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
# Returns: { "token": "eyJ...", "id": 1, "email": "user@example.com" }

# 3. Create task (use token from login)
TOKEN="eyJ..."
curl -X POST http://localhost:8080/api/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"task": "Buy milk"}'
# Returns: { "task": "Buy milk" }

# 4. List tasks
curl -X GET http://localhost:8080/api/tasks \
  -H "Authorization: Bearer $TOKEN"
# Returns: { "tasks": [{ "id": 1, "task": "Buy milk", "status": "pending" }] }

# 5. Update task status
curl -X PATCH http://localhost:8080/api/tasks/1/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status": "completed"}'
# Returns: { "message": "status updated" }

# 6. Delete task
curl -X DELETE http://localhost:8080/api/tasks/1 \
  -H "Authorization: Bearer $TOKEN"
# Returns: { "message": "Task deleted successfully" }
```

---

## Error Handling

### Common Errors

| Error | Cause | Solution |
|-------|-------|----------|
| `authorization header required` | Missing Authorization header | Add `Authorization: Bearer <token>` header |
| `invalid or expired token` | Token invalid or expired | Login again to get new token |
| `user account not found` | User deleted or doesn't exist | Register new account |
| `invalid credentials` | Wrong email or password | Check credentials |
| `email already exists` | Email already registered | Use different email |
| `invalid task ID` | Invalid ID format or < 1 | Use valid integer ID >= 1 |
| `Task not found` | Task doesn't exist or belongs to different user | Create task or check task ID |

### Error Response Format

All errors follow this format:
```json
{
  "error": "descriptive error message"
}
```

---

## Rate Limiting

Currently no rate limiting is implemented. Future version may include:
- Per-user rate limits
- Per-endpoint rate limits
- Time-window based throttling

---

## Pagination

Currently not implemented. Tasks endpoint returns all user tasks. Future version may include:
- `limit` query parameter
- `offset` query parameter
- Cursor-based pagination

---

## Filtering & Sorting

Not currently implemented. Future features:
- Filter tasks by status
- Filter tasks by date range
- Sort by creation date, status, etc.

---

## Swagger UI

Interactive API documentation available at:

```
http://localhost:8080/swagger/index.html
```

Browse and test all endpoints directly in the browser with:
- Request body builder
- Response examples
- Status code documentation

---

## Support & Issues

For API issues or questions:
- Check the [ARCHITECTURE.md](./ARCHITECTURE.md) for system design
- Review [DEVELOPMENT.md](./DEVELOPMENT.md) for setup help
- Open an issue on GitHub with detailed error messages