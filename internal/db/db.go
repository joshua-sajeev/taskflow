package db

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	once       sync.Once
)

// InitDB initializes and returns a singleton database connection
func InitDB() (*gorm.DB, error) {
	var err error

	once.Do(func() {
		// Load environment variables
		if err = godotenv.Load(); err != nil {
			logrus.Warn("No .env file found")
		}

		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			logrus.Fatal("DATABASE_URL is not set in .env")
		}

		// Open database connection
		dbInstance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			err = fmt.Errorf("failed to connect to database: %w", err)
			return
		}

		logrus.Info("Connected to PostgreSQL database")
	})

	return dbInstance, err
}

func CloseDB() {
	sqlDB, err := dbInstance.DB()
	if err != nil {
		logrus.Error("Failed to get database instance:", err)
		return
	}
	sqlDB.Close()
	logrus.Info("Database connection closed")
}
