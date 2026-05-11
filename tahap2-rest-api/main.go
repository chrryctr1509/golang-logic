package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/user/tahap2-rest-api/config"
	"github.com/user/tahap2-rest-api/internal/handler"
	"github.com/user/tahap2-rest-api/internal/middleware"
	"github.com/user/tahap2-rest-api/internal/repository"
	"github.com/user/tahap2-rest-api/internal/service"
	"github.com/user/tahap2-rest-api/internal/worker"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	_ = godotenv.Load()

	// Load config
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db := config.GetDB()
	if db == nil {
		log.Fatal("Failed to connect to database")
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	txRepo := repository.NewTransactionRepository(db)

	// Initialize worker
	tw := worker.NewTransferWorker(db, txRepo)
	transferQ := tw.Queue

	// Start worker in background
	ctx, cancel := context.WithCancel(context.Background())
	go tw.Start(ctx)
	defer cancel()

	// Initialize services
	authSvc := service.NewAuthService(userRepo)
	txSvc := service.NewTransactionService(db, txRepo, userRepo, transferQ)
	userSvc := service.NewUserService(userRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authSvc)
	txHandler := handler.NewTransactionHandler(txSvc)
	userHandler := handler.NewUserHandler(userSvc)

	// Setup Gin router
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.DebugMode)
	}
	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Public routes
	router.POST("/register", authHandler.Register)
	router.POST("/login", authHandler.Login)

	// Protected routes
	api := router.Group("/api/v1")
	api.Use(middleware.JWTAuth())

	api.POST("/topup", txHandler.TopUp)
	api.POST("/pay", txHandler.Payment)
	api.POST("/transfer", txHandler.Transfer)
	api.GET("/transactions", txHandler.GetTransactions)
	api.PUT("/profile", userHandler.UpdateProfile)

	port := config.GetConfig().AppPort
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
