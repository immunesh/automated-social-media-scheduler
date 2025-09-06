package database

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/immunesh/automated-social-media-scheduler/internal/models"
)

var DB *gorm.DB

func Initialize() error {
	// For development, we'll use SQLite. In production, switch to PostgreSQL
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./scheduler.db"
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	// Auto-migrate the schema
	return AutoMigrate()
}

func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.SocialAccount{},
		&models.Post{},
		&models.PostSocialAccount{},
		&models.Analytics{},
		&models.Subscription{},
	)
}

func GetDB() *gorm.DB {
	return DB
}