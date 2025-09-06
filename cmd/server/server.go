package server

import (
	"github.com/gin-gonic/gin"
	"github.com/immunesh/automated-social-media-scheduler/internal/database"
	"github.com/immunesh/automated-social-media-scheduler/internal/handlers"
	"github.com/immunesh/automated-social-media-scheduler/internal/middleware"
)

func Run(addr string) error {
	// Initialize database
	if err := database.Initialize(); err != nil {
		return err
	}

	// Set up Gin router
	r := gin.Default()

	// Load HTML templates
	r.LoadHTMLGlob("web/templates/*")
	r.Static("/static", "./web/static")

	// Setup routes
	setupRoutes(r)

	return r.Run(addr)
}

func setupRoutes(r *gin.Engine) {
	// Public routes
	r.GET("/", handlers.HomePage)
	r.GET("/register", handlers.RegisterPage)
	r.POST("/register", handlers.Register)
	r.GET("/login", handlers.LoginPage)
	r.POST("/login", handlers.Login)

	// API routes
	api := r.Group("/api")
	{
		// Auth endpoints
		api.POST("/auth/register", handlers.APIRegister)
		api.POST("/auth/login", handlers.APILogin)
		api.POST("/auth/logout", handlers.APILogout)
	}

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/dashboard", handlers.Dashboard)
		protected.GET("/posts", handlers.PostsPage)
		protected.POST("/posts", handlers.CreatePost)
		protected.GET("/analytics", handlers.AnalyticsPage)
		protected.GET("/settings", handlers.SettingsPage)
	}

	// Protected API routes
	apiAuth := api.Group("/")
	apiAuth.Use(middleware.AuthRequired())
	{
		// Posts endpoints
		apiAuth.GET("/posts", handlers.APIGetPosts)
		apiAuth.POST("/posts", handlers.APICreatePost)
		apiAuth.PUT("/posts/:id", handlers.APIUpdatePost)
		apiAuth.DELETE("/posts/:id", handlers.APIDeletePost)

		// Social media accounts
		apiAuth.GET("/social-accounts", handlers.APIGetSocialAccounts)
		apiAuth.POST("/social-accounts", handlers.APIConnectSocialAccount)
		apiAuth.DELETE("/social-accounts/:id", handlers.APIDisconnectSocialAccount)

		// Analytics
		apiAuth.GET("/analytics", handlers.APIGetAnalytics)

		// User profile
		apiAuth.GET("/profile", handlers.APIGetProfile)
		apiAuth.PUT("/profile", handlers.APIUpdateProfile)
	}
}