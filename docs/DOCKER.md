# Docker Setup for Trading Alchemist

This document explains how to use Docker for development and production deployment of the Trading Alchemist application.

## Quick Start

### Development Environment

```bash
# Start development environment with hot-reload
make docker-dev

# Or build and start
make docker-dev-build

# View logs
make docker-logs
```

### Production Environment

```bash
# Build and start production services
make docker-prod-build

# Or use the production compose file directly
docker-compose -f docker/compose/docker-compose.prod.yml up -d
```

## Project Structure

The Docker files are organized in a dedicated `docker/` directory following best practices:

```
docker/
├── app/
│   ├── Dockerfile          # Production-optimized multi-stage build
│   └── Dockerfile.dev      # Development image with hot-reload using Air
├── compose/
│   ├── docker-compose.yml           # Base configuration with app and PostgreSQL
│   ├── docker-compose.override.yml  # Development overrides (automatically used)
│   └── docker-compose.prod.yml      # Production configuration with secrets and resource limits
└── postgres/
    └── init.sql             # Database initialization script
```

### Additional Files

- **`.dockerignore`** - Excludes unnecessary files from build context
- **`.air.toml`** - Configuration for Air hot-reload tool (in project root)

### Benefits of This Structure

- **🗂️ Organization**: All Docker-related files are in one place
- **🔧 Separation of Concerns**: Different file types are in dedicated subdirectories
- **📖 Clarity**: Easy to understand what each directory contains
- **🚀 Scalability**: Can easily add more Docker configurations (e.g., staging)
- **🛡️ Security**: Better control over what gets copied into containers

## Available Commands

### Building Images

```bash
# Build production image
make docker-build

# Build development image
make docker-build-dev

# Build with custom script
./scripts/build.sh --tag v1.0.0 --push
```

### Running Services

```bash
# Development (with hot-reload)
make docker-dev                # Start dev environment
make docker-dev-build          # Build and start

# Production
make docker-prod               # Start production
make docker-prod-build         # Build and start production

# Basic services
make docker-up                 # Start all services
make docker-down               # Stop all services
make docker-restart            # Restart services
```

### Database Operations

```bash
# Run migrations in Docker
make docker-migrate-up
make docker-migrate-down

# Access database
docker-compose exec postgres psql -U postgres -d trading_alchemist_db

# Or connect from host (note the port 5433)
psql -h localhost -p 5433 -U postgres -d trading_alchemist_db
```

### Monitoring & Debugging

```bash
# View logs
make docker-logs

# Check service health
make docker-health

# View container stats
make docker-stats

# Run tests in container
make docker-test
```

## Environment Variables

### Development

Set in `docker/compose/docker-compose.override.yml`:

- `APP_ENV=development`
- `JWT_SECRET=dev-secret-key-change-in-production`
- `MAGIC_LINK_TTL=15m`

### Production

Use Docker secrets for sensitive data:

- `JWT_SECRET_FILE=/run/secrets/jwt_secret`
- `RESEND_API_KEY_FILE=/run/secrets/resend_api_key`
- `DB_PASSWORD_FILE=/run/secrets/db_password`

## Production Deployment

### 1. Prepare Secrets

```bash
# Create external secrets
echo "your-jwt-secret" | docker secret create jwt_secret -
echo "your-resend-api-key" | docker secret create resend_api_key -
echo "your-db-password" | docker secret create db_password -

```

### 2. Create External Volume

```bash
# Create external volume for production data
docker volume create postgres_data_prod
```

### 3. Deploy

```bash
# Deploy with production configuration
make deploy-prod

# Or manually
docker-compose -f docker/compose/docker-compose.prod.yml up -d
```

## Image Optimization

The production Dockerfile uses several optimization techniques:

- **Multi-stage build** - Separate build and runtime environments
- **Minimal base image** - Uses `scratch` for smallest possible image
- **Static binary** - CGO disabled for portability
- **Security** - Non-root user execution
- **Layer caching** - go.mod copied first for better build caching

## Health Checks

Both development and production configurations include health checks:

- **Application health** - HTTP endpoint check on `/health`
- **Database health** - PostgreSQL readiness check

## Networking

Services communicate through a dedicated Docker network:

- **Network name**: `trading_alchemist_network`
- **Driver**: bridge
- **Service resolution**: Services can reach each other by service name

## Volumes

### Development
- Source code mounted for hot-reload: `.:/app`
- Excluded directories: `/app/tmp`, `/app/vendor`

### Production
- External volume for PostgreSQL data: `postgres_data_prod`
- No source code mounts for security

## Troubleshooting

### Common Issues

1. **Port conflicts**
   ```bash
   # Stop conflicting services
   sudo lsof -i :5433  # Check PostgreSQL port (changed from 5432 to avoid conflicts)

   sudo lsof -i :8080  # Check app port
   
   # Note: PostgreSQL runs on port 5433 externally to avoid conflicts with system PostgreSQL
   ```

2. **Permission issues**
   ```bash
   # Make scripts executable
   chmod +x scripts/*.sh
   ```

3. **Build failures**
   ```bash
   # Clean Docker resources
   make docker-clean
   
   # Rebuild from scratch
   docker-compose build --no-cache
   ```

4. **Database connection issues**
   ```bash
   # Check database logs
   docker-compose logs postgres
   
   # Verify database is ready
   docker-compose exec postgres pg_isready -U postgres
   ```

### Useful Commands

```bash
# Shell into running container
docker-compose exec app sh

# View container processes
docker-compose top

# Inspect container configuration
docker-compose config

# Remove all containers and volumes
docker-compose down -v --remove-orphans
```

## CI/CD Integration

Example GitHub Actions workflow:

```yaml
name: Build and Push Docker Image

on:
  push:
    tags: ['v*']

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Build and push
        run: |
          ./scripts/build.sh \
            --tag ${{ github.ref_name }} \
            --registry ghcr.io/yourusername \
            --push
```

## Performance Considerations

### Resource Limits (Production)

- **App**: 512MB memory, 0.5 CPU limit
- **PostgreSQL**: 1GB memory, 1 CPU limit

### Scaling

For horizontal scaling, update the production compose file:

```yaml
services:
  app:
    deploy:
      replicas: 3  # Run 3 instances
``` 