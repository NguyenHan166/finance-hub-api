# Finance Hub API

A robust, scalable backend API for personal finance management built with Go.

## 🏗️ Architecture

This project follows **Clean Architecture** principles with clear separation of concerns:

```
server/
├── cmd/
│   └── api/              # Application entry point
├── internal/             # Private application code
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP request handlers
│   ├── services/        # Business logic layer
│   ├── repositories/    # Data access layer
│   ├── models/          # Domain models
│   ├── middleware/      # HTTP middleware
│   └── utils/           # Shared utilities
├── pkg/                 # Public shared libraries
│   ├── database/        # Database connection
│   ├── logger/          # Logging utilities
│   └── validator/       # Input validation
└── docs/                # API documentation
```

## 🚀 Features

- **RESTful API** with Gin framework
- **Clean Architecture** for maintainability
- **PostgreSQL** with Supabase integration
- **JWT Authentication** with Supabase Auth
- **Request Validation** and error handling
- **Structured Logging**
- **CORS** support
- **Hot Reload** with Air
- **Docker** support

## 📋 Prerequisites

- Go 1.22 or higher
- PostgreSQL (via Supabase)
- Make (optional)

## 🛠️ Setup

1. **Clone and navigate to server directory:**

    ```bash
    cd server
    ```

2. **Install dependencies:**

    ```bash
    go mod download
    ```

3. **Configure environment variables:**

    ```bash
    cp .env.example .env
    # Edit .env with your configuration
    ```

4. **Run the application:**

    ```bash
    # Development mode with hot reload
    make dev

    # Or run directly
    go run cmd/api/main.go
    ```

## 📝 Available Commands

```bash
# Run in development mode (with hot reload)
make dev

# Build the application
make build

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Run linter
make lint

# Clean build files
make clean

# Install development tools
make install-tools
```

## 🔑 Environment Variables

See `.env.example` for all available configuration options.

Key variables:

- `PORT`: Server port (default: 8080)
- `DB_*`: PostgreSQL/Supabase database configuration
- `SUPABASE_*`: Supabase API keys
- `JWT_SECRET`: Secret key for JWT signing
- `CORS_ALLOWED_ORIGINS`: Allowed origins for CORS

## 📚 API Documentation

Once the server is running, API documentation is available at:

- Swagger UI: `http://localhost:8080/api/v1/docs`

## 🏃 Development Workflow

1. Create feature branch
2. Implement changes in appropriate layer
3. Write tests
4. Run `make fmt` and `make lint`
5. Commit and push

## 📁 Layer Responsibilities

### Models (`internal/models`)

- Define domain entities
- No business logic
- Pure data structures

### Repositories (`internal/repositories`)

- Database operations (CRUD)
- Query building
- Transaction management
- No business logic

### Services (`internal/services`)

- Business logic implementation
- Orchestrate repository calls
- Data transformation
- Business rule validation

### Handlers (`internal/handlers`)

- HTTP request/response handling
- Input validation
- Call appropriate services
- Format responses

### Middleware (`internal/middleware`)

- Request/response interception
- Authentication/Authorization
- Logging, CORS, Rate limiting

## 🧪 Testing

```bash
# Run all tests
make test

# Run specific package tests
go test ./internal/services/...

# With coverage
make test-coverage
```

## 🐳 Docker

```bash
# Build image
docker build -t finance-hub-api .

# Run container
docker run -p 8080:8080 --env-file .env finance-hub-api
```

## 📈 Scaling Considerations

- Stateless design for horizontal scaling
- Database connection pooling
- Caching ready (Redis can be added)
- Rate limiting middleware
- Structured logging for monitoring

## 🤝 Contributing

1. Follow Go best practices
2. Maintain clean architecture boundaries
3. Write tests for new features
4. Update documentation

## 📄 License

MIT
