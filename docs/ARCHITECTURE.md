# Architecture Guide - Finance Hub

## System Overview

Finance Hub là ứng dụng quản lý tài chính cá nhân với kiến trúc **Client-Server, Clean Architecture**, tách biệt rõ ràng giữa Frontend (React), Backend (Golang), và Database (MongoDB).

```
┌─────────────────────────────────────────────────────────────┐
│                        User/Browser                          │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                 Frontend (React + TypeScript)                │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Features: Dashboard, Transactions, Budgets, AI Chat │   │
│  │  Services: API Clients, State Management             │   │
│  │  Components: Shadcn/UI, Custom Components            │   │
│  └──────────────────────────────────────────────────────┘   │
└────────────────────────┬────────────────────────────────────┘
                         │ HTTPS/REST API
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Backend (Golang + Gin Framework)                │
│  ┌──────────────────────────────────────────────────────┐   │
│  │            Handlers (HTTP Controllers)                │   │
│  │  ┌────────────────────────────────────────────────┐  │   │
│  │  │           Services (Business Logic)             │  │   │
│  │  │  ┌──────────────────────────────────────────┐  │  │   │
│  │  │  │      Repositories (Data Access)          │  │  │   │
│  │  │  │  ┌────────────────────────────────────┐  │  │  │   │
│  │  │  │  │   Models (Domain Entities)         │  │  │  │   │
│  │  │  │  └────────────────────────────────────┘  │  │  │   │
│  │  │  └──────────────────────────────────────────┘  │  │   │
│  │  └────────────────────────────────────────────────┘  │   │
│  └──────────────────────────────────────────────────────┘   │
└────────────────────────┬────────────────────────────────────┘
                         │ MongoDB Driver
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                   MongoDB Atlas (Cloud)                      │
│  Collections: users, accounts, transactions, categories...   │
└─────────────────────────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│           External Services (Supabase, VietQR, AI)           │
└─────────────────────────────────────────────────────────────┘
```

---

## Architecture Principles

### 1. Clean Architecture (Hexagonal)

Backend tuân theo **Clean Architecture** với 4 layers:

```
┌────────────────────────────────────────────────────────┐
│                    Handlers Layer                       │  ◄─── HTTP/REST Interface
│  - Request/Response parsing                             │
│  - Validation (basic)                                   │
│  - Route mapping                                        │
└─────────────────────┬──────────────────────────────────┘
                      │ calls
                      ▼
┌────────────────────────────────────────────────────────┐
│                   Services Layer                        │  ◄─── Business Logic
│  - Business rules                                       │
│  - Orchestration                                        │
│  - Complex validation                                   │
│  - Transaction coordination                             │
└─────────────────────┬──────────────────────────────────┘
                      │ calls
                      ▼
┌────────────────────────────────────────────────────────┐
│                 Repositories Layer                      │  ◄─── Data Access
│  - CRUD operations                                      │
│  - Query building                                       │
│  - Data mapping                                         │
└─────────────────────┬──────────────────────────────────┘
                      │ uses
                      ▼
┌────────────────────────────────────────────────────────┐
│                    Models Layer                         │  ◄─── Domain Entities
│  - Data structures                                      │
│  - DTOs                                                 │
│  - Request/Response types                               │
└────────────────────────────────────────────────────────┘
```

**Benefits:**

- ✅ Dễ test (mock dependencies)
- ✅ Loose coupling
- ✅ Clear separation of concerns
- ✅ Dễ maintain và scale

### 2. Dependency Injection

Tất cả dependencies được inject thông qua constructors:

```go
// Repository depends on DB
accountRepo := repositories.NewAccountRepository(db.Database)

// Service depends on Repository
accountService := services.NewAccountService(accountRepo)

// Handler depends on Service
accountHandler := handlers.NewAccountHandler(accountService)
```

### 3. Single Responsibility

Mỗi layer có trách nhiệm rõ ràng:

- **Handlers**: HTTP concerns only
- **Services**: Business logic only
- **Repositories**: Data access only
- **Models**: Data structures only

---

## Technology Stack

### Backend

- **Language**: Go 1.22
- **Framework**: Gin (HTTP router)
- **Database**: MongoDB (via official mongo-driver)
- **Auth**: Supabase JWT validation
- **Config**: godotenv (environment variables)
- **Logging**: Custom logger (pkg/logger)

### Frontend

- **Framework**: React 18 + TypeScript
- **Build Tool**: Vite
- **UI Library**: Shadcn/UI (Radix UI + Tailwind)
- **State Management**: React Context API
- **HTTP Client**: fetch API
- **Auth**: Supabase Client
- **Charts**: Recharts
- **Date**: date-fns
- **Icons**: Lucide React

### Infrastructure

- **Database**: MongoDB Atlas (Cloud)
- **Auth Provider**: Supabase
- **Deployment**: Docker (ready)
- **CI/CD**: GitHub Actions (ready)

---

## Project Structure

### Backend Structure

```
server/
├── cmd/
│   └── api/
│       └── main.go              # Entry point, dependency wiring
├── internal/                    # Private application code
│   ├── config/
│   │   └── config.go            # Configuration loading & validation
│   ├── handlers/
│   │   ├── router.go            # Route definitions
│   │   ├── account_handler.go
│   │   ├── transaction_handler.go
│   │   ├── category_handler.go
│   │   ├── budget_handler.go
│   │   └── ...
│   ├── services/
│   │   ├── account_service.go
│   │   ├── transaction_service.go
│   │   └── ...                  # Business logic
│   ├── repositories/
│   │   ├── account_repository.go
│   │   └── ...                  # Data access
│   ├── models/
│   │   └── models.go            # Domain models & DTOs
│   └── middleware/
│       ├── auth.go              # JWT authentication
│       ├── cors.go
│       └── logger.go
├── pkg/                         # Public, reusable packages
│   ├── database/
│   │   └── database.go          # MongoDB connection
│   ├── logger/
│   │   └── logger.go            # Custom logger
│   └── response/
│       └── response.go          # Standardized API responses
├── docs/                        # Documentation
├── .env.example
├── go.mod
├── go.sum
└── Dockerfile
```

**Design Principles:**

- `cmd/`: Application entry points
- `internal/`: Private code (không thể import từ ngoài)
- `pkg/`: Public, reusable packages
- `docs/`: All documentation

### Frontend Structure

```
my-finance-hub/
├── src/
│   ├── components/
│   │   ├── ui/                  # Shadcn/UI components
│   │   ├── layout/              # App layout components
│   │   ├── shared/              # Shared components
│   │   └── auth/                # Auth components
│   ├── features/                # Feature-based organization
│   │   ├── dashboard/
│   │   ├── transactions/
│   │   ├── budgets/
│   │   ├── reports/
│   │   ├── ai-chat/
│   │   └── settings/
│   ├── services/                # API clients
│   │   ├── AccountService.ts
│   │   ├── TransactionService.ts
│   │   └── ...
│   ├── contexts/                # React contexts
│   │   └── AuthContext.tsx
│   ├── hooks/                   # Custom hooks
│   ├── models/
│   │   └── index.ts             # TypeScript interfaces
│   ├── utils/                   # Utility functions
│   ├── lib/                     # Third-party configs
│   ├── App.tsx
│   └── main.tsx
├── public/
├── index.html
├── vite.config.ts
├── tailwind.config.ts
└── package.json
```

**Design Principles:**

- **Feature-based**: Code organized by features, not file types
- **Co-location**: Related code stays together
- **Separation**: UI components separate from business logic

---

## Key Design Patterns

### 1. Repository Pattern

Abstracts data access, dễ dàng swap database implementation.

```go
type AccountRepository interface {
    Create(userID string, req CreateAccountRequest) (*Account, error)
    GetByID(id, userID string) (*Account, error)
    GetAll(userID string, pagination PaginationQuery) ([]Account, int, error)
    Update(id, userID string, req UpdateAccountRequest) (*Account, error)
    Delete(id, userID string) error
    UpdateBalance(accountID string, delta float64) error
}

type accountRepositoryImpl struct {
    db *mongo.Database
}

func NewAccountRepository(db *mongo.Database) AccountRepository {
    return &accountRepositoryImpl{db: db}
}
```

### 2. Service Pattern

Encapsulates business logic, orchestrates multiple repositories.

```go
type AccountService struct {
    repo AccountRepository
}

func (s *AccountService) CreateAccount(userID string, req CreateAccountRequest) (*Account, error) {
    // Validation
    if !isValidAccountType(req.Type) {
        return nil, errors.New("invalid account type")
    }

    // Business logic
    // ...

    // Delegate to repository
    return s.repo.Create(userID, req)
}
```

### 3. Middleware Pattern

Chain of responsibilities cho HTTP request processing.

```go
router := gin.Default()
router.Use(middleware.Logger())
router.Use(middleware.CORS())

// Protected routes
auth := router.Group("/api/v1")
auth.Use(middleware.AuthMiddleware(cfg.Supabase))
{
    auth.GET("/accounts", accountHandler.GetAllAccounts)
}
```

### 4. Context Pattern (Frontend)

Global state management cho authentication và app state.

```tsx
const AppContext = createContext<AppContextType>()

export function AppProvider({ children }) {
  const [selectedMonth, setSelectedMonth] = useState("2026-02")
  const [selectedAccount, setSelectedAccount] = useState(null)

  return (
    <AppContext.Provider value={{ selectedMonth, selectedAccount, ... }}>
      {children}
    </AppContext.Provider>
  )
}
```

---

## Data Flow

### 1. Transaction Creation Flow

```
User fills form
    ↓
Frontend validates
    ↓
POST /api/v1/transactions
    ↓
Handler validates request
    ↓
Handler calls TransactionService.CreateTransaction()
    ↓
Service validates business rules (balance check)
    ↓
Service calls TransactionRepository.Create()
    ↓
Repository inserts to MongoDB
    ↓
Service calls AccountService.UpdateBalance()
    ↓
Repository updates account balance
    ↓
Response flows back through layers
    ↓
Frontend updates UI
```

### 2. Budget Alert Flow

```
Transaction created
    ↓
Service recalculates budget spent
    ↓
Service checks if spent > threshold
    ↓
If true: Service calls AlertService.CreateAlert()
    ↓
Alert stored in MongoDB
    ↓
Frontend polls /api/v1/alerts
    ↓
Alert displayed in UI
```

### 3. AI Chat Flow

```
User sends message
    ↓
POST /api/v1/ai/chat
    ↓
Handler calls AIChatService.ProcessMessage()
    ↓
Service analyzes intent (spending, budget, forecast...)
    ↓
Service fetches relevant data (transactions, budgets...)
    ↓
Service generates AI response
    ↓
Response with answer cards returned
    ↓
Frontend displays chat message + cards
```

---

## Authentication Flow

### Supabase JWT Authentication

```
1. User registers/logs in via Supabase
   ↓
2. Supabase returns JWT token
   ↓
3. Frontend stores token in memory/localStorage
   ↓
4. Frontend includes token in Authorization header
   ↓
5. Backend middleware validates JWT
   ↓
6. Middleware extracts user_id from JWT claims
   ↓
7. user_id injected into Gin context
   ↓
8. Handlers access user_id via c.Get("user_id")
```

**JWT Claims:**

```json
{
    "sub": "user-uuid", // user_id
    "email": "user@example.com",
    "role": "authenticated",
    "iat": 1234567890,
    "exp": 1234571490
}
```

**Middleware Implementation:**

```go
func AuthMiddleware(cfg SupabaseConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extract token from Authorization header
        token := extractToken(c)

        // Validate JWT with Supabase public key
        claims, err := validateJWT(token, cfg.JWTSecret)
        if err != nil {
            c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
            return
        }

        // Inject user_id into context
        c.Set("user_id", claims.Sub)
        c.Next()
    }
}
```

---

## Error Handling

### Backend Error Strategy

```go
// 1. Service returns domain errors
func (s *Service) DoSomething() error {
    if condition {
        return errors.New("business rule violated")
    }
    return nil
}

// 2. Handler translates to HTTP responses
func (h *Handler) HandleRequest(c *gin.Context) {
    err := h.service.DoSomething()
    if err != nil {
        if errors.Is(err, ErrNotFound) {
            response.NotFoundResponse(c, "Resource")
            return
        }
        response.InternalErrorResponse(c, err)
        return
    }
    response.SuccessResponse(c, 200, "Success", data)
}
```

### Frontend Error Strategy

```tsx
try {
    const result = await TransactionService.createTransaction(data);
    toast.success("Transaction created!");
} catch (error) {
    if (error.status === 400) {
        toast.error(error.message); // Show validation error
    } else if (error.status === 401) {
        navigate("/login"); // Redirect to login
    } else {
        toast.error("An error occurred"); // Generic error
    }
}
```

---

## State Management (Frontend)

### Global State (Context API)

**Use for:**

- User authentication state
- Selected month (app-wide filter)
- Selected account
- Theme preference

```tsx
const AppContext = createContext({
    user: null,
    selectedMonth: "2026-02",
    selectedAccount: null,
    setSelectedMonth: () => {},
    setSelectedAccount: () => {},
});
```

### Local State (useState)

**Use for:**

- Component-specific state
- Form inputs
- UI toggles (modal open/close)

```tsx
const [isModalOpen, setIsModalOpen] = useState(false)
const [formData, setFormData] = useState({...})
```

### Server State (React Query - Future)

**Consider for:**

- Caching API responses
- Auto-refetch on window focus
- Optimistic updates

---

## API Client Pattern (Frontend)

### Service Layer Pattern

```typescript
class AccountServiceClass {
    private baseURL = "/api/v1";

    async listAccounts(): Promise<Account[]> {
        const response = await fetch(`${this.baseURL}/accounts`, {
            headers: {
                Authorization: `Bearer ${getToken()}`,
            },
        });

        if (!response.ok) {
            throw new ApiError(response.status, await response.json());
        }

        return response.json();
    }
}

export const AccountService = new AccountServiceClass();
```

**Benefits:**

- Centralized API calls
- Easy to mock for testing
- Type-safe with TypeScript
- Reusable across components

---

## Security Considerations

### Backend Security

1. **JWT Validation**: All protected routes validate Supabase JWT
2. **User Isolation**: All queries filter by `user_id` from JWT
3. **Input Validation**: Validate all inputs at handler level
4. **SQL/NoSQL Injection**: Use parameterized queries (MongoDB BSON)
5. **CORS**: Configure allowed origins
6. **Rate Limiting**: Implement per-user rate limits
7. **HTTPS Only**: Force HTTPS in production

### Frontend Security

1. **Token Storage**: Store JWT in memory (or httpOnly cookie)
2. **XSS Prevention**: React auto-escapes by default
3. **CSRF**: Not applicable (using JWT, not cookies)
4. **Validation**: Client-side validation + server-side validation
5. **Sanitization**: Sanitize user inputs (especially file uploads)

---

## Performance Optimization

### Backend Optimizations

1. **Database Indexing**: Index frequently queried fields
2. **Connection Pooling**: Reuse MongoDB connections
3. **Caching**: Cache budget calculations, forecast results
4. **Pagination**: Always paginate large result sets
5. **Batch Operations**: Bulk delete/update when possible

### Frontend Optimizations

1. **Code Splitting**: Lazy load features
2. **Image Optimization**: Compress images, use WebP
3. **Debouncing**: Debounce search inputs
4. **Memoization**: Use React.memo for expensive components
5. **Virtual Scrolling**: For large transaction lists

### MongoDB Optimizations

```javascript
// Compound index for common query
db.transactions.createIndex({
    user_id: 1,
    date_time_iso: -1,
});

// Denormalize for performance
// Store category_name in transaction (read optimization)
```

---

## Testing Strategy

### Backend Testing

```go
// Unit tests for services
func TestAccountService_CreateAccount(t *testing.T) {
    mockRepo := &MockAccountRepository{}
    service := services.NewAccountService(mockRepo)

    account, err := service.CreateAccount("user-id", validRequest)

    assert.NoError(t, err)
    assert.NotNil(t, account)
}

// Integration tests for repositories
func TestAccountRepository_Create(t *testing.T) {
    db := setupTestDB()
    repo := repositories.NewAccountRepository(db)

    account, err := repo.Create("user-id", validRequest)

    assert.NoError(t, err)
    // Verify in database
}
```

### Frontend Testing

```tsx
// Component tests
it("renders transaction list", () => {
    render(<TransactionList transactions={mockData} />);
    expect(screen.getByText("Phở 24")).toBeInTheDocument();
});

// Service mocking
jest.mock("@/services/TransactionService");
```

---

## Deployment Architecture

### Development

```
localhost:5173 (Frontend - Vite)
    ↓
localhost:8080 (Backend - Go)
    ↓
MongoDB Atlas (Cloud)
```

### Production (Docker)

```
┌──────────────────────────────────────┐
│         Nginx Reverse Proxy           │
│  - SSL Termination                    │
│  - Static file serving (frontend)     │
│  - /api/* → backend:8080              │
└──────────────┬───────────────────────┘
               │
    ┌──────────┴──────────┐
    ▼                     ▼
┌─────────┐         ┌─────────┐
│ Frontend│         │ Backend │
│Container│         │Container│
│(Nginx)  │         │(Go)     │
└─────────┘         └────┬────┘
                         │
                         ▼
                  ┌──────────────┐
                  │ MongoDB Atlas│
                  └──────────────┘
```

**Docker Compose:**

```yaml
services:
    backend:
        build: ./server
        ports:
            - "8080:8080"
        environment:
            - MONGODB_URI=${MONGODB_URI}
            - SUPABASE_URL=${SUPABASE_URL}

    frontend:
        build: ./my-finance-hub
        ports:
            - "80:80"
        depends_on:
            - backend
```

---

## Monitoring & Logging

### Backend Logging

```go
logger.Log.Info.Printf("Transaction created: %s", transaction.ID)
logger.Log.Error.Printf("Database error: %v", err)
```

### Metrics to Track

- Request latency (p50, p95, p99)
- Error rate by endpoint
- Database query performance
- Active users
- Transaction volume

### Tools (Future)

- **Logging**: ELK Stack (Elasticsearch, Logstash, Kibana)
- **Monitoring**: Prometheus + Grafana
- **Tracing**: Jaeger
- **Alerts**: PagerDuty

---

## Scalability Considerations

### Horizontal Scaling

- Backend: Stateless, scale with multiple containers
- Database: MongoDB sharding (khi cần)
- Load balancer: Nginx or AWS ALB

### Vertical Scaling

- MongoDB: Increase cluster tier
- Backend: Increase container resources

### Caching Strategy

```
Redis Cache Layer
    ↓
- User session data
- Budget calculations (5 min TTL)
- Forecast results (1 hour TTL)
- Category list (cache invalidate on update)
```

---

## Future Architecture Enhancements

### 1. Microservices (nếu cần scale)

```
API Gateway
    ├── Account Service
    ├── Transaction Service
    ├── Analytics Service
    ├── AI Service
    └── Notification Service
```

### 2. Event-Driven Architecture

```
Transaction Created Event
    ├── Update Account Balance
    ├── Recalculate Budget
    ├── Check Alert Thresholds
    ├── Trigger Forecast Regeneration
    └── Log to Analytics
```

### 3. GraphQL API (alternative to REST)

Single endpoint, client-specified queries.

### 4. Real-time Updates

```
WebSocket Server
    ↓
Real-time notifications
- New transaction alerts
- Budget threshold warnings
- AI insights
```

---

**Last Updated**: February 27, 2026  
**Version**: 1.0
