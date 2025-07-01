# Trading Alchemist

A modern Go-based authentication system with magic link functionality using Clean Architecture principles.

## ✨ Features

- 🚀 **Fiber Web Framework** - Fast HTTP framework for Go
- 🗄️ **PostgreSQL** - Robust database with proper indexing
- 📧 **Email Integration** - Resend integration for sending magic links
- 🔒 **JWT Authentication** - Secure token-based authentication
- 🐳 **Docker Support** - Easy development environment setup
- 🏗️ **Clean Architecture** - Domain-driven design with proper separation
- 📦 **SQLC Integration** - Type-safe SQL queries
- 🔄 **Database Migrations** - Version-controlled schema changes
- 📖 **Swagger Documentation** - Interactive API documentation with OpenAPI

## 🏗️ Architecture

This project follows Clean Architecture principles with clear separation of concerns:

```
├── cmd/api/              # Application entry point
├── internal/
│   ├── application/      # Use cases and DTOs
│   ├── domain/          # Business logic and entities
│   ├── infrastructure/  # External concerns (DB, email)
│   └── presentation/    # HTTP handlers and routes
└── pkg/                 # Shared utilities
```

## 🚀 Development Setup

You have two options for setting up the development environment:

### Option 1: Devbox (Recommended)

[Devbox](https://www.jetify.com/devbox) provides a reproducible development environment with all tools pre-installed.

```bash
# Install Devbox (if not already installed)
curl -fsSL https://get.jetify.com/devbox | bash

# Enter development environment
devbox shell

# Setup project (first time)
devbox run install-tools
devbox run setup
devbox run db-setup

# Start development server
devbox run dev
```

For detailed Devbox instructions, see [.devbox/README.md](.devbox/README.md).

### Option 2: Manual Setup

#### Prerequisites

- Go 1.24.3+
- Docker and Docker Compose
- PostgreSQL client tools
- migrate CLI tool

#### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd trading-alchemist
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Setup environment**
   ```bash
   cp configs/app.env.example configs/app.env
   # Edit configs/app.env with your settings
   ```

4. **Start services**
   ```bash
   # This will start PostgreSQL container
   make docker-up
   
   # Run database migrations
   make migrate-up
   
   # Generate SQLC code
   make sqlc-generate
   ```

5. **Run the application**
   ```bash
   make run
   ```

   The API will be available at `http://localhost:8080`

## 🐳 Docker Development

For containerized development:

```bash
# Start development environment with hot-reload
make docker-dev

# Or build and start
make docker-dev-build

# View logs
make docker-logs
```

This command will:
- Start PostgreSQL container
- Run database migrations
- Generate SQLC code
- Start the application with hot-reload

## 📝 Configuration

Copy the example environment file and modify as needed:

```bash
cp configs/app.env.example configs/app.env
```

### Key Configuration Options

```bash
# Server
SERVER_HOST=localhost
SERVER_PORT=8080

# Database (Note: Using port 5433 to avoid conflicts)
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=trading_alchemist_db

# JWT
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
JWT_TTL=24h

# Email (Resend)
RESEND_API_KEY=re_xxxxxxxxx
FROM_EMAIL=noreply@yourdomain.com
FROM_NAME=Trading Alchemist

# Application
APP_NAME=Trading Alchemist
APP_ENV=development
APP_BASE_URL=http://localhost:8080
MAGIC_LINK_TTL=15m
```

## 🗄️ Database

The project uses PostgreSQL with SQLC for type-safe queries.

### Migrations

```bash
# Run migrations
make migrate-up

# Create new migration
make migrate-create name=add_new_table

# Rollback migrations
make migrate-down
```

### SQLC Code Generation

```bash
# Generate Go code from SQL queries
make sqlc-generate
```

## 🧪 Testing

```bash
# Run all tests
make test

# Run tests in Docker
make docker-test
```

## 📊 API Endpoints

### Authentication

- `POST /api/auth/send-magic-link` - Send magic link to email
- `POST /api/auth/verify` - Verify magic link and get JWT token

### Users

- `GET /api/users/profile` - Get user profile (requires authentication)

### Health

- `GET /health` - Application health check

## 🏃‍♂️ Available Commands

Use `make help` to see all available commands:

```bash
# Development
make run              # Run application locally
make test             # Run tests
make build            # Build binary
make clean            # Clean build artifacts

# Docker
make docker-up        # Start services
make docker-down      # Stop services
make docker-dev       # Development with hot-reload
make docker-logs      # View logs

# Database
make migrate-up       # Run migrations
make migrate-down     # Rollback migrations
make sqlc-generate    # Generate code from SQL

# Setup helpers
make dev-setup        # Complete local setup
make docker-dev-setup # Complete Docker setup
```

## 🚀 Production Deployment

See [docs/DOCKER.md](docs/DOCKER.md) for detailed Docker deployment instructions.

Quick production deployment:

```bash
# Build and deploy with Docker
make deploy-prod
```

## 📚 Documentation

- [Docker Setup Guide](docs/DOCKER.md) - Comprehensive Docker documentation
- [Devbox Setup Guide](.devbox/README.md) - Devbox development environment
- [API Documentation](docs/api/) - Interactive Swagger API documentation
- **Swagger UI:** Available at `http://localhost:8080/docs` when server is running
- **API Schema:** Available at `http://localhost:8080/swagger.json`

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
