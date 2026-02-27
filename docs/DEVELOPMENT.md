# Development Guide

## Prerequisites

- Go 1.22+
- PostgreSQL 15+ (or Supabase account)
- Make (optional but recommended)
- Air (for hot reload)
- Docker & Docker Compose (optional)

## Getting Started

### 1. Clone and Setup

```bash
cd server
cp .env.example .env
# Edit .env with your configuration
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Install Development Tools

```bash
make install-tools
```

This installs:

- `air` - Hot reload for development
- `golangci-lint` - Code linter

### 4. Database Setup

#### Option A: Using Supabase (Recommended for Production)

1. Create a Supabase project
2. Copy the connection details to `.env`
3. Run migrations in Supabase SQL Editor:
    - `migrations/001_create_tables.sql`
    - `migrations/002_seed_default_categories.sql`

#### Option B: Local PostgreSQL

```bash
# Using Docker Compose
docker-compose up -d postgres

# Or install PostgreSQL locally and create database
createdb finance_hub
psql finance_hub < migrations/001_create_tables.sql
psql finance_hub < migrations/002_seed_default_categories.sql
```

### 5. Run the Application

```bash
# Development mode with hot reload
make dev

# Or run directly
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

## Project Structure

```
server/
├── cmd/
│   └── api/              # Application entry point
│       └── main.go       # Main application file
│
├── internal/             # Private application code
│   ├── config/          # Configuration management
│   │   └── config.go
│   │
│   ├── handlers/        # HTTP request handlers (Controller layer)
│   │   ├── account_handler.go
│   │   ├── category_handler.go
│   │   ├── transaction_handler.go
│   │   ├── health_handler.go
│   │   └── router.go
│   │
│   ├── services/        # Business logic layer
│   │   ├── account_service.go
│   │   ├── category_service.go
│   │   └── transaction_service.go
│   │
│   ├── repositories/    # Data access layer
│   │   ├── account_repository.go
│   │   ├── category_repository.go
│   │   └── transaction_repository.go
│   │
│   ├── models/          # Domain models and DTOs
│   │   └── models.go
│   │
│   ├── middleware/      # HTTP middleware
│   │   └── middleware.go
│   │
│   └── utils/           # Utility functions
│       └── helpers.go
│
├── pkg/                 # Public shared libraries
│   ├── database/        # Database connection
│   │   └── database.go
│   ├── logger/          # Logging utilities
│   │   └── logger.go
│   └── response/        # HTTP response helpers
│       └── response.go
│
├── migrations/          # Database migrations
│   ├── 001_create_tables.sql
│   └── 002_seed_default_categories.sql
│
├── docs/                # Documentation
│   └── API.md
│
├── .env.example         # Environment variables template
├── .air.toml           # Air configuration (hot reload)
├── docker-compose.yml  # Docker Compose configuration
├── Dockerfile          # Docker image definition
├── go.mod              # Go module definition
├── Makefile            # Development commands
└── README.md           # Project documentation
```

## Architecture

This project follows **Clean Architecture** principles:

### Layer Dependencies

```
Handlers → Services → Repositories → Database
```

Each layer only depends on the layer below it:

1. **Handlers** (Presentation Layer)
    - Handle HTTP requests/responses
    - Validate input
    - Call services
    - Format responses

2. **Services** (Business Logic Layer)
    - Implement business rules
    - Orchestrate repository calls
    - Transform data
    - No HTTP concerns

3. **Repositories** (Data Access Layer)
    - Database operations (CRUD)
    - Query building
    - No business logic

4. **Models** (Domain Layer)
    - Define entities
    - Data structures
    - No logic

### Benefits

- **Testable**: Each layer can be tested independently
- **Maintainable**: Clear separation of concerns
- **Scalable**: Easy to add new features
- **Flexible**: Easy to swap implementations

## Development Workflow

### 1. Adding a New Feature

Example: Adding a Budget feature

#### Step 1: Define Model

```go
// internal/models/models.go
type Budget struct {
    ID         uuid.UUID
    UserID     uuid.UUID
    CategoryID uuid.UUID
    Amount     float64
    Period     string
    // ...
}
```

#### Step 2: Create Repository

```go
// internal/repositories/budget_repository.go
type BudgetRepository struct {
    db *sql.DB
}

func (r *BudgetRepository) Create(...) (*models.Budget, error) {
    // Database logic
}
```

#### Step 3: Create Service

```go
// internal/services/budget_service.go
type BudgetService struct {
    repo *repositories.BudgetRepository
}

func (s *BudgetService) CreateBudget(...) (*models.Budget, error) {
    // Business logic
    return s.repo.Create(...)
}
```

#### Step 4: Create Handler

```go
// internal/handlers/budget_handler.go
type BudgetHandler struct {
    service *services.BudgetService
}

func (h *BudgetHandler) CreateBudget(c *gin.Context) {
    // HTTP logic
    result, err := h.service.CreateBudget(...)
    response.SuccessResponse(c, http.StatusCreated, "Budget created", result)
}
```

#### Step 5: Register Routes

```go
// internal/handlers/router.go
budgets := protected.Group("/budgets")
{
    budgets.POST("", budgetHandler.CreateBudget)
    budgets.GET("", budgetHandler.GetAllBudgets)
    // ...
}
```

### 2. Testing

```bash
# Run all tests
make test

# Run specific package tests
go test ./internal/services/...

# Run tests with coverage
make test-coverage
```

### 3. Code Quality

```bash
# Format code
make fmt

# Run linter
make lint
```

### 4. Before Committing

```bash
make fmt
make lint
make test
```

## Common Tasks

### Database Migration

1. Create new migration file in `migrations/`
2. Run migration on database:
    ```bash
    psql $DATABASE_URL < migrations/003_your_migration.sql
    ```

### Adding Dependencies

```bash
go get github.com/pkg/errors
go mod tidy
```

### Debugging

1. Use VS Code debugger with launch configuration
2. Or add debug logs:
    ```go
    logger.Log.Debug.Printf("Debug info: %+v", data)
    ```

## Environment Variables

See `.env.example` for all available variables.

Key variables for development:

- `ENV=development` - Development mode
- `DB_SSLMODE=disable` - For local PostgreSQL
- `LOG_LEVEL=debug` - Verbose logging

## Docker Development

```bash
# Start all services
docker-compose up

# Rebuild after code changes
docker-compose up --build

# Stop all services
docker-compose down

# View logs
docker-compose logs -f api
```

## Troubleshooting

### Port already in use

```bash
# Find and kill process using port 8080
lsof -ti:8080 | xargs kill -9
```

### Database connection issues

1. Check PostgreSQL is running
2. Verify credentials in `.env`
3. Check firewall settings

### Hot reload not working

```bash
# Reinstall Air
go install github.com/cosmtrek/air@latest
```

## Best Practices

1. **Always validate input** in handlers
2. **Keep business logic in services** not handlers
3. **Use repositories for all database access**
4. **Write tests** for new features
5. **Follow naming conventions**
6. **Use meaningful commit messages**
7. **Update documentation** when adding features

## Resources

- [Go Documentation](https://golang.org/doc/)
- [Gin Framework](https://gin-gonic.com/)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Supabase Documentation](https://supabase.com/docs)
