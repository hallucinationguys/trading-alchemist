package main

import (
	"context"
	"log"
	"trading-alchemist/internal/config"
	"trading-alchemist/internal/infrastructure/database"
	"trading-alchemist/internal/infrastructure/database/seeder"
)

func main() {
	log.Println("Starting database seeder...")

	// Load configuration
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Create context
	ctx := context.Background()

	// Setup database connection
	dbPool, err := database.NewConnection(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(dbPool)

	// Setup database service
	dbService := database.NewService(dbPool)

	// Run the seeder
	seeder.Seed(dbService)

	log.Println("Seeder finished successfully.")
} 