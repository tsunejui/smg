# Social Media Growth Engine (SMG)

A comprehensive full-stack application for managing social media content growth with a web-based admin panel and a mobile app.

## 🏗️ Architecture

### Components:
- **Web Admin Panel**: Next.js with Tailwind CSS
- **Mobile App**: Flutter (iOS + Android)
- **Backend API**: Go with Gin framework
- **Database**: PostgreSQL with Prisma ORM
- **Cache**: Redis
- **Background Jobs**: Go scheduler service

### Repository Structure:
```
root/
├── apps/
│   ├── web/               # Next.js Admin Panel
│   └── app/               # Flutter Mobile App
├── cmd/                   # Go CLI tools and services
│   ├── api/               # REST API server
│   ├── scheduler/         # Background job scheduler
│   └── migrate/           # Database migration tool
├── pkg/                   # Shared Go modules
│   ├── models/            # Data models
│   ├── services/          # Business logic
│   ├── handlers/          # HTTP handlers
│   ├── middleware/        # HTTP middleware
│   └── config/            # Configuration
├── scripts/               # Setup scripts
├── Makefile              # Build and dev commands
├── docker-compose.yml    # Docker services
└── docs/                 # Documentation
```

## 🚀 Quick Start

### Prerequisites:
- Docker and Docker Compose
- Go 1.21+
- Node.js 18+
- Flutter 3.0+

### 1. Clone and Setup:
```bash
git clone <repository-url>
cd smg
make env          # Create .env file
make setup        # Install dependencies
```

### 2. Start Services:
```bash
make start        # Start all services with Docker
```

### 3. Run Database Migrations:
```bash
make migrate      # Run Prisma migrations
```

### 4. Access the Application:
- Web Admin: http://localhost:3000
- API Server: http://localhost:8080
- Database Studio: `make studio`

## 📱 Features

### Web Admin Panel:
- **Authentication**: Google OAuth2 and Email/Password
- **User Dashboard** (中文界面):
  - 主題管理 (Topic Management)
  - 媒體連結 (Media Connections)
  - 文章管理 (Content Management)
  - 系統設定 (System Settings)
- **Admin Backend**:
  - 用戶管理 (User Management)
  - 社群媒體設定 (Platform Settings)
  - 系統設定 (System Configuration)

### Mobile App:
- QR Code login
- Topic-based content browsing
- Content approval and reposting
- Offline capability

### Backend Services:
- RESTful API with JWT authentication
- Background job scheduling
- Real-time data synchronization
- Social media platform integration

## 🛠️ Development

### Available Make Commands:
```bash
# Development
make web          # Start Next.js dev server
make api          # Start Go API server
make scheduler    # Start background scheduler
make run-flutter  # Start Flutter app

# Database
make migrate      # Run Prisma migrations
make migrate-go   # Run Go migrations
make studio       # Open Prisma Studio
make reset        # Reset database

# Building
make build        # Build web app
make build-go     # Build Go services
make build-flutter # Build Flutter APK

# Testing
make test-go      # Run Go tests
make test-flutter # Run Flutter tests

# Utilities
make format       # Format all code
make clean        # Clean build artifacts
make logs         # View Docker logs
```

### Environment Variables:
Copy `.env.example` to `.env` and update the values:
- Database connection strings
- JWT secrets
- Social media API keys
- SMTP settings for email notifications

## 📚 Documentation

For detailed documentation, refer to the `/docs` directory:
- [Google OAuth Setup](docs/google-auth.md)
- [PostgreSQL Guide](docs/postgre-guideline.md)
- [Go Development Guide](docs/golang-guideline.md)
- [Next.js Guide](docs/nextjs-guideline.md)
- [Development Environment Setup](docs/dev-env.md)

## 🔧 Development

### Code Style
- **Go**: Follow Uber Go Style Guide
- **TypeScript**: Use ESLint and Prettier
- **Flutter**: Use `flutter format`

### Recommended IDE Setup
- VS Code with Prisma, ESLint, Prettier extensions
- GoLand or VS Code with Go extension
- Android Studio or VS Code with Flutter extension

## 📄 License

This project is licensed under the MIT License.