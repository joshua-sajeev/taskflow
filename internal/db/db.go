package db

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB initializes the database connection
func InitDB() (*gorm.DB, error) {
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		logrus.Fatal("DATABASE_URL is not set in .env")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	logrus.Info("Connected to PostgreSQL database")
	return db, nil
}
