.PHONY: start stop web migrate clean install api scheduler build-go

# Start all services using Docker Compose
start:
	docker-compose up --build -d

# Stop all services
stop:
	docker-compose down

# Start the web development server
web:
	cd apps/web && npm run dev

# Start the Go API server
api:
	cd cmd/api && go run main.go

# Start the scheduler service
scheduler:
	cd cmd/scheduler && go run main.go

# Run database migrations using Prisma
migrate:
	cd apps/web && npx prisma migrate dev

# Run database migrations using Go migrate tool
migrate-go:
	cd cmd/migrate && go run main.go up

# Create new migration
migrate-create:
	cd cmd/migrate && go run main.go create $(name)

# Generate Prisma client
generate:
	cd apps/web && npx prisma generate

# Reset database
reset:
	cd apps/web && npx prisma migrate reset

# Install dependencies
install:
	cd apps/web && npm install
	go mod tidy

# Install Flutter dependencies
install-flutter:
	cd apps/app && flutter pub get

# Clean build artifacts
clean:
	cd apps/web && rm -rf .next node_modules
	cd apps/app && flutter clean
	docker-compose down --volumes

# Setup development environment
setup: install install-flutter generate
	@echo "Development environment setup complete"

# View container logs
logs:
	docker-compose logs -f

# Start Prisma Studio for database management
studio:
	cd apps/web && npx prisma studio

# Format code
format:
	cd apps/web && npm run format
	cd apps/app && flutter format lib/
	go fmt ./...

# Build application for production
build:
	cd apps/web && npm run build

# Build Go services
build-go:
	go build -o bin/api cmd/api/main.go
	go build -o bin/scheduler cmd/scheduler/main.go
	go build -o bin/migrate cmd/migrate/main.go

# Run Go tests
test-go:
	go test ./...

# Run Flutter tests
test-flutter:
	cd apps/app && flutter test

# Build Flutter app
build-flutter:
	cd apps/app && flutter build apk

# Run Flutter app
run-flutter:
	cd apps/app && flutter run

# Check Go code
check-go:
	go vet ./...
	golint ./...

# Create .env file from example
env:
	cp .env.example .env
	@echo "Created .env file. Please update the values as needed."