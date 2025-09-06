package models

import (
	"time"
)

// User represents a user account
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Plan      string    `json:"plan" gorm:"default:free"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	SocialAccounts []SocialAccount `json:"social_accounts"`
	Posts          []Post          `json:"posts"`
}

// SocialAccount represents connected social media accounts
type SocialAccount struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id"`
	Platform     string    `json:"platform"` // facebook, twitter, instagram, linkedin
	AccountName  string    `json:"account_name"`
	AccountID    string    `json:"account_id"`
	AccessToken  string    `json:"-"`
	RefreshToken string    `json:"-"`
	ExpiresAt    time.Time `json:"expires_at"`
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relationships
	User  User   `json:"user" gorm:"foreignKey:UserID"`
	Posts []Post `json:"posts" gorm:"many2many:post_social_accounts;"`
}

// Post represents a social media post
type Post struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id"`
	Title       string    `json:"title"`
	Content     string    `json:"content" gorm:"type:text"`
	MediaURLs   string    `json:"media_urls" gorm:"type:text"` // JSON array as string
	Hashtags    string    `json:"hashtags"`
	Status      string    `json:"status" gorm:"default:draft"` // draft, scheduled, published, failed
	ScheduledAt time.Time `json:"scheduled_at"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relationships
	User           User            `json:"user" gorm:"foreignKey:UserID"`
	SocialAccounts []SocialAccount `json:"social_accounts" gorm:"many2many:post_social_accounts;"`
	Analytics      []Analytics     `json:"analytics"`
}

// PostSocialAccount represents the many-to-many relationship between posts and social accounts
type PostSocialAccount struct {
	PostID          uint   `json:"post_id"`
	SocialAccountID uint   `json:"social_account_id"`
	Status          string `json:"status" gorm:"default:pending"` // pending, published, failed
	ExternalID      string `json:"external_id"`                   // ID from social platform
	ErrorMessage    string `json:"error_message"`
	PublishedAt     time.Time `json:"published_at"`
	
	Post          Post          `json:"post" gorm:"foreignKey:PostID"`
	SocialAccount SocialAccount `json:"social_account" gorm:"foreignKey:SocialAccountID"`
}

// Analytics represents engagement data for posts
type Analytics struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	PostID         uint      `json:"post_id"`
	SocialAccountID uint     `json:"social_account_id"`
	Platform       string    `json:"platform"`
	Likes          int       `json:"likes" gorm:"default:0"`
	Comments       int       `json:"comments" gorm:"default:0"`
	Shares         int       `json:"shares" gorm:"default:0"`
	Clicks         int       `json:"clicks" gorm:"default:0"`
	Impressions    int       `json:"impressions" gorm:"default:0"`
	Reach          int       `json:"reach" gorm:"default:0"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	// Relationships
	Post          Post          `json:"post" gorm:"foreignKey:PostID"`
	SocialAccount SocialAccount `json:"social_account" gorm:"foreignKey:SocialAccountID"`
}

// Subscription represents user subscription information
type Subscription struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	Plan      string    `json:"plan"`        // free, basic, professional, agency
	Status    string    `json:"status"`      // active, canceled, expired
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	User User `json:"user" gorm:"foreignKey:UserID"`
}