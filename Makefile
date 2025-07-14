.PHONY: start stop web migrate clean install

# Start all services using Docker Compose
start:
	docker-compose up --build -d

# Stop all services
stop:
	docker-compose down

# Start the web development server
web:
	cd apps/web && npm run dev

# Run database migrations
migrate:
	cd apps/web && npx prisma migrate dev

# Generate Prisma client
generate:
	cd apps/web && npx prisma generate

# Reset database
reset:
	cd apps/web && npx prisma migrate reset

# Install dependencies
install:
	cd apps/web && npm install

# Clean build artifacts
clean:
	cd apps/web && rm -rf .next node_modules
	docker-compose down --volumes

# Setup development environment
setup: install generate
	@echo "Development environment setup complete"

# View container logs
logs:
	docker-compose logs -f

# Start Prisma Studio for database management
studio:
	cd apps/web && npx prisma studio

# Format code using prettier
format:
	cd apps/web && npm run format

# Build application for production
build:
	cd apps/web && npm run build