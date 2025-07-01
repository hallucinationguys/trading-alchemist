package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"trading-alchemist/internal/application/usecases"
	"trading-alchemist/internal/config"
	"trading-alchemist/internal/infrastructure/database"
	"trading-alchemist/internal/infrastructure/email"
	postgres "trading-alchemist/internal/infrastructure/repositories/postgres"
	server "trading-alchemist/internal/presentation/http"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create context for database connection
	ctx := context.Background()

	// Initialize database connection
	db, err := database.NewConnection(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	magicLinkRepo := postgres.NewMagicLinkRepository(db)

	// Initialize services
	dbService := database.NewService(db)
	emailService, err := email.NewEmailService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize email service: %v", err)
	}

	// Initialize use cases
	authUseCase := usecases.NewAuthUseCase(userRepo, magicLinkRepo, emailService, cfg, dbService)

	// Initialize HTTP server
	srv := server.NewServer(cfg, authUseCase, userRepo)

	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
		log.Printf("Starting server on %s", addr)
		if err := srv.Start(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
} 