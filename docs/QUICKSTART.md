# Quick Start Guide - Finance Hub

This guide will walk you through setting up Finance Hub from scratch and making your first transactions. Perfect for new developers or users wanting to test the application quickly.

---

## ⏱️ Time Required

- **Backend Setup**: 10 minutes
- **Frontend Setup**: 5 minutes
- **Database Setup**: 5 minutes
- **First Transaction**: 2 minutes

**Total**: ~20-25 minutes

---

## 📋 Prerequisites Checklist

Before starting, ensure you have:

- [ ] Go 1.22+ installed ([Download](https://go.dev/dl/))
- [ ] Node.js 18+ and npm installed ([Download](https://nodejs.org/))
- [ ] Git installed
- [ ] Code editor (VS Code recommended)
- [ ] MongoDB Atlas account (free tier) ([Sign up](https://www.mongodb.com/cloud/atlas/register))
- [ ] Supabase account (free tier) ([Sign up](https://supabase.com/))

---

## 🚀 Step-by-Step Setup

### Step 1: Clone the Repository (2 min)

```bash
# Clone the repository
git clone <repository-url>
cd FinanceTracking

# Verify structure
ls
# You should see: server/ my-finance-hub/ README.md
```

---

### Step 2: Setup MongoDB Atlas (5 min)

#### 2.1 Create MongoDB Cluster

1. Go to [MongoDB Atlas](https://cloud.mongodb.com/)
2. Sign in / Sign up
3. Click **"Build a Database"**
4. Choose **"M0 Free"** tier
5. Select a cloud provider and region (closest to you)
6. Name your cluster: `Cluster0`
7. Click **"Create"**

#### 2.2 Setup Database Access

1. In the left sidebar, click **"Database Access"**
2. Click **"Add New Database User"**
3. Choose **"Password"** authentication
4. Username: `nvhan166` (or your choice)
5. Password: `han1662003` (or generate secure password)
6. Database User Privileges: **"Read and write to any database"**
7. Click **"Add User"**

#### 2.3 Setup Network Access

1. In the left sidebar, click **"Network Access"**
2. Click **"Add IP Address"**
3. Click **"Allow Access from Anywhere"** (for development)
    - IP Address: `0.0.0.0/0`
4. Click **"Confirm"**

#### 2.4 Get Connection String

1. Go back to **"Database"** in left sidebar
2. Click **"Connect"** on your cluster
3. Choose **"Connect your application"**
4. Driver: **"Go"**, Version: **"1.13 or later"**
5. Copy the connection string:
    ```
    mongodb+srv://nvhan166:<password>@cluster0.evbdltl.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0
    ```
6. Replace `<password>` with your actual password
7. Add database name: `.../fmp_app?retryWrites=...`

**Final URI:**

```
mongodb+srv://nvhan166:han1662003@cluster0.evbdltl.mongodb.net/fmp_app?retryWrites=true&w=majority&appName=Cluster0
```

---

### Step 3: Setup Supabase Authentication (5 min)

#### 3.1 Create Supabase Project

1. Go to [Supabase Dashboard](https://supabase.com/dashboard)
2. Sign in / Sign up
3. Click **"New Project"**
4. Organization: Create new or select existing
5. Project Name: `finance-hub`
6. Database Password: Generate strong password (save it!)
7. Region: Choose closest to you
8. Click **"Create new project"** (wait ~2 minutes)

#### 3.2 Get Supabase Credentials

1. Once project is ready, go to **Settings** (left sidebar)
2. Click **"API"**
3. Copy these values:
    - **Project URL**: `https://xxxxx.supabase.co`
    - **anon public key**: `eyJhbGc...` (long JWT token)

#### 3.3 Create Test User (Optional)

1. Go to **Authentication** → **Users** in left sidebar
2. Click **"Add User"**
3. Email: `test@example.com`
4. Password: `Test1234!`
5. Click **"Create user"**

---

### Step 4: Backend Setup (5 min)

```bash
# Navigate to backend folder
cd server

# Install Go dependencies
go mod download

# This will download:
# - gin-gonic/gin v1.10.0
# - mongodb/mongo-go-driver v1.14.0
# - joho/godotenv
```

#### 4.1 Create Environment File

```bash
# Create .env file
touch .env

# Or on Windows:
# type nul > .env
```

#### 4.2 Edit .env File

Open `server/.env` in your editor and add:

```env
# MongoDB
MONGODB_URI=mongodb+srv://nvhan166:han1662003@cluster0.evbdltl.mongodb.net/fmp_app?retryWrites=true&w=majority&appName=Cluster0

# Supabase
SUPABASE_URL=https://xxxxx.supabase.co
SUPABASE_KEY=eyJhbGc...your-anon-key...

# Server
PORT=8080
GIN_MODE=debug

# CORS
ALLOWED_ORIGINS=http://localhost:5173
```

**Replace:**

- MongoDB URI with your actual connection string
- `SUPABASE_URL` with your Project URL
- `SUPABASE_KEY` with your anon public key

#### 4.3 Run Backend

```bash
# Run backend
go run main.go

# You should see:
# [GIN-debug] Listening on :8080
```

**Test it works:**
Open http://localhost:8080/health in browser

Expected response:

```json
{
    "status": "ok",
    "timestamp": "2026-02-27T10:30:00Z"
}
```

---

### Step 5: Frontend Setup (5 min)

Open a **new terminal** (keep backend running):

```bash
# Navigate to frontend folder
cd my-finance-hub

# Install dependencies
npm install

# This will take 1-2 minutes
```

#### 5.1 Create Frontend .env

```bash
# Create .env file
touch .env

# Or on Windows:
# type nul > .env
```

#### 5.2 Edit .env File

Open `my-finance-hub/.env` in your editor and add:

```env
VITE_SUPABASE_URL=https://xxxxx.supabase.co
VITE_SUPABASE_ANON_KEY=eyJhbGc...your-anon-key...
VITE_API_URL=http://localhost:8080/api/v1
```

**Use the same Supabase values from Step 3.**

#### 5.3 Run Frontend

```bash
# Run dev server
npm run dev

# You should see:
# VITE ready in 500ms
# ➜ Local: http://localhost:5173/
```

**Open browser**: http://localhost:5173

You should see the Finance Hub login page! 🎉

---

### Step 6: Initialize Database (5 min)

The database needs default categories and indexes.

#### 6.1 Connect to MongoDB

Open [MongoDB Compass](https://www.mongodb.com/products/compass) (GUI) or use `mongosh` (CLI):

**Using mongosh:**

```bash
# Install mongosh if not installed
# https://www.mongodb.com/docs/mongodb-shell/install/

# Connect
mongosh "mongodb+srv://nvhan166:han1662003@cluster0.evbdltl.mongodb.net/fmp_app"
```

#### 6.2 Create Indexes

Copy and paste this in `mongosh`:

```javascript
// Users collection indexes
db.users.createIndex({ email: 1 }, { unique: true });

// Accounts collection indexes
db.accounts.createIndex({ userId: 1 });
db.accounts.createIndex({ userId: 1, isActive: 1 });

// Transactions collection indexes
db.transactions.createIndex({ userId: 1, dateTimeISO: -1 });
db.transactions.createIndex({ userId: 1, accountId: 1 });
db.transactions.createIndex({ userId: 1, categoryId: 1 });
db.transactions.createIndex({ userId: 1, type: 1, dateTimeISO: -1 });

// Categories collection indexes
db.categories.createIndex({ userId: 1 });
db.categories.createIndex({ userId: 1, isDefault: 1 });

// Budgets collection indexes
db.budgets.createIndex({ userId: 1, month: 1 });
db.budgets.createIndex(
    { userId: 1, categoryId: 1, month: 1 },
    { unique: true },
);

// Alerts collection indexes
db.alerts.createIndex({ userId: 1, status: 1, createdAt: -1 });

// Chat messages collection indexes
db.chat_messages.createIndex({ userId: 1, createdAt: -1 });

print("✅ All indexes created!");
```

#### 6.3 Seed Default Categories

```javascript
// Income categories
db.categories.insertMany([
    {
        _id: "cat_income_salary",
        userId: "default",
        name: "Lương",
        type: "income",
        icon: "💰",
        color: "#10B981",
        isDefault: true,
        createdAt: new Date(),
        updatedAt: new Date(),
    },
    {
        _id: "cat_income_bonus",
        userId: "default",
        name: "Thưởng",
        type: "income",
        icon: "🎁",
        color: "#3B82F6",
        isDefault: true,
        createdAt: new Date(),
        updatedAt: new Date(),
    },
    {
        _id: "cat_income_investment",
        userId: "default",
        name: "Đầu tư",
        type: "income",
        icon: "📈",
        color: "#8B5CF6",
        isDefault: true,
        createdAt: new Date(),
        updatedAt: new Date(),
    },
    {
        _id: "cat_income_other",
        userId: "default",
        name: "Thu nhập khác",
        type: "income",
        icon: "💵",
        color: "#6B7280",
        isDefault: true,
        createdAt: new Date(),
        updatedAt: new Date(),
    },
]);

// Expense categories
db.categories.insertMany([
    {
        _id: "cat_expense_food",
        userId: "default",
        name: "Ăn uống",
        type: "expense",
        icon: "🍜",
        color: "#EF4444",
        isDefault: true,
        createdAt: new Date(),
        updatedAt: new Date(),
    },
    {
        _id: "cat_expense_transport",
        userId: "default",
        name: "Di chuyển",
        type: "expense",
        icon: "🚗",
        color: "#F59E0B",
        isDefault: true,
        createdAt: new Date(),
        updatedAt: new Date(),
    },
    {
        _id: "cat_expense_shopping",
        userId: "default",
        name: "Mua sắm",
        type: "expense",
        icon: "🛍️",
        color: "#EC4899",
        isDefault: true,
        createdAt: new Date(),
        updatedAt: new Date(),
    },
    {
        _id: "cat_expense_house",
        userId: "default",
        name: "Nhà cửa",
        type: "expense",
        icon: "🏠",
        color: "#8B5CF6",
        isDefault: true,
        createdAt: new Date(),
        updatedAt: new Date(),
    },
    {
        _id: "cat_expense_entertainment",
        userId: "default",
        name: "Giải trí",
        type: "expense",
        icon: "🎮",
        color: "#3B82F6",
        isDefault: true,
        createdAt: new Date(),
        updatedAt: new Date(),
    },
    {
        _id: "cat_expense_health",
        userId: "default",
        name: "Sức khỏe",
        type: "expense",
        icon: "⚕️",
        color: "#10B981",
        isDefault: true,
        createdAt: new Date(),
        updatedAt: new Date(),
    },
    {
        _id: "cat_expense_education",
        userId: "default",
        name: "Giáo dục",
        type: "expense",
        icon: "📚",
        color: "#6366F1",
        isDefault: true,
        createdAt: new Date(),
        updatedAt: new Date(),
    },
    {
        _id: "cat_expense_other",
        userId: "default",
        name: "Chi phí khác",
        type: "expense",
        icon: "💸",
        color: "#6B7280",
        isDefault: true,
        createdAt: new Date(),
        updatedAt: new Date(),
    },
]);

print("✅ Default categories created!");
```

---

### Step 7: Create Your First Account & Transaction (2 min)

#### 7.1 Register / Login

1. Open http://localhost:5173 in browser
2. Click **"Sign Up"** (or use test user from Step 3.3)
3. Enter email and password
4. Click **"Create Account"**
5. You'll be redirected to Dashboard

#### 7.2 Create Your First Account

1. In the sidebar, click **"Accounts"**
2. Click **"Add Account"** button (top right)
3. Choose **"Cash"**
4. Fill in:
    - Name: `Ví tiền mặt`
    - Initial Balance: `5000000` (5 triệu)
    - Icon: 💵 (or your choice)
    - Color: Green
5. Click **"Save"**

#### 7.3 Add Your First Transaction

1. In the sidebar, click **"Transactions"**
2. Click **"Add Transaction"** button
3. Choose **"Expense"** tab
4. Fill in:
    - Amount: `150000` (150k)
    - Account: Select `Ví tiền mặt`
    - Category: Select `🍜 Ăn uống`
    - Merchant: `Phở 24`
    - Date & Time: Now (default)
    - Note: `Ăn trưa`
5. Click **"Save"**

#### 7.4 View Dashboard

1. Click **"Dashboard"** in sidebar
2. You should see:
    - Total Expense: 150,000 ₫
    - Account balance: 4,850,000 ₫ (5M - 150k)
    - Spend Trend chart
    - Top Categories (Ăn uống)

**🎉 Congratulations! You've successfully set up Finance Hub!**

---

## ✅ Verification Checklist

After setup, verify everything works:

- [ ] Backend running on http://localhost:8080
- [ ] Frontend running on http://localhost:5173
- [ ] Can register/login successfully
- [ ] Can create account
- [ ] Can add transaction
- [ ] Dashboard shows data correctly
- [ ] No console errors in browser DevTools
- [ ] No errors in backend terminal

---

## 🎯 Next Steps

Now that you're set up, try these:

### Explore Features

1. **Create a Budget**: Go to Budgets → Add budget for "Ăn uống" (2M/month)
2. **Add More Transactions**: Add 5-10 transactions across different categories
3. **View Reports**: Check Reports page for analytics
4. **Try AI Chat**: Ask "Tháng này tôi chi bao nhiêu?"

### Development Tasks

1. **Add a Bank Account**: Use VietQR integration
2. **Create Custom Category**: Settings → Categories → Add custom category
3. **Set Up Alerts**: Add more expenses to trigger budget alert
4. **Export Report**: Try exporting to PDF/Excel

### Learn the Codebase

1. Read [ARCHITECTURE.md](./ARCHITECTURE.md) - Understand system design
2. Read [API_DOCUMENTATION.md](./API_DOCUMENTATION.md) - Learn API endpoints
3. Read [DEVELOPMENT_GUIDE.md](./DEVELOPMENT_GUIDE.md) - Development workflow

---

## 🐛 Troubleshooting

### Backend Issues

**Problem**: `mongo: no reachable servers`

```
Solution:
1. Check MongoDB URI in .env is correct
2. Verify network access (0.0.0.0/0) in MongoDB Atlas
3. Check username/password are correct
```

**Problem**: `Supabase: Invalid JWT`

```
Solution:
1. Verify SUPABASE_KEY is the "anon public" key (not service_role key)
2. Check SUPABASE_URL is correct
3. Restart backend after .env changes
```

**Problem**: `Port 8080 already in use`

```
Solution:
# Find and kill process on port 8080
# macOS/Linux:
lsof -ti:8080 | xargs kill -9

# Windows:
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Or change PORT in .env to 8081
```

### Frontend Issues

**Problem**: `Network Error` when calling API

```
Solution:
1. Check backend is running (http://localhost:8080/health)
2. Check VITE_API_URL in frontend .env
3. Check CORS: ALLOWED_ORIGINS in backend .env includes http://localhost:5173
4. Clear browser cache and reload
```

**Problem**: `Module not found` errors

```
Solution:
# Delete node_modules and reinstall
rm -rf node_modules package-lock.json
npm install
```

**Problem**: Supabase auth not working

```
Solution:
1. Check VITE_SUPABASE_URL and VITE_SUPABASE_ANON_KEY are correct
2. Restart Vite dev server after .env changes:
   Ctrl+C then npm run dev
3. Clear browser localStorage:
   F12 → Application → Local Storage → Clear
```

### Database Issues

**Problem**: Categories not showing

```
Solution:
1. Connect to MongoDB via mongosh or Compass
2. Check categories collection exists and has default categories:
   db.categories.find({ isDefault: true }).count()
3. If 0, re-run seed script from Step 6.3
```

**Problem**: Slow queries

```
Solution:
1. Verify indexes are created:
   db.transactions.getIndexes()
2. If missing, re-run index creation from Step 6.2
```

---

## 🎓 Learning Resources

### Backend (Go)

- [Go by Example](https://gobyexample.com/)
- [Gin Framework Docs](https://gin-gonic.com/docs/)
- [MongoDB Go Driver Tutorial](https://www.mongodb.com/docs/drivers/go/current/)

### Frontend (React)

- [React Docs](https://react.dev/)
- [TypeScript Handbook](https://www.typescriptlang.org/docs/handbook/intro.html)
- [Shadcn/UI Components](https://ui.shadcn.com/)

### Database

- [MongoDB University](https://university.mongodb.com/) - Free courses
- [MongoDB Schema Design Best Practices](https://www.mongodb.com/developer/products/mongodb/schema-design-anti-pattern-summary/)

---

## 💡 Tips for Development

### Backend Development

```bash
# Use Air for hot reload
go install github.com/cosmtrek/air@latest
air

# Format code
go fmt ./...

# Run linter
golangci-lint run
```

### Frontend Development

```bash
# Type checking
npm run type-check

# Lint and fix
npm run lint
npm run lint:fix

# Format with Prettier
npm run format
```

### Database GUI

Use **MongoDB Compass** for visual database management:

1. Download: https://www.mongodb.com/products/compass
2. Connect using your MongoDB URI
3. Browse collections, run queries, create indexes visually

---

## 📞 Getting Help

If you're stuck:

1. **Check Documentation**:
    - [DEVELOPMENT_GUIDE.md](./DEVELOPMENT_GUIDE.md) - Most common issues
    - [API_DOCUMENTATION.md](./API_DOCUMENTATION.md) - API reference
2. **Check Logs**:
    - Backend: Terminal where `go run main.go` is running
    - Frontend: Browser DevTools Console (F12)
    - MongoDB: Atlas Monitoring tab

3. **Search Issues**: Check if someone else had the same problem

4. **Ask for Help**: Open a GitHub issue with:
    - What you tried to do
    - What happened instead
    - Error messages (full stack trace)
    - Your environment (OS, Go version, Node version)

---

## 🎉 Success!

You now have a fully functional personal finance app running locally!

**What you've accomplished:**

- ✅ Set up MongoDB Atlas database
- ✅ Configured Supabase authentication
- ✅ Running Go backend with Gin framework
- ✅ Running React frontend with TypeScript
- ✅ Created your first account and transaction
- ✅ Dashboard displaying real-time data

**Ready to build?** Head to [DEVELOPMENT_GUIDE.md](./DEVELOPMENT_GUIDE.md) to learn how to add new features!

---

**Happy Coding! 🚀**
