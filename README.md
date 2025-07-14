# Social Media Growth Engine

A full-stack application integrating a web-based admin panel and mobile app to help users collect and repost content related to their favorite topics from various social media platforms, thereby growing their social media following.

## 🤖 AI-Generated Code

**This entire codebase was created using AI to test the feasibility of "vibe coding" - a development approach where code is generated through natural language descriptions and iterative refinement with AI assistance.**

This project serves as a proof of concept for AI-powered software development, demonstrating how complex full-stack applications can be built entirely through AI code generation. The implementation includes:

- Complete application architecture and setup
- Database schema design and migrations
- Authentication system with multiple providers
- Responsive UI with Chinese localization
- Docker containerization and development environment
- Full test suite and deployment configuration

All code, documentation, and configuration files in this repository were generated through AI assistance, showcasing the current capabilities and potential of AI in software development.

## 🏗️ Architecture

### Technology Stack
- **Frontend**: Next.js + TypeScript + Tailwind CSS
- **Backend**: Go + PostgreSQL
- **Mobile App**: Flutter (iOS + Android)
- **Database**: PostgreSQL with Prisma ORM
- **Authentication**: NextAuth.js (Google OAuth + Email/Password)
- **Development Environment**: Docker Compose

### Project Structure
```
root/
├── apps/
│   ├── web/               # Next.js Admin Panel
│   └── app/               # Flutter Mobile App
├── cmd/                   # Go CLI Tools
├── pkg/                   # Shared Go Modules
├── scripts/               # Scripts and Tools
├── docs/                  # Documentation
├── Makefile              # Project Commands
└── docker-compose.yml    # Docker Configuration
```

## 🚀 Quick Start

### Prerequisites
- Node.js 18+
- Docker & Docker Compose
- PostgreSQL (or use Docker)

### Installation Steps

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd smg
   ```

2. **Set up environment variables**
   ```bash
   cp apps/web/.env.local.example apps/web/.env.local
   # Edit .env.local file and set your environment variables
   ```

3. **Start the database**
   ```bash
   make start
   ```

4. **Generate Prisma client**
   ```bash
   make generate
   ```

5. **Run database migrations**
   ```bash
   make migrate
   ```

6. **Start development server**
   ```bash
   make web
   ```

## 📱 Features

### User Dashboard
- **Topic Management**: Create and manage topics of interest
- **Media Links**: Connect various social media accounts
- **Article Management**: Browse and manage collected articles
- **System Settings**: Configure automatic reposting logic

### Admin Dashboard
- **User Management**: Manage all user accounts
- **Social Media Settings**: Configure platform API settings
- **System Settings**: Configure SMTP and other system settings

## 🛠️ Available Commands

| Command | Description |
|---------|-------------|
| `make start` | Start all services (Docker Compose) |
| `make stop` | Stop all services |
| `make web` | Start web development server |
| `make migrate` | Run database migrations |
| `make generate` | Generate Prisma client |
| `make studio` | Start Prisma Studio |
| `make install` | Install dependencies |
| `make clean` | Clean build artifacts |

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