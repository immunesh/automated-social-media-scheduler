package services

import (
	"log"
	"time"

	"gorm.io/gorm"
	"github.com/immunesh/automated-social-media-scheduler/internal/database"
	"github.com/immunesh/automated-social-media-scheduler/internal/models"
)

// SchedulerService handles post scheduling
type SchedulerService struct {
	db *gorm.DB
}

func NewSchedulerService() *SchedulerService {
	return &SchedulerService{
		db: database.GetDB(),
	}
}

// StartScheduler starts the post scheduling service
func (s *SchedulerService) StartScheduler() {
	ticker := time.NewTicker(1 * time.Minute) // Check every minute
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.processScheduledPosts()
		}
	}
}

// processScheduledPosts finds and processes posts that are due to be published
func (s *SchedulerService) processScheduledPosts() {
	var posts []models.Post
	now := time.Now()

	// Find posts that are scheduled to be published now or in the past
	result := s.db.Preload("SocialAccounts").
		Where("status = ? AND scheduled_at <= ?", "scheduled", now).
		Find(&posts)

	if result.Error != nil {
		log.Printf("Error fetching scheduled posts: %v", result.Error)
		return
	}

	for _, post := range posts {
		s.publishPost(&post)
	}
}

// publishPost publishes a post to the connected social media platforms
func (s *SchedulerService) publishPost(post *models.Post) {
	log.Printf("Publishing post: %s", post.Title)

	// Update post status to publishing
	s.db.Model(post).Update("status", "publishing")

	success := true
	
	for _, account := range post.SocialAccounts {
		err := s.publishToAccount(post, &account)
		if err != nil {
			log.Printf("Error publishing to %s (%s): %v", account.Platform, account.AccountName, err)
			success = false
			
			// Update the relationship status
			s.db.Model(&models.PostSocialAccount{}).
				Where("post_id = ? AND social_account_id = ?", post.ID, account.ID).
				Updates(models.PostSocialAccount{
					Status: "failed",
					ErrorMessage: err.Error(),
				})
		} else {
			// Update the relationship status
			s.db.Model(&models.PostSocialAccount{}).
				Where("post_id = ? AND social_account_id = ?", post.ID, account.ID).
				Updates(models.PostSocialAccount{
					Status: "published",
					PublishedAt: time.Now(),
				})
		}
	}

	// Update overall post status
	var finalStatus string
	if success {
		finalStatus = "published"
	} else {
		finalStatus = "failed"
	}

	s.db.Model(post).Updates(models.Post{
		Status:      finalStatus,
		PublishedAt: time.Now(),
	})
}

// publishToAccount publishes a post to a specific social media account
func (s *SchedulerService) publishToAccount(post *models.Post, account *models.SocialAccount) error {
	// This is a placeholder for actual social media API integration
	// In a real implementation, you would:
	// 1. Use the account's access token to authenticate with the platform API
	// 2. Format the post content according to platform requirements
	// 3. Upload any media files
	// 4. Create the post via API
	// 5. Handle platform-specific errors and rate limits
	
	log.Printf("Simulating publication to %s account '%s'", account.Platform, account.AccountName)
	
	// Simulate API call delay
	time.Sleep(100 * time.Millisecond)
	
	// For demo purposes, we'll just log success
	// In reality, this would make HTTP requests to platform APIs
	return nil
}

// AnalyticsService handles fetching and updating post analytics
type AnalyticsService struct {
	db *gorm.DB
}

func NewAnalyticsService() *AnalyticsService {
	return &AnalyticsService{
		db: database.GetDB(),
	}
}

// UpdateAnalytics fetches latest analytics data for published posts
func (a *AnalyticsService) UpdateAnalytics() {
	log.Println("Updating analytics data...")
	
	var posts []models.Post
	result := a.db.Preload("SocialAccounts").
		Where("status = ?", "published").
		Find(&posts)

	if result.Error != nil {
		log.Printf("Error fetching published posts: %v", result.Error)
		return
	}

	for _, post := range posts {
		for _, account := range post.SocialAccounts {
			a.fetchAnalyticsForPost(&post, &account)
		}
	}
}

// fetchAnalyticsForPost fetches analytics data for a specific post on a platform
func (a *AnalyticsService) fetchAnalyticsForPost(post *models.Post, account *models.SocialAccount) {
	// This is a placeholder for actual analytics API integration
	// In a real implementation, you would:
	// 1. Use the account's access token to authenticate with platform analytics API
	// 2. Fetch engagement metrics (likes, comments, shares, impressions, reach)
	// 3. Store or update the analytics data in the database
	
	log.Printf("Fetching analytics for post '%s' on %s", post.Title, account.Platform)
	
	// For demo purposes, create sample analytics data
	analytics := models.Analytics{
		PostID:          post.ID,
		SocialAccountID: account.ID,
		Platform:        account.Platform,
		Likes:          10,  // Would be fetched from API
		Comments:       2,   // Would be fetched from API
		Shares:         1,   // Would be fetched from API
		Impressions:    100, // Would be fetched from API
		Reach:          80,  // Would be fetched from API
	}

	// Check if analytics record already exists
	var existing models.Analytics
	if err := a.db.Where("post_id = ? AND social_account_id = ?", post.ID, account.ID).
		First(&existing).Error; err != nil {
		// Create new record
		a.db.Create(&analytics)
	} else {
		// Update existing record
		a.db.Model(&existing).Updates(analytics)
	}
}