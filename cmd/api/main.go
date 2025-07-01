package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"trading-alchemist/internal/config"
	"trading-alchemist/internal/infrastructure/database"
	"trading-alchemist/internal/infrastructure/email"
	server "trading-alchemist/internal/presentation/http"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Create context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to database
	db, err := database.NewConnection(ctx, cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close(db)

	// Initialize email service (Resend only)
	emailService, err := email.NewEmailService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize email service: %v", err)
	}

	// Create HTTP server
	httpServer := server.NewServer(cfg, db, emailService)

	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
		log.Printf("Starting server on %s", addr)
		
		if err := httpServer.Start(addr); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
} 