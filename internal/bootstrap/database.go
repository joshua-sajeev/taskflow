package bootstrap

import (
	"taskflow/internal/db"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InitDatabase() (*gorm.DB, error) {
	database, err := db.InitDB()
	if err != nil {
		logrus.Warn("Couldn't Initialize Database")
		return nil, err
	}
	return database, nil
}
