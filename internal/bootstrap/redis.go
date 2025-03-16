package bootstrap

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var RedisClient *redis.Client

func InitRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Redis default address
		DB:   0,                // Default DB
	})

	// Test connection
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		panic("Failed to connect to Redis: " + err.Error())
	}

	logrus.Info("Connected to Redis")
	return client
}
