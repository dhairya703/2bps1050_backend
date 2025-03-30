package config

import (
	"github.com/redis/go-redis/v9"
	"context"
	"log"
	"os"
)

// RedisClient is the Redis connection
var RedisClient *redis.Client

// Initialize Redis connection
func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"), // e.g., "localhost:6379"
		Password: "",                      // No password by default
		DB:       0,                        // Default DB
	})

	// Ping Redis to check connection
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}
