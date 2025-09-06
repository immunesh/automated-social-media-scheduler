package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/immunesh/automated-social-media-scheduler/internal/auth"
	"github.com/immunesh/automated-social-media-scheduler/internal/database"
	"github.com/immunesh/automated-social-media-scheduler/internal/models"
)

// API Auth Handlers

// APIRegister handles API registration
func APIRegister(c *gin.Context) {
	var req struct {
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=6"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"token": token,
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"plan":       user.Plan,
		},
	})
}

// APILogin handles API login
func APILogin(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if !auth.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"plan":       user.Plan,
		},
	})
}

// APILogout handles API logout
func APILogout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// API Post Handlers

// APIGetPosts returns all posts for the authenticated user
func APIGetPosts(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var posts []models.Post
	result := database.DB.Where("user_id = ?", userID).
		Preload("SocialAccounts").
		Order("created_at desc").
		Find(&posts)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// APICreatePost creates a new post
func APICreatePost(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Title       string    `json:"title" binding:"required"`
		Content     string    `json:"content" binding:"required"`
		Hashtags    string    `json:"hashtags"`
		MediaURLs   string    `json:"media_urls"`
		ScheduledAt time.Time `json:"scheduled_at"`
		Platforms   []uint    `json:"platforms"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create post
	post := models.Post{
		UserID:      userID.(uint),
		Title:       req.Title,
		Content:     req.Content,
		Hashtags:    req.Hashtags,
		MediaURLs:   req.MediaURLs,
		Status:      "draft",
		ScheduledAt: req.ScheduledAt,
	}

	if err := database.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	// Associate with social accounts if provided
	if len(req.Platforms) > 0 {
		var socialAccounts []models.SocialAccount
		database.DB.Where("id IN ? AND user_id = ?", req.Platforms, userID).Find(&socialAccounts)
		
		if err := database.DB.Model(&post).Association("SocialAccounts").Append(&socialAccounts); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to associate social accounts"})
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{"post": post})
}

// APIUpdatePost updates an existing post
func APIUpdatePost(c *gin.Context) {
	userID, _ := c.Get("user_id")
	postID := c.Param("id")

	var post models.Post
	if err := database.DB.Where("id = ? AND user_id = ?", postID, userID).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var req struct {
		Title       string    `json:"title"`
		Content     string    `json:"content"`
		Hashtags    string    `json:"hashtags"`
		MediaURLs   string    `json:"media_urls"`
		ScheduledAt time.Time `json:"scheduled_at"`
		Status      string    `json:"status"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update post
	if err := database.DB.Model(&post).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}

// APIDeletePost deletes a post
func APIDeletePost(c *gin.Context) {
	userID, _ := c.Get("user_id")
	postID := c.Param("id")

	result := database.DB.Where("id = ? AND user_id = ?", postID, userID).Delete(&models.Post{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}

// API Social Account Handlers

// APIGetSocialAccounts returns all connected social accounts
func APIGetSocialAccounts(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var accounts []models.SocialAccount
	if err := database.DB.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch social accounts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"social_accounts": accounts})
}

// APIConnectSocialAccount connects a new social media account
func APIConnectSocialAccount(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Platform    string `json:"platform" binding:"required"`
		AccountName string `json:"account_name" binding:"required"`
		AccountID   string `json:"account_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create social account
	account := models.SocialAccount{
		UserID:      userID.(uint),
		Platform:    req.Platform,
		AccountName: req.AccountName,
		AccountID:   req.AccountID,
		IsActive:    true,
	}

	if err := database.DB.Create(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect social account"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"social_account": account})
}

// APIDisconnectSocialAccount disconnects a social media account
func APIDisconnectSocialAccount(c *gin.Context) {
	userID, _ := c.Get("user_id")
	accountID := c.Param("id")

	result := database.DB.Where("id = ? AND user_id = ?", accountID, userID).Delete(&models.SocialAccount{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disconnect social account"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Social account not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Social account disconnected successfully"})
}

// API Analytics Handlers

// APIGetAnalytics returns analytics data
func APIGetAnalytics(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var analytics []models.Analytics
	result := database.DB.Joins("JOIN posts ON analytics.post_id = posts.id").
		Where("posts.user_id = ?", userID).
		Preload("Post").
		Preload("SocialAccount").
		Find(&analytics)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch analytics"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"analytics": analytics})
}

// API Profile Handlers

// APIGetProfile returns user profile
func APIGetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := database.DB.Preload("SocialAccounts").First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// APIUpdateProfile updates user profile
func APIUpdateProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Model(&user).Updates(req).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}