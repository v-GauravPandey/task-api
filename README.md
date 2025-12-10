# task-api
Task Api Repository

# Task Management API (Go + Fiber)

A production-ready Task Management REST API implemented in Go with Fiber, following DDD-style layering and TDD principles. Uses an in-memory repository (no external DB required).

## Features

✅ Create, Read, Update, Delete tasks  
✅ List tasks with pagination and status filtering  
✅ Input validation (required fields, future dates, status enum)  
✅ Error handling with meaningful messages  
✅ DDD layering: Domain → Repository → Service → Handler  
✅ In-memory data store  
✅ Unit + Integration tests  
✅ Postman collection for easy testing  

## Tech Stack

- **Go** 1.22+
- **Fiber** v2 (lightweight HTTP framework)
- **UUID** for task ID generation
- **testify** for assertions in tests

## Installation & Setup

### Prerequisites

- Go 1.22 or higher installed
- Postman (for testing, optional but recommended)

### Clone & Install Dependencies

```bash
git clone https://github.com/yourname/task-api.git
cd task-api

go mod download
go mod tidy
```

## Running the Server

```bash
go run ./cmd/server
```

The API will be available at: **`http://localhost:8080`**

## Running Tests

### All tests:
```bash
go test ./...
```

### Specific test file:
```bash
go test ./tests -v
```

### With coverage:
```bash
go test ./... -cover
```

---

## API Endpoints

All endpoints are prefixed with `/api`.

### 1. Create Task
**POST** `/api/tasks`

**Request Body:**
```json
{
  "title": "Complete report",
  "description": "Quarterly sales report",
  "status": "PENDING",
  "due_date": "2025-12-31T23:59:59Z"
}
```

**Required Fields:**
- `title` (string, non-empty)
- `due_date` (string, ISO8601 format, must be in the future)

**Optional Fields:**
- `description` (string)
- `status` (enum: `PENDING`, `IN_PROGRESS`, `DONE`; default: `PENDING`)

**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Complete report",
  "description": "Quarterly sales report",
  "status": "PENDING",
  "due_date": "2025-12-31T23:59:59Z"
}
```

**Error (400 Bad Request):**
```json
{
  "error": "title is required"
}
```

---

### 2. Get Task by ID
**GET** `/api/tasks/{id}`

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Complete report",
  "description": "Quarterly sales report",
  "status": "PENDING",
  "due_date": "2025-12-31T23:59:59Z"
}
```

**Error (404 Not Found):**
```json
{
  "error": "task not found"
}
```

---

### 3. Update Task
**PUT** `/api/tasks/{id}`

**Request Body (all fields optional for partial update):**
```json
{
  "title": "Updated title",
  "description": "Updated description",
  "status": "IN_PROGRESS",
  "due_date": "2025-12-25T18:00:00Z"
}
```

**Response (200 OK):** Returns updated task object

**Error (404 Not Found):** If task doesn't exist

---

### 4. Delete Task
**DELETE** `/api/tasks/{id}`

**Response (204 No Content):** Empty response body

**Error (404 Not Found):** If task doesn't exist

---

### 5. List All Tasks
**GET** `/api/tasks`

**Query Parameters:**
- `status` (optional): Filter by status (`PENDING`, `IN_PROGRESS`, `DONE`)
- `page` (optional, default=1): Page number for pagination
- `page_size` (optional, default=10): Number of items per page

**Examples:**
```
GET /api/tasks
GET /api/tasks?status=PENDING
GET /api/tasks?page=2&page_size=20
GET /api/tasks?status=IN_PROGRESS&page=1&page_size=5
```

**Response (200 OK):**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "title": "Task 1",
    "description": "Description 1",
    "status": "PENDING",
    "due_date": "2025-12-20T10:00:00Z"
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "title": "Task 2",
    "description": "Description 2",
    "status": "IN_PROGRESS",
    "due_date": "2025-12-31T23:59:59Z"
  }
]
```

---

## Testing with Postman

### Import Collection

1. Open **Postman**
2. Click **Import** → **Upload Files**
3. Select the provided `task-api-postman.json` file
4. Collection will appear in your workspace

### Set Environment Variable

1. In Postman, click the **Environment** icon (eye icon)
2. Edit **base_url** variable if your server runs on a different port:
   - Default: `http://localhost:8080`
   - Change if needed and save

### Run Requests

The collection includes:
- ✅ **Create Task** (with auto-capture of task ID)
- ✅ **Create Task - Minimal** (required fields only)
- ✅ **Create Task - Missing Title** (validation test)
- ✅ **Create Task - Past Due Date** (validation test)
- ✅ **Get Task by ID** (uses saved task_id)
- ✅ **Get Task - Not Found** (error case)
- ✅ **Update Task** (full update)
- ✅ **Update Task - Partial** (only title)
- ✅ **Update Task - Not Found** (error case)
- ✅ **List All Tasks**
- ✅ **List Tasks - Paginated**
- ✅ **List Tasks - Filtered by Status**
- ✅ **List Tasks - Filtered & Paginated**
- ✅ **Delete Task**
- ✅ **Delete Task - Not Found** (error case)

Each request includes **test assertions** that validate the response.

### Sample Workflow

1. **Create Task** - Returns task ID, auto-saves to `{{task_id}}`
2. **Get Task** - Fetches the created task using saved ID
3. **Update Task** - Modifies the task
4. **List Tasks** - Views all tasks with optional filters
5. **Delete Task** - Removes the task

---

## Project Structure

```
task-api/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point
├── internal/
│   ├── domain/
│   │   ├── task.go                 # Domain model + validation constants
│   │   └── task_service.go         # Business logic (interfaces + impl)
│   ├── repository/
│   │   ├── task_repository.go      # Repository interface
│   │   └── task_repository_memory.go # In-memory implementation
│   └── transport/
│       └── http/
│           ├── task_handler.go     # HTTP handlers
│           └── router.go           # Fiber app setup
├── pkg/
│   └── errors/
│       └── errors.go               # Custom error types
├── tests/
│   ├── task_service_test.go        # Unit tests
│   └── task_handler_integration_test.go # Integration tests
├── go.mod
├── go.sum
├── task-api-postman.json           # Postman collection
└── README.md
```

---

## Design Principles

### Domain-Driven Design (DDD)

- **Domain Layer**: `task.go` - Pure domain model with validation rules
- **Repository Layer**: `task_repository*.go` - Data access abstraction
- **Service Layer**: `task_service.go` - Business logic orchestration
- **Handler Layer**: `task_handler.go` - HTTP request/response mapping

### Clean Code

- Clear naming conventions
- Single Responsibility Principle
- Dependency Injection
- Interface-based design for testability
- Error handling with custom error types

### Test-Driven Development

- Unit tests for service logic
- Integration tests for HTTP endpoints
- Mock in-memory repository for isolation
- Test helpers and fixtures

---

## Error Handling

### Validation Errors (400 Bad Request)
```json
{
  "error": "title is required"
}
```

### Not Found Errors (404 Not Found)
```json
{
  "error": "task not found"
}
```

### Internal Server Errors (500)
```json
{
  "error": "internal error"
}
```

---

## Development & Contribution

### Running with Hot Reload (optional)

Install `air` for hot reloading:
```bash
go install github.com/cosmtrek/air@latest
air
```

### Adding New Features

1. Define domain logic in `domain/`
2. Implement repository methods if needed
3. Add service methods
4. Create HTTP handlers
5. Write tests first (TDD)
6. Update Postman collection if adding endpoints

---

## Assumptions

- In-memory storage (data lost on restart)
- Single-threaded for simplicity (but uses sync.RWMutex for safety)
- UUID v4 for task IDs
- ISO8601 date format (RFC3339)
- No authentication/authorization

---

## Future Enhancements

- [ ] PostgreSQL persistence
- [ ] JWT authentication
- [ ] Rate limiting
- [ ] Logging middleware
- [ ] Request validation middleware
- [ ] Swagger/OpenAPI docs
- [ ] Docker containerization
- [ ] CI/CD pipeline
- [ ] Deployment guide

---

## License

MIT

---

## Quick Start Commands

```bash
# Clone
git clone https://github.com/gauravpandey771/task-api.git && cd task-api

# Install deps
go mod tidy

# Run server
go run ./cmd/server

# Test (in another terminal)
go test ./... -v

# Import Postman collection and test all endpoints
```

Server running on `http://localhost:8080`
