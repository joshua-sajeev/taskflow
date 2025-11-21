package database

import (
	"context"
	"fmt"
	"log"
	"taskflow/pkg"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	defaultMaxRetries = 10
	defaultRetryDelay = 2 * time.Second
)

type Config struct {
	User       string
	Password   string
	Host       string
	Port       string
	Database   string
	MaxRetries int
	RetryDelay time.Duration
	LogLevel   logger.LogLevel
}

func LoadConfigFromEnv() *Config {
	return &Config{
		User:       pkg.GetEnv("MYSQL_USER", "appuser"),
		Password:   pkg.GetEnv("MYSQL_PASSWORD", "apppassword"),
		Host:       pkg.GetEnv("MYSQL_HOST", "mysql"),
		Port:       pkg.GetEnv("MYSQL_PORT", "3306"),
		Database:   pkg.GetEnv("MYSQL_DATABASE", "taskdb"),
		MaxRetries: defaultMaxRetries,
		RetryDelay: defaultRetryDelay,
		LogLevel:   logger.Info,
	}
}

func ConnectDB(cfg *Config) (*gorm.DB, error) {
	if cfg == nil {
		cfg = LoadConfigFromEnv()
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&timeout=500ms",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)

	log.Printf("Connecting to database %s@%s:%s/%s", cfg.User, cfg.Host, cfg.Port, cfg.Database)

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(cfg.LogLevel),
	}

	var gdb *gorm.DB
	var err error

	for i := 0; i < cfg.MaxRetries; i++ {
		gdb, err = gorm.Open(mysql.Open(dsn), gormConfig)
		if err == nil {
			sqlDB, dbErr := gdb.DB()
			if dbErr == nil {

				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
				pingErr := sqlDB.PingContext(ctx)
				cancel()

				if pingErr == nil {
					log.Println("Database connected successfully")

					sqlDB.SetMaxIdleConns(10)
					sqlDB.SetMaxOpenConns(100)
					sqlDB.SetConnMaxLifetime(time.Hour)

					return gdb, nil
				}

				err = pingErr
			} else {
				err = dbErr
			}
		}

		log.Printf("Database connection attempt %d/%d failed: %v. Retrying in %v...",
			i+1, cfg.MaxRetries, err, cfg.RetryDelay)
		time.Sleep(cfg.RetryDelay)
	}

	return nil, fmt.Errorf(
		"failed to connect to database after %d attempts: %w",
		cfg.MaxRetries, err,
	)
}

func MigrateModels(db *gorm.DB, models ...any) error {
	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("auto-migration failed: %w", err)
	}
	log.Println("Database migration completed successfully")
	return nil
}
