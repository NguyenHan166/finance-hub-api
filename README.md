# Finance Hub - Personal Finance Management App

A comprehensive personal finance tracking application built with Go, MongoDB, and React. Track expenses, manage budgets, analyze spending patterns, and get AI-powered financial insights.

[🇻🇳 Đọc bằng tiếng Việt](./server/docs/README_vi.md)

---

## 🚀 Features

### Core Functionality

- **💰 Multi-Account Management**: Track cash, bank accounts, and credit cards
- **💸 Transaction Tracking**: Record income, expenses, and transfers with categorization
- **📁 Smart Categories**: Default categories + custom category support with icons
- **🎯 Budget Planning**: Set monthly budgets and track spending progress
- **📊 Advanced Reports**: Comprehensive analytics with charts and trends
- **🤖 AI Assistant**: Chat with AI for financial insights and recommendations
- **🔔 Smart Alerts**: Automatic notifications for budget limits and unusual spending
- **📈 Forecasting**: Predict future expenses based on historical data

### Additional Features

- VietQR integration for easy bank account setup
- Multi-currency support (VND, USD)
- Dark mode support
- Export reports to PDF/Excel
- Mobile-responsive design
- Real-time data synchronization

---

## 🏗️ Architecture

**Backend**: Clean Architecture with 4 layers

```
📦 Backend (Go)
├── handlers/      # HTTP request handlers (Gin)
├── services/      # Business logic layer
├── repositories/  # Data access layer (MongoDB)
└── models/        # Domain models
```

**Frontend**: React with TypeScript

```
📦 Frontend (React + TypeScript)
├── components/    # Reusable UI components
├── features/      # Feature-based modules
├── services/      # API client services
├── contexts/      # React Context (Auth, etc.)
└── pages/         # Route pages
```

**Database**: MongoDB Atlas

- 9 collections: users, accounts, transactions, categories, budgets, recurring_transactions, alerts, forecasts, chat_messages
- Indexes optimized for common queries

---

## 🛠️ Tech Stack

### Backend

- **Language**: Go 1.22+
- **Framework**: Gin v1.10.0
- **Database**: MongoDB driver v1.14.0
- **Auth**: Supabase JWT
- **Config**: godotenv

### Frontend

- **Framework**: React 18 + TypeScript
- **Build Tool**: Vite
- **UI Library**: Shadcn/UI (Radix + Tailwind CSS)
- **Charts**: Recharts
- **State Management**: Context API
- **Forms**: React Hook Form + Zod

### DevOps

- **Database**: MongoDB Atlas (Cloud)
- **Auth Provider**: Supabase
- **Hot Reload**: Air (Go) + Vite HMR (React)

---

## 📋 Prerequisites

- **Go**: 1.22 or higher ([Download](https://go.dev/dl/))
- **Node.js**: 18+ and npm ([Download](https://nodejs.org/))
- **MongoDB**: Atlas account or local MongoDB instance
- **Supabase**: Account for authentication ([Sign up](https://supabase.com/))

---

## ⚡ Quick Start

### 1. Clone Repository

```bash
git clone <repository-url>
cd FinanceTracking
```

### 2. Backend Setup

```bash
cd server

# Install dependencies
go mod download

# Create .env file
cp .env.example .env

# Edit .env with your credentials
# MONGODB_URI=mongodb+srv://...
# SUPABASE_URL=...
# SUPABASE_KEY=...

# Run backend
go run main.go

# Or with hot reload (install Air first)
air
```

Backend runs on http://localhost:8080

### 3. Frontend Setup

```bash
cd my-finance-hub

# Install dependencies
npm install

# Run development server
npm run dev
```

Frontend runs on http://localhost:5173

### 4. Database Setup

See [QUICKSTART.md](./server/docs/QUICKSTART.md) for detailed step-by-step setup including:

- Creating MongoDB collections
- Setting up indexes
- Seeding default categories

---

## 📚 Documentation

Comprehensive documentation is available in the `server/docs/` folder:

| Document                                                             | Description                                       |
| -------------------------------------------------------------------- | ------------------------------------------------- |
| [QUICKSTART.md](./server/docs/QUICKSTART.md)                         | **Start here!** Step-by-step setup guide (20 min) |
| [API_DOCUMENTATION.md](./server/docs/API_DOCUMENTATION.md)           | Complete REST API reference with all endpoints    |
| [DATABASE_SCHEMA.md](./server/docs/DATABASE_SCHEMA.md)               | MongoDB schema for all 9 collections              |
| [ARCHITECTURE.md](./server/docs/ARCHITECTURE.md)                     | System architecture and design patterns           |
| [DEVELOPMENT_GUIDE.md](./server/docs/DEVELOPMENT_GUIDE.md)           | Developer workflow and best practices             |
| [FEATURE_SPECIFICATIONS.md](./server/docs/FEATURE_SPECIFICATIONS.md) | Detailed feature specs for all modules            |

---

## 🔑 Environment Variables

### Backend (.env)

```env
# MongoDB
MONGODB_URI=mongodb+srv://username:password@cluster.mongodb.net/fmp_app

# Supabase Auth
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_KEY=your-anon-key

# Server
PORT=8080
GIN_MODE=debug

# CORS (optional)
ALLOWED_ORIGINS=http://localhost:5173
```

### Frontend (.env)

```env
VITE_SUPABASE_URL=https://your-project.supabase.co
VITE_SUPABASE_ANON_KEY=your-anon-key
VITE_API_URL=http://localhost:8080/api/v1
```

---

## 🧪 Testing

### Backend Tests

```bash
cd server

# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./services
```

### Frontend Tests

```bash
cd my-finance-hub

# Run Vitest tests
npm run test

# Run with coverage
npm run test:coverage

# Run in watch mode
npm run test:watch
```

---

## 🏃 Common Tasks

### Add a New API Endpoint

1. Define model in `models/`
2. Create repository methods in `repositories/`
3. Implement business logic in `services/`
4. Add handler in `handlers/`
5. Register route in `main.go`

See [DEVELOPMENT_GUIDE.md](./server/docs/DEVELOPMENT_GUIDE.md#adding-a-new-endpoint) for detailed steps.

### Add a New UI Component

```bash
cd my-finance-hub

# Use shadcn CLI to add component
npx shadcn-ui@latest add <component-name>

# Or create custom component in src/components/
```

### Create Database Migration

MongoDB is schemaless, but for index migrations:

```javascript
// In mongosh
db.transactions.createIndex({ userId: 1, dateTimeISO: -1 });
```

---

## 📈 Project Structure

```
FinanceTracking/
├── server/                    # Backend (Go)
│   ├── handlers/             # HTTP handlers
│   ├── services/             # Business logic
│   ├── repositories/         # Data access
│   ├── models/               # Domain models
│   ├── middleware/           # Auth, CORS, etc.
│   ├── docs/                 # Documentation
│   ├── main.go               # Entry point
│   ├── go.mod                # Dependencies
│   └── .env                  # Config (not committed)
│
├── my-finance-hub/           # Frontend (React)
│   ├── src/
│   │   ├── components/       # UI components
│   │   ├── features/         # Feature modules
│   │   ├── services/         # API clients
│   │   ├── contexts/         # React contexts
│   │   ├── pages/            # Route pages
│   │   ├── hooks/            # Custom hooks
│   │   ├── lib/              # Utilities
│   │   └── models/           # TypeScript types
│   ├── public/               # Static assets
│   ├── package.json          # Dependencies
│   └── vite.config.ts        # Vite config
│
└── README.md                 # This file
```

---

## 🐛 Troubleshooting

### Backend won't start

- **CORS errors**: Check `ALLOWED_ORIGINS` in `.env`
- **MongoDB connection failed**: Verify `MONGODB_URI` and network access in MongoDB Atlas
- **Port already in use**: Change `PORT` in `.env` or kill process on port 8080

### Frontend build fails

- **Module not found**: Run `npm install` again
- **Vite errors**: Clear cache: `rm -rf node_modules/.vite && npm run dev`
- **API calls fail**: Check `VITE_API_URL` points to correct backend

### Authentication issues

- **JWT invalid**: Verify `SUPABASE_KEY` matches in backend and frontend `.env`
- **User not found**: Check Supabase Users table
- **Token expired**: Token expires after 1 hour, re-login required

See [DEVELOPMENT_GUIDE.md](./server/docs/DEVELOPMENT_GUIDE.md#troubleshooting) for more solutions.

---

## 🚀 Deployment

### Backend (Go)

```bash
# Build binary
go build -o finance-hub-api main.go

# Run in production
GIN_MODE=release ./finance-hub-api
```

Deploy to:

- **Railway**, **Render**, **Fly.io**: Docker or binary
- **AWS EC2**: Upload binary + systemd service
- **Google Cloud Run**: Containerized deployment

### Frontend (React)

```bash
# Build for production
npm run build

# Output in dist/
# Deploy dist/ folder to:
```

Deploy to:

- **Vercel**, **Netlify**: Auto-deploy from Git
- **AWS S3 + CloudFront**: Static hosting
- **nginx**: Serve `dist/` folder

### Database

- **MongoDB Atlas**: Already cloud-hosted
- **Backups**: Enable automated backups in Atlas

---

## 🤝 Contributing

1. Fork the repository
2. Create feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Open Pull Request

### Code Style

- **Go**: Follow standard Go conventions, use `gofmt`
- **TypeScript/React**: ESLint + Prettier configured
- **Commits**: Use conventional commits (feat:, fix:, docs:, etc.)

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 👥 Authors

- **Development Team** - Finance Hub

---

## 🙏 Acknowledgments

- [Gin Web Framework](https://gin-gonic.com/)
- [MongoDB Go Driver](https://www.mongodb.com/docs/drivers/go/current/)
- [React](https://react.dev/)
- [Shadcn/UI](https://ui.shadcn.com/)
- [Supabase](https://supabase.com/)
- [VietQR](https://vietqr.io/)

---

## 📞 Support

- **Documentation**: [server/docs/](./server/docs/)
- **Issues**: [GitHub Issues](<repository-url>/issues)
- **Email**: support@financehub.com

---

## 🗺️ Roadmap

### Current Version: 1.0 (MVP)

- ✅ Account management
- ✅ Transaction tracking
- ✅ Categories
- ✅ Budgets
- ✅ Basic reports

### Version 1.1 (Q2 2026)

- 🔄 AI Chat Assistant
- 🔄 Advanced forecasting
- 🔄 Recurring transactions
- 🔄 Mobile app (React Native)

### Version 1.2 (Q3 2026)

- 📅 Bank integration (auto-import transactions)
- 📅 Multi-user accounts (family sharing)
- 📅 Investment tracking
- 📅 Bill reminders

### Future

- Split expenses (group payments)
- Cryptocurrency tracking
- Tax report generation
- Financial goal tracking

---

## 📊 Project Stats

![GitHub stars](https://img.shields.io/github/stars/yourusername/finance-hub)
![GitHub forks](https://img.shields.io/github/forks/yourusername/finance-hub)
![GitHub issues](https://img.shields.io/github/issues/yourusername/finance-hub)
![License](https://img.shields.io/github/license/yourusername/finance-hub)

---

**Built with ❤️ by Finance Hub Team**
