# Todo Service

Go backend service for managing todo items with file upload capabilities. Built with Clean Architecture principles, featuring MySQL, Redis streaming, and comprehensive testing with GoMock.


**
Implement gqlgen, graphQL, beeORM
**

**Create graphQL:**
```json
mutation CreateTodo($input: CreateTodoInput!) {
  createTodo(input: $input) {
    todo{
      id
      description
      dueDate
      fileId
      createdAt
    }
     error{
      field
      message
    }
  }
}
```

**get graphQL by id :**
```json
query GetTodo($id: ID!) {
  todo(id: $id) {
    id
    description
    dueDate
    createdAt
    
  }
}

```
```json
{"id": "1"}

```

===================================================================

```json

{
  "input": {
    "description": "sohel item!",
    "dueDate": "2025-12-27T10:00:00Z",
    "fileId": "file123"
  }
}
```

**Search graphQL:**
```json
query search{
  searchTodos(query: "item") {
    total
    results {
      id
      description
    }
  }
}

```



## Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)
- [Testing](#testing)
- [Project Structure](#project-structure)

---

## Features

- Create and manage todo items
- File upload with metadata storage
- Redis Stream integration for event publishing
- MySQL database
- Comprehensive unit tests with GoMock
- Performance benchmarking
- Docker containerization
- Automatic database migrations
- Health check endpoints
- Clean Architecture design

---

## Architecture

The application follows Clean Architecture principles with clear separation of concerns:

```
GraphQL Request
    ↓
internal/api/graphql/handler/ (Infrastructure - Input Adapter)
    ↓
internal/port/interface.go (Interface Layer - Input Port)
    ↓
internal/usecase/todo/ (Application Layer)
    ↓
internal/port/port.go (Interface Layer - Output Port)
    ↓
internal/repository/beeorm/ (Infrastructure - Output Adapter)
    ↓
MySQL Database
```

### Layers:

- **Domain**: Pure domain entities with no business logic
- **Use Cases**: Business logic independent of frameworks
- **Interfaces**: Abstraction contracts for services
- **Infrastructure**: Database, storage, and messaging implementations

---

## Prerequisites

- Go 1.25
- Docker & Docker Compose
- Make
- Install mockgen (for autogenerate mock) on your computer 

``` bash 
go install github.com/golang/mock/mockgen@latest 
```

---

## Installation

### Step 1: Clone the Repository

```bash
git clone <repository-url>
cd todo-service
```


### Env file

Please create .env file in the project root and values like this And added a .env.sample file included in the project root.

```bash

# GO APP Configuration
APP_PORT=8080

# Database Configuration
DB_HOST=mysql
DB_PORT=3306
DB_USER=todo_user
DB_PASSWORD=todo_password
DB_NAME=todo_db

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_ADDR=redis:6379
REDIS_STREAM=todos

# S3 Configuration (LocalStack)
S3_ENDPOINT=http://localstack:4566
S3_BUCKET=todo-bucket
AWS_ACCESS_KEY_ID=test
AWS_SECRET_ACCESS_KEY=test
AWS_DEFAULT_REGION=us-east-1


# Application Configuration
MAX_FILE_SIZE=5242880
LOG_LEVEL=info

```

you can change port and config in the .env file. 

---

## Running the Application

### Using Make (Recommended)

The project includes a Makefile with convenient commands:

```bash
# Start all services with Docker Compose (builds image first)
make run

# Run unit tests
make test

# Run performance benchmarks
make benchmark

# Generate mock implementations
make generate-mocks
```

### Using Docker Compose Directly

```bash
# Start all services (MySQL, Redis, LocalStack, Go App)
docker compose up --build


# Stop all services
docker compose down
```


### Health Check

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "time": "2025-10-18T14:54:44Z"
}
```

---

## API Endpoints

### 1. Create Todo

Create a new todo item with optional file reference.

**Endpoint:** `POST /todo`

**Request Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
  "description": "Buy groceries from market",
  "dueDate": "2025-10-25T15:30:00Z",
  "fileId": "optional-file-id"
}
```

**Example using curl:**
```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "description": "Buy groceries from market",
    "dueDate": "2025-10-25T15:30:00Z",
  }' \
  http://localhost:8080/todo
```

**Success Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "description": "Buy groceries from market",
  "dueDate": "2025-10-25T15:30:00Z",
  "createdAt": "2025-10-18T14:54:44Z",
  "updatedAt": "2025-10-18T14:54:44Z"
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "Description is required"
}
```

**Error Response (500 Internal Server Error):**
```json
{
  "error": "Failed to create todo: database connection error"
}
```

---

### 2. Upload File

Upload a file and get a file ID reference.

**Endpoint:** `POST /upload`

**Request Headers:**
```
Content-Type: multipart/form-data
```

**Request Parameters:**
- `file` (required): The file to upload (max 5MB)

**Example using curl:**
```bash
# Create a test file
echo "This is test content" > test.txt

# Upload the file
curl -X POST \
  -F "file=@test.txt" \
  -F "uploadedBy=john_doe" \
  http://localhost:8080/upload
```

**Success Response (201 Created):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "fileName": "123e4567-e89b-12d3-a456-426614174000.txt",
  "originalName": "test.txt",
  "contentType": "text/plain",
  "fileSize": 20,
  "fileHash": "a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6",
  "uploadedAt": "2025-10-18T14:54:44Z"
}
```

**Error Response (400 Bad Request - File Too Large):**
```json
{
  "error": "Upload failed: file too large: max size is 5242880 bytes"
}
```

**Error Response (400 Bad Request - Invalid File Type):**
```json
{
  "error": "Upload failed: invalid file type: only image/* and text/* are allowed"
}
```

**Error Response (500 Internal Server Error):**
```json
{
  "error": "Upload failed: failed to upload file to S3: service unavailable"
}
```

---

## Complete Workflow Example

### Step 1: Upload a File

```bash
echo "Project requirements" > project.txt

UPLOAD_RESPONSE=$(curl -s -X POST \
  -F "file=@project.txt" \
  http://localhost:8080/upload)

# Extract file ID
FILE_ID=$(echo "$UPLOAD_RESPONSE" | grep -o '"id":"[^"]*' | cut -d'"' -f4)
echo "Uploaded File ID: $FILE_ID"

```

Important: correcct file ID is needed, it uses as valid file refrence.

### Step 2: Create Todo with File Reference

```bash
curl -X POST \
  -H "Content-Type: application/json" \
  -d "{
    \"description\": \"Review project requirements\",
    \"dueDate\": \"2025-10-20T10:00:00Z\",
    \"fileId\": \"$FILE_ID\"
  }" \
  http://localhost:8080/todo
```

---

## Testing

### Using Make

## Makefile Commands Reference

The project includes a convenient Makefile with the following targets:

### Available Commands

```makefile
.PHONY: run test benchmark generate-mocks

# Start all services with Docker Compose
make run
  - Builds Docker image
  - Starts MySQL, Redis, LocalStack, and Go App
  - Output: Services running on http://localhost:8080

# Run unit tests
make test
  - Runs all tests in ./test/unit/...
  - Shows verbose output
  - Output: Test results and coverage

# Run performance benchmarks
make benchmark
  - Executes benchmarks in ./test/bench/...
  - Shows memory allocation statistics
  - Output: Benchmark results with ns/op, B/op, allocs/op

# Generate mock implementations
make generate-mocks
  - Creates GoMock mock implementations
  - Destination: test/mock/genrate-mocks.go
  - Generates: MockTodoRepository, MockFileRepository, MockS3Repository, MockRedisStreamRepository
```

### Project structure.

```
todo-service/
├── cmd/
│   └── main.go                 # Application entry point
|   └── cmd.go             
├── internal/
│   ├── domain/
│   │   ├── entity/            # Domain models
│   │   └── interface/         # Service interfaces
│   ├── usecase/               # Business logic
│   ├── handler/               # HTTP handlers
│   └── infrastructure/        # Implementation
│       ├── repository/        # Database operations
│       ├── storage/           # S3 storage
│       ├── stream/            # Redis streaming
│       ├── migration/         # Database migrations
│       └── helper/            # Utility functions
├── test/
│   ├── unit/                  # Unit tests
│   ├── bench/                 # Benchmarks
│   └── mock/                  # Mock implementations
├── migrations/                # SQL migration files
├── docker-compose.yml         # Docker services
├── Dockerfile                 # Go app container
├── go.mod & go.sum            # Dependencies
└── README.md                  # This file
```



---

The following file types are allowed for upload:

- **Images**: `image/jpeg`, `image/png`, `image/gif`, `image/webp`
- **Documents**: `text/plain`, `text/csv`, `application/pdf`

Maximum file size: 5MB

---


## Development

### Adding New Endpoints

1. Create use case in `internal/usecase/`
2. Add handler in `internal/handler/`
3. Register route in `handler.RegisterRoutes()`
4. Write tests in `test/unit/`

### Database Migrations

Add new migrations in `migrations/` folder following the naming convention:
```
000003_migration_name.up.sql
000003_migration_name.down.sql
```

Migrations run automatically on application startup.

---



### Check logs on Docker

Check bucket storage on localstack

```bash
docker ps
docker exec -it <localstack-container-id> awslocal s3 ls
```

check file how many files uploaded to localstack

```bash
docker exec -it <localstack-container-id> awslocal s3 ls s3://todo-bucket 
```

check redis cli

```bash
docker exec -it <redis-container-id> redis-cli
```

then enter the redis cli and give this command for check stream message
```bash
XRANGE todos - +
```
