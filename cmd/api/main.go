package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"trading-alchemist/internal/application/auth"
	"trading-alchemist/internal/config"
	"trading-alchemist/internal/infrastructure/database"
	"trading-alchemist/internal/infrastructure/email"
	"trading-alchemist/internal/infrastructure/llm/agent"
	server "trading-alchemist/internal/presentation/http"
)

func main() {
	// Load configuration
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Configuration validation failed: %v", err)
	}

	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup database connection
	dbPool, err := database.NewConnection(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(dbPool)

	// Setup database service
	dbService := database.NewService(dbPool)

	// Setup email service
	emailService, err := email.NewEmailService(cfg)
	if err != nil {
		log.Fatalf("Failed to create email service: %v", err)
	}

	// Setup LLM service
	llmService, err := agent.NewLLMService()
	if err != nil {
		log.Fatalf("Failed to create LLM service: %v", err)
	}

	// Initialize use cases - repositories are now managed through dbService
	authUseCase := auth.NewAuthUseCase(emailService, cfg, dbService)

	// Initialize HTTP server
	httpServer := server.NewServer(cfg, authUseCase, dbService, llmService)

	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
		log.Printf("Starting server on %s", addr)
		if err := httpServer.Start(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Shutdown context with a timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited properly")
} 