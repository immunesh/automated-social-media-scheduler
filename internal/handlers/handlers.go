package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/immunesh/automated-social-media-scheduler/internal/auth"
	"github.com/immunesh/automated-social-media-scheduler/internal/database"
	"github.com/immunesh/automated-social-media-scheduler/internal/models"
)

// HomePage renders the home page
func HomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Social Media Scheduler",
	})
}

// RegisterPage renders the registration page
func RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title": "Register",
	})
}

// LoginPage renders the login page
func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

// Register handles user registration
func Register(c *gin.Context) {
	var req struct {
		Email     string `form:"email" json:"email" binding:"required,email"`
		Password  string `form:"password" json:"password" binding:"required,min=6"`
		FirstName string `form:"first_name" json:"first_name"`
		LastName  string `form:"last_name" json:"last_name"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"title": "Register",
			"error": err.Error(),
		})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.HTML(http.StatusBadRequest, "register.html", gin.H{
			"title": "Register",
			"error": "User already exists",
		})
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"title": "Register",
			"error": "Failed to create user",
		})
		return
	}

	// Create user
	user := models.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Plan:      "free",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "register.html", gin.H{
			"title": "Register",
			"error": "Failed to create user",
		})
		return
	}

	c.Redirect(http.StatusFound, "/login")
}

// Login handles user login
func Login(c *gin.Context) {
	var req struct {
		Email    string `form:"email" json:"email" binding:"required,email"`
		Password string `form:"password" json:"password" binding:"required"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{
			"title": "Login",
			"error": err.Error(),
		})
		return
	}

	// Find user
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"title": "Login",
			"error": "Invalid credentials",
		})
		return
	}

	// Check password
	if !auth.CheckPasswordHash(req.Password, user.Password) {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"title": "Login",
			"error": "Invalid credentials",
		})
		return
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"title": "Login",
			"error": "Failed to generate token",
		})
		return
	}

	// Set cookie
	c.SetCookie("auth_token", token, 3600*24, "/", "", false, true)
	c.Redirect(http.StatusFound, "/dashboard")
}

// Dashboard renders the main dashboard
func Dashboard(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	// Get user data
	var user models.User
	database.DB.Preload("SocialAccounts").Preload("Posts").First(&user, userID)
	
	// Get recent posts
	var recentPosts []models.Post
	database.DB.Where("user_id = ?", userID).Order("created_at desc").Limit(5).Find(&recentPosts)
	
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"title": "Dashboard",
		"user":  user,
		"posts": recentPosts,
	})
}

// PostsPage renders the posts management page
func PostsPage(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var posts []models.Post
	database.DB.Where("user_id = ?", userID).Preload("SocialAccounts").Order("created_at desc").Find(&posts)
	
	var socialAccounts []models.SocialAccount
	database.DB.Where("user_id = ? AND is_active = ?", userID, true).Find(&socialAccounts)
	
	c.HTML(http.StatusOK, "posts.html", gin.H{
		"title":          "Posts",
		"posts":          posts,
		"socialAccounts": socialAccounts,
	})
}

// CreatePost handles post creation
func CreatePost(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var req struct {
		Title       string   `form:"title" json:"title" binding:"required"`
		Content     string   `form:"content" json:"content" binding:"required"`
		Hashtags    string   `form:"hashtags" json:"hashtags"`
		ScheduledAt string   `form:"scheduled_at" json:"scheduled_at"`
		Platforms   []string `form:"platforms" json:"platforms"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse scheduled time
	var scheduledAt time.Time
	if req.ScheduledAt != "" {
		var err error
		scheduledAt, err = time.Parse("2006-01-02T15:04", req.ScheduledAt)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid scheduled time"})
			return
		}
	}

	// Create post
	post := models.Post{
		UserID:      userID.(uint),
		Title:       req.Title,
		Content:     req.Content,
		Hashtags:    req.Hashtags,
		Status:      "draft",
		ScheduledAt: scheduledAt,
	}

	if err := database.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.Redirect(http.StatusFound, "/posts")
}

// AnalyticsPage renders the analytics page
func AnalyticsPage(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	// Get analytics data
	var analytics []models.Analytics
	database.DB.Joins("JOIN posts ON analytics.post_id = posts.id").
		Where("posts.user_id = ?", userID).
		Preload("Post").
		Find(&analytics)
	
	c.HTML(http.StatusOK, "analytics.html", gin.H{
		"title":     "Analytics",
		"analytics": analytics,
	})
}

// SettingsPage renders the settings page
func SettingsPage(c *gin.Context) {
	userID, _ := c.Get("user_id")
	
	var user models.User
	database.DB.Preload("SocialAccounts").First(&user, userID)
	
	c.HTML(http.StatusOK, "settings.html", gin.H{
		"title": "Settings",
		"user":  user,
	})
}