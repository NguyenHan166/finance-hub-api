package main

import (
	"finance-hub-api/internal/config"
	"finance-hub-api/internal/handlers"
	"finance-hub-api/internal/repositories"
	"finance-hub-api/internal/services"
	"finance-hub-api/pkg/database"
	"finance-hub-api/pkg/logger"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	logger.Log.Info.Printf("🚀 Starting Finance Hub API")
	logger.Log.Info.Printf("Environment: %s", cfg.Server.Env)
	logger.Log.Info.Printf("Port: %s", cfg.Server.Port)

	// Connect to database
	dbCfg := database.Config{
		URI:      cfg.Database.URI,
		Database: cfg.Database.Database,
	}
	
	db, err := database.NewConnection(dbCfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	accountRepo := repositories.NewAccountRepository(db.Database)
	transactionRepo := repositories.NewTransactionRepository(db.Database)
	categoryRepo := repositories.NewCategoryRepository(db.Database)

	// Initialize services
	accountService := services.NewAccountService(accountRepo)
	transactionService := services.NewTransactionService(transactionRepo, accountRepo)
	categoryService := services.NewCategoryService(categoryRepo)

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	accountHandler := handlers.NewAccountHandler(accountService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Setup router
	router := handlers.NewRouter(
		cfg,
		healthHandler,
		accountHandler,
		transactionHandler,
		categoryHandler,
	)

	engine := router.Setup()

	// Start server
	serverAddr := fmt.Sprintf(":%s", cfg.Server.Port)
	logger.Log.Info.Printf("✨ Server starting on http://localhost%s", serverAddr)
	logger.Log.Info.Printf("📚 API Documentation: http://localhost%s/api/%s", serverAddr, cfg.Server.APIVersion)

	// Graceful shutdown
	go func() {
		if err := engine.Run(serverAddr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info.Println("🛑 Shutting down server...")
	logger.Log.Info.Println("👋 Server stopped")
}
