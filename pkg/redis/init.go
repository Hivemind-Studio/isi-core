package redis

import (
	"context"
	"github.com/Hivemind-Studio/isi-core/utils"
	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
)

var (
	redisClient  *redis.Client
	redisOptions *redis.Options
)

func InitRedis(addr string, password string, db int) *redis.Client {
	logger := utils.NewLogger()

	// Set up Redis options
	redisOptions = &redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		PoolSize: 10,
	}

	// Create a new Redis client
	redisClient = redis.NewClient(redisOptions)

	// Test the connection
	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		logger.Error("Failed Connect Redis %v", err)
		return nil
	}

	logger.Success("Connected to Redis %s", addr)
	return redisClient
}

func GetRedisClient() *redis.Client {
	if redisClient == nil {
		log.Fatal("Redis client is not initialized")
	}
	return redisClient
}
