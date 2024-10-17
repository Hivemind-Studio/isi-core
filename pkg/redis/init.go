package redis

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
)

// Global variables to hold the Redis client and options
var (
	redisClient  *redis.Client
	redisOptions *redis.Options
)

// InitRedis initializes the Redis client with the provided configuration.
func InitRedis(addr string, password string, db int) error {
	logger := utils.NewLogger()

	// Set up Redis options
	redisOptions = &redis.Options{
		Addr:     addr,     // Redis server address (e.g., "localhost:6379")
		Password: password, // Redis server password (if any)
		DB:       db,       // Redis database number
		PoolSize: 10,       // Connection pool size
	}

	// Create a new Redis client
	redisClient = redis.NewClient(redisOptions)

	// Test the connection
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		logger.Error("Failed Connect Redis %v", err)
		return err
	}

	logger.Success("Connected to Redis %s", addr)
	return nil
}

// GetRedisClient returns the initialized Redis client.
func GetRedisClient() *redis.Client {
	if redisClient == nil {
		log.Fatal("Redis client is not initialized")
	}
	return redisClient
}
