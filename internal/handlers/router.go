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
	authHandler        *AuthHandler
	accountHandler     *AccountHandler
	transactionHandler *TransactionHandler
	categoryHandler    *CategoryHandler
}

// NewRouter creates a new router
func NewRouter(
	cfg *config.Config,
	healthHandler *HealthHandler,
	authHandler *AuthHandler,
	accountHandler *AccountHandler,
	transactionHandler *TransactionHandler,
	categoryHandler *CategoryHandler,
) *Router {
	return &Router{
		cfg:                cfg,
		healthHandler:      healthHandler,
		authHandler:        authHandler,
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
		// Authentication routes (no auth required)
		auth := api.Group("/auth")
		{
			// Apply strict rate limiting to sensitive endpoints
			auth.POST("/register", middleware.StrictRateLimitMiddleware(), r.authHandler.Register)
			auth.POST("/login", middleware.StrictRateLimitMiddleware(), r.authHandler.Login)
			auth.POST("/forgot-password", middleware.StrictRateLimitMiddleware(), r.authHandler.RequestPasswordReset)
			auth.POST("/reset-password", middleware.StrictRateLimitMiddleware(), r.authHandler.ResetPassword)
			auth.POST("/resend-verification-email", middleware.StrictRateLimitMiddleware(), r.authHandler.ResendVerificationEmail)
			
			// Moderate rate limiting for other auth endpoints
			auth.POST("/refresh", r.authHandler.RefreshToken)
			auth.GET("/google", r.authHandler.InitiateGoogleOAuth)
			auth.GET("/google/callback", r.authHandler.HandleGoogleCallback)
			auth.POST("/google/token", r.authHandler.VerifyGoogleToken)
			auth.POST("/verify-email", r.authHandler.VerifyEmail)

			// Protected auth routes
			authProtected := auth.Group("")
			authProtected.Use(middleware.AuthMiddleware(r.cfg.JWT.Secret))
			{
				authProtected.GET("/profile", r.authHandler.GetProfile)
				authProtected.POST("/change-password", r.authHandler.ChangePassword)
				authProtected.POST("/logout", r.authHandler.Logout)
				authProtected.POST("/send-verification-email", r.authHandler.SendVerificationEmail)
			}
		}

		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(r.cfg.JWT.Secret))
		protected.Use(middleware.ModerateRateLimitMiddleware()) // Rate limit for all API endpoints
		{
			// Account routes
			accounts := protected.Group("/accounts")
			{
				accounts.GET("/summary", r.accountHandler.GetAccountSummary) // Must be before /:id
				accounts.GET("/banks", r.accountHandler.GetBanks)            // Must be before /:id
				accounts.POST("", r.accountHandler.CreateAccount)
				accounts.GET("", r.accountHandler.GetAllAccounts)
				accounts.GET("/:id", r.accountHandler.GetAccount)
				accounts.PUT("/:id", r.accountHandler.UpdateAccount)
				accounts.DELETE("/:id", r.accountHandler.DeleteAccount)
			}

			// Transaction routes
			transactions := protected.Group("/transactions")
			{
				// Special routes first (before /:id to avoid conflicts)
				transactions.GET("/recent", r.transactionHandler.GetRecentTransactions)
				transactions.GET("/summary", r.transactionHandler.GetTransactionSummary)
				transactions.PUT("/bulk/category", r.transactionHandler.BulkUpdateCategory)
				transactions.DELETE("/bulk", r.transactionHandler.BulkDelete)
				
				// Standard CRUD routes
				transactions.POST("", r.transactionHandler.CreateTransaction)
				transactions.GET("", r.transactionHandler.GetAllTransactions)
				transactions.GET("/:id", r.transactionHandler.GetTransaction)
				transactions.PUT("/:id", r.transactionHandler.UpdateTransaction)
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
