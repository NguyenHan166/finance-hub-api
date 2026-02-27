package handlers

import (
	"finance-hub-api/internal/config"
	"finance-hub-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Router sets up all routes
type Router struct {
	cfg                *config.Config
	healthHandler      *HealthHandler
	accountHandler     *AccountHandler
	transactionHandler *TransactionHandler
	categoryHandler    *CategoryHandler
}

// NewRouter creates a new router
func NewRouter(
	cfg *config.Config,
	healthHandler *HealthHandler,
	accountHandler *AccountHandler,
	transactionHandler *TransactionHandler,
	categoryHandler *CategoryHandler,
) *Router {
	return &Router{
		cfg:                cfg,
		healthHandler:      healthHandler,
		accountHandler:     accountHandler,
		transactionHandler: transactionHandler,
		categoryHandler:    categoryHandler,
	}
}

// Setup sets up all routes
func (r *Router) Setup() *gin.Engine {
	// Create Gin router
	if r.cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware(r.cfg.CORS.AllowedOrigins))

	// Health check routes (no auth required)
	router.GET("/health", r.healthHandler.Health)
	router.GET("/ready", r.healthHandler.Ready)

	// API routes
	api := router.Group("/api/" + r.cfg.Server.APIVersion)
	{
		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(r.cfg.JWT.Secret))
		{
			// Account routes
			accounts := protected.Group("/accounts")
			{
				accounts.POST("", r.accountHandler.CreateAccount)
				accounts.GET("", r.accountHandler.GetAllAccounts)
				accounts.GET("/:id", r.accountHandler.GetAccount)
				accounts.PUT("/:id", r.accountHandler.UpdateAccount)
				accounts.DELETE("/:id", r.accountHandler.DeleteAccount)
			}

			// Transaction routes
			transactions := protected.Group("/transactions")
			{
				transactions.POST("", r.transactionHandler.CreateTransaction)
				transactions.GET("", r.transactionHandler.GetAllTransactions)
				transactions.GET("/:id", r.transactionHandler.GetTransaction)
				transactions.DELETE("/:id", r.transactionHandler.DeleteTransaction)
			}

			// Category routes
			categories := protected.Group("/categories")
			{
				categories.POST("", r.categoryHandler.CreateCategory)
				categories.GET("", r.categoryHandler.GetAllCategories)
				categories.GET("/:id", r.categoryHandler.GetCategory)
				categories.DELETE("/:id", r.categoryHandler.DeleteCategory)
			}
		}
	}

	return router
}
