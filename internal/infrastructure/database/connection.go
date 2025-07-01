package database

import (
	"context"
	"fmt"
	"time"

	"trading-alchemist/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewConnection creates a new PostgreSQL connection pool
func NewConnection(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	// Build database URL
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	dbConfig, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// Set connection pool configuration
	dbConfig.MaxConns = int32(cfg.Database.MaxConns)
	dbConfig.MinConns = int32(cfg.Database.MinConns)
	
	// Parse connection lifetime settings
	maxConnLife, err := time.ParseDuration(cfg.Database.MaxConnLife)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_CONN_LIFE configuration: %w", err)
	}
	dbConfig.MaxConnLifetime = maxConnLife
	
	maxConnIdle, err := time.ParseDuration(cfg.Database.MaxConnIdle)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_MAX_CONN_IDLE configuration: %w", err)
	}
	dbConfig.MaxConnIdleTime = maxConnIdle

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

// Close closes the database connection pool
func Close(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
} 