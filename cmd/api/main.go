package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vendix/internal/config"
	"vendix/internal/database"
	"vendix/internal/logger"
	"vendix/internal/server"

	_ "vendix/docs" // Swagger documentation

	"github.com/joho/godotenv"
)

// @title Vendix API
// @version 1.0
// @description Electronic Billing SaaS for Dominican Republic
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@billing.example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize logger
	logger.InitLogger()
	defer logger.Sync()

	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	logger.Info("Database connection established")

	// Run migrations for public schema
	if err := database.RunPublicMigrations(db); err != nil {
		logger.Fatal("Failed to run public migrations", "error", err)
	}

	logger.Info("Public migrations completed")

	// Initialize server
	srv := server.NewServer(cfg, db)

	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf(":%d", cfg.AppPort)
		logger.Info("Starting server", "address", addr)
		if err := srv.Listen(addr); err != nil {
			logger.Fatal("Server failed to start", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.ShutdownWithContext(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited gracefully")
}
