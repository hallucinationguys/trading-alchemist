#!/bin/sh

# Docker entrypoint script for Trading Alchemist

set -e

# Function to wait for service to be ready
wait_for_service() {
    local host=$1
    local port=$2
    local service_name=$3
    
    echo "Waiting for $service_name to be ready on $host:$port..."
    
    while ! nc -z "$host" "$port"; do
        echo "Waiting for $service_name..."
        sleep 2
    done
    
    echo "$service_name is ready!"
}

# Function to run database migrations
run_migrations() {
    echo "Running database migrations..."
    
    # Wait for database to be ready
    wait_for_service "${DB_HOST:-postgres}" "${DB_PORT:-5432}" "PostgreSQL"
    
    # Run migrations if migrate tool is available
    if command -v migrate >/dev/null 2>&1; then
        migrate -path /app/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL_MODE}" up
        echo "Migrations completed successfully!"
    else
        echo "Migrate tool not found, skipping migrations"
    fi
}

# Function for health check
health_check() {
    echo "Performing health check..."
    
    # Simple health check - verify the application is responsive
    if command -v wget >/dev/null 2>&1; then
        wget --quiet --tries=1 --spider "http://localhost:${SERVER_PORT:-8080}/health" && echo "Health check passed"
    elif command -v curl >/dev/null 2>&1; then
        curl -f "http://localhost:${SERVER_PORT:-8080}/health" && echo "Health check passed"
    else
        echo "No HTTP client available for health check"
        exit 1
    fi
}

# Main execution logic
case "$1" in
    "migrate")
        run_migrations
        ;;
    "health-check"|"--health-check")
        health_check
        ;;
    "")
        # Default behavior - start the application
        echo "Starting Trading Alchemist..."
        
        # Wait for dependencies if in Docker environment
        if [ "${DOCKER_ENV:-false}" = "true" ]; then
            wait_for_service "${DB_HOST:-postgres}" "${DB_PORT:-5432}" "PostgreSQL"
        fi
        
        # Start the application
        exec /app/main
        ;;
    *)
        # Pass through any other commands
        exec "$@"
        ;;
esac 