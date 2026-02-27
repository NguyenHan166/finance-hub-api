# Development Guide - Finance Hub

## Table of Contents

1. [Getting Started](#getting-started)
2. [Development Environment](#development-environment)
3. [Project Setup](#project-setup)
4. [Development Workflow](#development-workflow)
5. [Code Style & Conventions](#code-style--conventions)
6. [Testing](#testing)
7. [Debugging](#debugging)
8. [Common Tasks](#common-tasks)
9. [Troubleshooting](#troubleshooting)

---

## Getting Started

### Prerequisites

**Backend:**

- Go 1.22 or later ([download](https://go.dev/dl/))
- MongoDB Atlas account (hoặc local MongoDB)
- Git

**Frontend:**

- Node.js 18+ và npm/yarn/bun ([download](https://nodejs.org/))
- Git

**Recommended Tools:**

- VS Code với Go extension
- MongoDB Compass (GUI for MongoDB)
- Postman or Thunder Client (API testing)
- Docker Desktop (optional)

---

## Development Environment

### Backend Setup

#### 1. Clone Repository

```bash
git clone <repository-url>
cd server
```

#### 2. Install Dependencies

```bash
go mod download
```

#### 3. Configure Environment

Create `.env` file:

```bash
cp .env.example .env
```

Edit `.env`:

```env
# Server Configuration
SERVER_PORT=8080
SERVER_ENV=development
API_VERSION=v1

# MongoDB Configuration
MONGODB_URI=mongodb+srv://nvhan166:han1662003@cluster0.evbdltl.mongodb.net/fmp_app?appName=Cluster0
MONGODB_DATABASE=fmp_app

# Supabase Configuration
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_JWT_SECRET=your-jwt-secret

# CORS
ALLOWED_ORIGINS=http://localhost:5173,http://localhost:3000
```

#### 4. Run Backend

```bash
# Development mode
go run cmd/api/main.go

# Or with hot reload (using Air)
air

# Build for production
go build -o finance-hub-api cmd/api/main.go
./finance-hub-api
```

Backend runs on `http://localhost:8080`

### Frontend Setup

#### 1. Navigate to Frontend Directory

```bash
cd my-finance-hub
```

#### 2. Install Dependencies

```bash
# Using npm
npm install

# Using yarn
yarn install

# Using bun
bun install
```

#### 3. Configure Environment

Create `.env.local`:

```env
VITE_SUPABASE_URL=https://your-project.supabase.co
VITE_SUPABASE_ANON_KEY=your-anon-key
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

#### 4. Run Frontend

```bash
# Development mode
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

Frontend runs on `http://localhost:5173`

---

## Project Setup

### Initial Database Setup

#### 1. Create MongoDB Collections

```javascript
// Connect to MongoDB Atlas via Compass or mongosh
use fmp_app

// Create collections (auto-created on first insert, but good to define)
db.createCollection("users")
db.createCollection("accounts")
db.createCollection("categories")
db.createCollection("transactions")
db.createCollection("budgets")
db.createCollection("alerts")

// Create indexes (important!)
db.accounts.createIndex({ user_id: 1, created_at: -1 })
db.transactions.createIndex({ user_id: 1, date_time_iso: -1 })
db.categories.createIndex({ user_id: 1, type: 1 })
db.budgets.createIndex({ user_id: 1, month: 1 })
```

#### 2. Seed Default Categories

Run seed script (create if needed):

```javascript
// seed.js
const defaultCategories = [
    { name: "Lương", type: "income", icon: "💰", color: "#10B981" },
    { name: "Ăn uống", type: "expense", icon: "🍜", color: "#F59E0B" },
    { name: "Di chuyển", type: "expense", icon: "🚗", color: "#3B82F6" },
    // ... more categories
];

// Insert for test user
const userId = "test-user-uuid";
defaultCategories.forEach((cat) => {
    db.categories.insertOne({
        _id: generateUUID(),
        user_id: userId,
        ...cat,
        is_default: true,
        created_at: new Date(),
        updated_at: new Date(),
    });
});
```

#### 3. Create Test User in Supabase

1. Go to Supabase dashboard
2. Authentication → Users → Add User
3. Note down user ID for testing

---

## Development Workflow

### Feature Development Flow

#### 1. Create Feature Branch

```bash
git checkout -b feature/transaction-filters
```

#### 2. Backend Development

```bash
# 1. Add model/DTO in internal/models/models.go
type TransactionFilters struct {
    Month      string `form:"month"`
    AccountID  string `form:"account_id"`
    CategoryID string `form:"category_id"`
}

# 2. Add repository method in internal/repositories/
func (r *TransactionRepository) GetFiltered(userID string, filters TransactionFilters) ([]Transaction, error) {
    // MongoDB query with filters
}

# 3. Add service method in internal/services/
func (s *TransactionService) GetFilteredTransactions(userID string, filters TransactionFilters) ([]Transaction, error) {
    // Business logic + validation
    return s.repo.GetFiltered(userID, filters)
}

# 4. Add handler in internal/handlers/
func (h *TransactionHandler) GetFiltered(c *gin.Context) {
    var filters TransactionFilters
    c.ShouldBindQuery(&filters)

    userID, _ := c.Get("user_id")
    transactions, err := h.service.GetFilteredTransactions(userID.(string), filters)

    response.SuccessResponse(c, 200, "Success", transactions)
}

# 5. Register route in internal/handlers/router.go
transactionGroup.GET("/", transactionHandler.GetFiltered)

# 6. Test manually with Postman/curl
curl "http://localhost:8080/api/v1/transactions?month=2026-02" \
  -H "Authorization: Bearer <token>"
```

#### 3. Frontend Development

```bash
# 1. Update model in src/models/index.ts
export interface TransactionFilters {
  month?: string
  accountId?: string
  categoryId?: string
}

# 2. Update service in src/services/TransactionService.ts
async listTransactions(filters?: TransactionFilters): Promise<Transaction[]> {
  const params = new URLSearchParams(filters)
  const response = await fetch(`${this.baseURL}/transactions?${params}`)
  return response.json()
}

# 3. Create/update component in src/features/transactions/
function TransactionFilterBar() {
  const [filters, setFilters] = useState<TransactionFilters>({})

  const handleFilterChange = (key, value) => {
    setFilters(prev => ({ ...prev, [key]: value }))
  }

  return (
    <div className="filter-bar">
      <Select onValueChange={v => handleFilterChange('month', v)}>
        {/* Month options */}
      </Select>
    </div>
  )
}

# 4. Use in page
function TransactionsPage() {
  const [filters, setFilters] = useState({})
  const [transactions, setTransactions] = useState([])

  useEffect(() => {
    TransactionService.listTransactions(filters)
      .then(setTransactions)
  }, [filters])

  return (
    <>
      <TransactionFilterBar onFilterChange={setFilters} />
      <TransactionList transactions={transactions} />
    </>
  )
}
```

#### 4. Testing

```bash
# Backend tests
cd server
go test ./internal/services/...
go test ./internal/repositories/...

# Frontend tests
cd my-finance-hub
npm run test
```

#### 5. Commit & Push

```bash
git add .
git commit -m "feat: add transaction filtering by month, account, category"
git push origin feature/transaction-filters
```

#### 6. Create Pull Request

1. Go to GitHub
2. Create PR from feature branch to main
3. Request review
4. Merge after approval

---

## Code Style & Conventions

### Backend (Go)

#### Naming Conventions

```go
// Exported (public) - PascalCase
type AccountService struct {}
func NewAccountService() {}

// Unexported (private) - camelCase
type accountRepository struct {}
func validateInput() {}

// Constants - UPPER_SNAKE_CASE or PascalCase
const MaxPageSize = 100
const DEFAULT_TIMEOUT = 10

// Interfaces - noun or adjective
type Reader interface {}
type AccountRepository interface {}
```

#### File Organization

```go
package services

import (
    // Standard library first
    "context"
    "fmt"
    "time"

    // External packages
    "github.com/gin-gonic/gin"
    "go.mongodb.org/mongo-driver/mongo"

    // Internal packages
    "finance-hub-api/internal/models"
    "finance-hub-api/internal/repositories"
)

// Constants
const DefaultLimit = 10

// Types
type AccountService struct {
    repo repositories.AccountRepository
}

// Constructor
func NewAccountService(repo repositories.AccountRepository) *AccountService {
    return &AccountService{repo: repo}
}

// Methods
func (s *AccountService) GetAccount(id string) (*models.Account, error) {
    // Implementation
}
```

#### Error Handling

```go
// Return errors, don't panic
func DoSomething() error {
    if err != nil {
        return fmt.Errorf("failed to do something: %w", err)
    }
    return nil
}

// Handle errors immediately
result, err := DoSomething()
if err != nil {
    return nil, err
}

// Use custom error types for domain errors
var ErrNotFound = errors.New("resource not found")
var ErrInsufficientBalance = errors.New("insufficient balance")
```

#### Documentation

```go
// Package documentation
// Package services contains business logic for the application.
package services

// Type documentation
// AccountService handles account-related business logic.
type AccountService struct {}

// Method documentation
// CreateAccount creates a new account for the user.
// Returns an error if the account type is invalid or database operation fails.
func (s *AccountService) CreateAccount(userID string, req CreateAccountRequest) (*Account, error) {
    // Implementation
}
```

### Frontend (TypeScript/React)

#### Naming Conventions

```tsx
// Components - PascalCase
function TransactionList() {}
function AccountCard() {}

// Hooks - camelCase with "use" prefix
function useTransactions() {}
function useAuth() {}

// Services - PascalCase class, singleton export
class TransactionServiceClass {}
export const TransactionService = new TransactionServiceClass();

// Constants - UPPER_SNAKE_CASE
const API_BASE_URL = "http://localhost:8080";
const MAX_FILE_SIZE = 5 * 1024 * 1024;

// Types/Interfaces - PascalCase
interface Transaction {}
type TransactionType = "income" | "expense";
```

#### File Organization

```tsx
// Imports - organized by source
import { useState, useEffect } from "react"; // React
import { format } from "date-fns"; // External libs
import { Button } from "@/components/ui/button"; // UI components
import { TransactionService } from "@/services"; // Services
import { Transaction } from "@/models"; // Types
import { formatCurrency } from "@/utils/format"; // Utils

// Types (if not imported)
interface Props {
    transactions: Transaction[];
    onSelect: (id: string) => void;
}

// Component
export function TransactionList({ transactions, onSelect }: Props) {
    // Hooks
    const [selected, setSelected] = useState<string | null>(null);

    // Effects
    useEffect(() => {
        // Side effects
    }, []);

    // Event handlers
    const handleClick = (id: string) => {
        setSelected(id);
        onSelect(id);
    };

    // Render
    return (
        <div className="transaction-list">
            {transactions.map((tx) => (
                <TransactionItem
                    key={tx.id}
                    transaction={tx}
                    onClick={handleClick}
                />
            ))}
        </div>
    );
}
```

#### TypeScript Best Practices

```tsx
// Use interfaces for object shapes
interface Account {
    id: string;
    name: string;
    balance: number;
}

// Use type for unions and primitives
type AccountType = "cash" | "bank" | "credit";
type Status = "loading" | "success" | "error";

// Avoid "any" - use proper types
// ❌ Bad
const data: any = await fetch();

// ✅ Good
const data: Account[] = await AccountService.listAccounts();

// Use optional chaining and nullish coalescing
const name = account?.name ?? "Unknown";
const balance = account?.balance ?? 0;
```

#### React Best Practices

```tsx
// Use functional components
function MyComponent() {} // ✅
class MyComponent {} // ❌ Avoid

// Destructure props
function Card({ title, amount }: Props) {} // ✅
function Card(props: Props) {} // ❌

// Use meaningful variable names
const [isModalOpen, setIsModalOpen] = useState(false); // ✅
const [open, setOpen] = useState(false); // ❌

// Extract complex logic to custom hooks
function useTransactions(filters: TransactionFilters) {
    const [transactions, setTransactions] = useState([]);
    const [loading, setLoading] = useState(false);

    useEffect(() => {
        setLoading(true);
        TransactionService.listTransactions(filters)
            .then(setTransactions)
            .finally(() => setLoading(false));
    }, [filters]);

    return { transactions, loading };
}

// Use composition over prop drilling
<AppProvider>
    <TransactionsPage /> {/* Can access context */}
</AppProvider>;
```

---

## Testing

### Backend Testing

#### Unit Tests (Services)

```go
// internal/services/account_service_test.go
package services

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock repository
type MockAccountRepository struct {
    mock.Mock
}

func (m *MockAccountRepository) Create(userID string, req CreateAccountRequest) (*Account, error) {
    args := m.Called(userID, req)
    return args.Get(0).(*Account), args.Error(1)
}

// Test
func TestAccountService_CreateAccount_Success(t *testing.T) {
    // Arrange
    mockRepo := new(MockAccountRepository)
    service := NewAccountService(mockRepo)

    expectedAccount := &Account{ID: "123", Name: "Test"}
    mockRepo.On("Create", "user-id", mock.Anything).Return(expectedAccount, nil)

    // Act
    result, err := service.CreateAccount("user-id", CreateAccountRequest{
        Name: "Test",
        Type: "cash",
        Balance: 1000,
    })

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expectedAccount, result)
    mockRepo.AssertExpectations(t)
}

func TestAccountService_CreateAccount_InvalidType(t *testing.T) {
    service := NewAccountService(nil)

    _, err := service.CreateAccount("user-id", CreateAccountRequest{
        Name: "Test",
        Type: "invalid",
        Balance: 1000,
    })

    assert.Error(t, err)
    assert.Contains(t, err.Error(), "invalid account type")
}
```

#### Integration Tests (Repositories)

```go
// internal/repositories/account_repository_test.go
package repositories

import (
    "context"
    "testing"
    "go.mongodb.org/mongo-driver/mongo"
)

func setupTestDB(t *testing.T) *mongo.Database {
    // Connect to test database
    client, _ := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
    db := client.Database("test_fmp_app")

    // Clean up after test
    t.Cleanup(func() {
        db.Drop(context.Background())
        client.Disconnect(context.Background())
    })

    return db
}

func TestAccountRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    repo := NewAccountRepository(db)

    account, err := repo.Create("user-id", CreateAccountRequest{
        Name: "Test Account",
        Type: "cash",
        Balance: 1000,
    })

    assert.NoError(t, err)
    assert.NotEmpty(t, account.ID)

    // Verify in database
    found, _ := repo.GetByID(account.ID, "user-id")
    assert.Equal(t, account.Name, found.Name)
}
```

#### Run Tests

```bash
# All tests
go test ./...

# Specific package
go test ./internal/services

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...
```

### Frontend Testing

#### Component Tests

```tsx
// src/components/TransactionList.test.tsx
import { render, screen } from "@testing-library/react";
import { TransactionList } from "./TransactionList";

describe("TransactionList", () => {
    const mockTransactions = [
        { id: "1", merchant: "Phở 24", amount: 50000, type: "expense" },
        { id: "2", merchant: "Salary", amount: 20000000, type: "income" },
    ];

    it("renders transaction list", () => {
        render(<TransactionList transactions={mockTransactions} />);

        expect(screen.getByText("Phở 24")).toBeInTheDocument();
        expect(screen.getByText("Salary")).toBeInTheDocument();
    });

    it("calls onSelect when transaction clicked", () => {
        const handleSelect = jest.fn();
        render(
            <TransactionList
                transactions={mockTransactions}
                onSelect={handleSelect}
            />,
        );

        screen.getByText("Phở 24").click();
        expect(handleSelect).toHaveBeenCalledWith("1");
    });
});
```

#### Service Tests (Mocking)

```tsx
// src/services/TransactionService.test.ts
import { TransactionService } from "./TransactionService";

// Mock fetch
global.fetch = jest.fn();

describe("TransactionService", () => {
    afterEach(() => {
        jest.clearAllMocks();
    });

    it("fetches transactions successfully", async () => {
        const mockData = [{ id: "1", amount: 1000 }](
            fetch as jest.Mock,
        ).mockResolvedValueOnce({
            ok: true,
            json: async () => mockData,
        });

        const result = await TransactionService.listTransactions();

        expect(result).toEqual(mockData);
        expect(fetch).toHaveBeenCalledWith("/api/v1/transactions");
    });
});
```

#### Run Tests

```bash
# All tests
npm run test

# Watch mode
npm run test:watch

# Coverage
npm run test:coverage
```

---

## Debugging

### Backend Debugging

#### VS Code Launch Configuration

Create `.vscode/launch.json`:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch Backend",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/server/cmd/api",
            "env": {
                "SERVER_PORT": "8080"
            },
            "args": []
        }
    ]
}
```

Set breakpoints in VS Code and press F5.

#### Manual Debugging

```go
// Add print statements
fmt.Printf("Debug: userID = %s\n", userID)

// Use logger
logger.Log.Info.Printf("Creating account for user: %s", userID)

// Check values at runtime
if account == nil {
    logger.Log.Error.Println("Account is nil!")
}
```

### Frontend Debugging

#### Browser DevTools

- Console: `console.log()`, `console.error()`
- Network tab: Inspect API requests/responses
- React DevTools: Inspect component state

#### VS Code Debugging

```tsx
// Add debugger statement
function handleSubmit() {
    debugger; // Execution will pause here
    TransactionService.createTransaction(data);
}
```

---

## Common Tasks

### Add New API Endpoint

**Backend:**

```bash
# 1. Add model/DTO
# 2. Add repository method
# 3. Add service method
# 4. Add handler
# 5. Register route in router.go

# Example: GET /api/v1/accounts/:id/balance
func (h *AccountHandler) GetBalance(c *gin.Context) {
    id := c.Param("id")
    userID, _ := c.Get("user_id")

    account, err := h.service.GetAccount(id, userID.(string))
    if err != nil {
        response.NotFoundResponse(c, "Account")
        return
    }

    response.SuccessResponse(c, 200, "Success", gin.H{
        "balance": account.Balance,
        "account_id": account.ID,
    })
}

# In router.go
accountGroup.GET("/:id/balance", accountHandler.GetBalance)
```

**Frontend:**

```tsx
// Add service method
async getAccountBalance(id: string): Promise<number> {
  const response = await fetch(`${this.baseURL}/accounts/${id}/balance`)
  const data = await response.json()
  return data.balance
}

// Use in component
const [balance, setBalance] = useState(0)

useEffect(() => {
  AccountService.getAccountBalance(accountId)
    .then(setBalance)
}, [accountId])
```

### Add New UI Component

```bash
# 1. Create component file
# src/components/shared/AccountBadge.tsx

export function AccountBadge({ account }: { account: Account }) {
  return (
    <div className="inline-flex items-center gap-2 px-3 py-1 bg-blue-100 rounded-full">
      <span>{account.icon}</span>
      <span className="font-medium">{account.name}</span>
    </div>
  )
}

# 2. Export from index
# src/components/shared/index.ts
export { AccountBadge } from "./AccountBadge"

# 3. Use in other components
import { AccountBadge } from "@/components/shared"

<AccountBadge account={selectedAccount} />
```

### Add Database Index

```javascript
// Connect to MongoDB
use fmp_app

// Create index
db.transactions.createIndex({
  user_id: 1,
  merchant: 1
})

// Verify index
db.transactions.getIndexes()

// Check query uses index (explain)
db.transactions.find({
  user_id: "user-id",
  merchant: "Phở 24"
}).explain("executionStats")
```

### Database Migration

```javascript
// Add new field to existing documents
db.accounts.updateMany(
    { is_active: { $exists: false } },
    { $set: { is_active: true } },
);

// Rename field
db.transactions.updateMany(
    {},
    { $rename: { old_field_name: "new_field_name" } },
);

// Remove field
db.accounts.updateMany({}, { $unset: { deprecated_field: "" } });
```

---

## Troubleshooting

### Backend Issues

#### "Cannot connect to MongoDB"

```bash
# Check connection string
echo $MONGODB_URI

# Test connection
mongosh "$MONGODB_URI"

# Check firewall/IP whitelist on MongoDB Atlas
```

#### "JWT validation failed"

```bash
# Check Supabase JWT secret
echo $SUPABASE_JWT_SECRET

# Verify token format
# Should be: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# Check token expiration
# Decode JWT at jwt.io and check "exp" claim
```

#### "Module not found"

```bash
# Sync dependencies
go mod tidy

# Verify go.mod
cat go.mod

# Clear cache if needed
go clean -modcache
```

### Frontend Issues

#### "API request failed (CORS error)"

```bash
# Check backend CORS configuration
# In .env:
ALLOWED_ORIGINS=http://localhost:5173

# Check browser console for exact error
```

#### "Module not found"

```bash
# Re-install dependencies
rm -rf node_modules package-lock.json
npm install

# Check import path
# Should use @ alias for src:
import { Button } from "@/components/ui/button"
```

#### "Build error"

```bash
# Clear cache
rm -rf node_modules/.vite

# Check TypeScript errors
npm run type-check

# Verbose build output
npm run build -- --debug
```

### Database Issues

#### "Slow queries"

```javascript
// Enable profiling
db.setProfilingLevel(2)

// Check slow queries
db.system.profile.find().sort({ ts: -1 }).limit(5)

// Analyze with explain
db.transactions.find({...}).explain("executionStats")

// Add missing indexes
```

#### "Duplicate key error"

```javascript
// Check unique indexes
db.accounts.getIndexes();

// Drop problematic index if needed
db.accounts.dropIndex("index_name");

// Re-create with correct uniqueness
```

---

## Useful Commands Cheatsheet

### Backend

```bash
# Development
go run cmd/api/main.go
air  # Hot reload

# Build
go build -o app cmd/api/main.go

# Test
go test ./...
go test -v ./internal/services
go test -cover ./...

# Dependencies
go mod tidy
go mod download
go get package@version

# Format
go fmt ./...
gofmt -w .

# Lint
golangci-lint run
```

### Frontend

```bash
# Development
npm run dev
npm run dev -- --host  # Expose to network

# Build
npm run build
npm run preview

# Test
npm run test
npm run test:watch
npm run test:coverage

# Lint & Format
npm run lint
npm run format

# Type check
npm run type-check
```

### Git

```bash
# Feature branch
git checkout -b feature/name
git add .
git commit -m "feat: description"
git push origin feature/name

# Update from main
git checkout main
git pull
git checkout feature/name
git rebase main

# Squash commits
git rebase -i HEAD~3
```

### MongoDB

```bash
# Connect
mongosh "$MONGODB_URI"

# Database commands
use fmp_app
show collections
db.transactions.find().limit(5)
db.transactions.countDocuments({ user_id: "xxx" })

# Indexes
db.transactions.getIndexes()
db.transactions.createIndex({ user_id: 1 })
db.transactions.dropIndex("index_name")
```

---

## Resources

### Documentation

- [Go Documentation](https://go.dev/doc/)
- [Gin Framework](https://gin-gonic.com/docs/)
- [MongoDB Go Driver](https://www.mongodb.com/docs/drivers/go/current/)
- [React Documentation](https://react.dev/)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/)
- [Tailwind CSS](https://tailwindcss.com/docs)

### Tools

- [MongoDB Compass](https://www.mongodb.com/products/compass)
- [Postman](https://www.postman.com/)
- [VS Code Go Extension](https://marketplace.visualstudio.com/items?itemName=golang.Go)

---

**Last Updated**: February 27, 2026
