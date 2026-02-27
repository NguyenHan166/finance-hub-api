# 🚀 Quick Start Guide

## Setup trong 5 phút

### 1. Cài đặt Dependencies

```bash
cd server
go mod download
```

### 2. Cấu hình Environment

```bash
cp .env.example .env
```

Chỉnh sửa `.env` với thông tin Supabase của bạn:

```env
# Database từ Supabase
DB_HOST=db.xxx.supabase.co
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=postgres
DB_SSLMODE=require

# Supabase Keys
SUPABASE_URL=https://xxx.supabase.co
SUPABASE_ANON_KEY=your-anon-key
SUPABASE_SERVICE_ROLE_KEY=your-service-role-key

# JWT Secret
JWT_SECRET=your-super-secret-key
```

### 3. Chạy Database Migrations

Trong Supabase SQL Editor, chạy các file:

1. `migrations/001_create_tables.sql`
2. `migrations/002_seed_default_categories.sql`

### 4. Cài đặt Development Tools (Optional)

```bash
make install-tools
```

Hoặc manual:

```bash
go install github.com/cosmtrek/air@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### 5. Chạy Server

**Option A: Hot Reload (Recommended)**

```bash
make dev
```

**Option B: Run Direct**

```bash
go run cmd/api/main.go
```

Server sẽ chạy tại: `http://localhost:8080`

### 6. Test API

```bash
# Health check
curl http://localhost:8080/health

# Hoặc mở browser
http://localhost:8080/health
```

## 📁 Cấu trúc Project

```
server/
├── cmd/api/              # Entry point
├── internal/             # Private code
│   ├── config/          # Configuration
│   ├── handlers/        # HTTP handlers
│   ├── services/        # Business logic
│   ├── repositories/    # Database access
│   ├── models/          # Data models
│   ├── middleware/      # HTTP middleware
│   └── utils/           # Utilities
├── pkg/                 # Public libraries
├── migrations/          # SQL migrations
└── docs/                # Documentation
```

## 🔧 Lệnh Make hữu ích

```bash
make dev              # Run với hot reload
make build            # Build binary
make test             # Run tests
make test-coverage    # Test với coverage report
make fmt              # Format code
make lint             # Run linter
make clean            # Xóa build files
```

## 🐳 Docker Development

Nếu muốn chạy toàn bộ với Docker:

```bash
# Start tất cả services (PostgreSQL + API)
docker-compose up

# Chạy background
docker-compose up -d

# Xem logs
docker-compose logs -f api

# Stop
docker-compose down
```

## 📝 API Endpoints

### Authentication Required (JWT Token)

**Accounts:**

- `POST /api/v1/accounts` - Tạo account mới
- `GET /api/v1/accounts` - Lấy danh sách accounts
- `GET /api/v1/accounts/:id` - Lấy account theo ID
- `PUT /api/v1/accounts/:id` - Cập nhật account
- `DELETE /api/v1/accounts/:id` - Xóa account

**Transactions:**

- `POST /api/v1/transactions` - Tạo transaction
- `GET /api/v1/transactions` - Lấy danh sách transactions
- `GET /api/v1/transactions/:id` - Lấy transaction theo ID
- `DELETE /api/v1/transactions/:id` - Xóa transaction

**Categories:**

- `POST /api/v1/categories` - Tạo category
- `GET /api/v1/categories` - Lấy danh sách categories
- `GET /api/v1/categories/:id` - Lấy category theo ID
- `DELETE /api/v1/categories/:id` - Xóa category

Xem chi tiết tại: [docs/API.md](docs/API.md)

## 🧪 Test với Postman/cURL

### 1. Get JWT Token từ Supabase

Login từ frontend hoặc Supabase Dashboard để lấy access token.

### 2. Create Account

```bash
curl -X POST http://localhost:8080/api/v1/accounts \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Main Checking",
    "type": "checking",
    "balance": 1000.00,
    "currency": "USD",
    "bank_name": "Chase"
  }'
```

### 3. Get All Accounts

```bash
curl http://localhost:8080/api/v1/accounts \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🎯 Next Steps

1. **Integrate với Frontend:**
    - Update frontend service URLs
    - Configure CORS trong `.env`
    - Implement JWT authentication

2. **Thêm Features:**
    - Reports/Analytics
    - Budget tracking
    - Recurring transactions
    - File uploads (receipts)

3. **Production Deployment:**
    - Setup CI/CD
    - Configure production environment
    - Setup monitoring
    - Add rate limiting

## 📚 Documentation

- [API Documentation](docs/API.md)
- [Development Guide](docs/DEVELOPMENT.md)
- [README](README.md)

## ❓ Troubleshooting

**Port 8080 đã được sử dụng:**

```bash
# Windows
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Linux/Mac
lsof -ti:8080 | xargs kill -9
```

**Database connection error:**

- Kiểm tra Supabase credentials trong `.env`
- Verify database đang chạy
- Check firewall settings

**Hot reload không hoạt động:**

```bash
go install github.com/cosmtrek/air@latest
```

## 🎉 Hoàn tất!

Server đã sẵn sàng. Happy coding! 🚀

Nếu có vấn đề, check:

1. Logs trong terminal
2. [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md)
3. [docs/API.md](docs/API.md)
